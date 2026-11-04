package utils

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setRequestLimits(t *testing.T, timeout time.Duration, megabytes int64) {
	t.Helper()
	oldTimeout, oldSize := RequestTimeout, MaxResponseMB
	t.Cleanup(func() { RequestTimeout, MaxResponseMB = oldTimeout, oldSize })
	RequestTimeout, MaxResponseMB = timeout, megabytes
}

type countedBody struct {
	io.Reader
	read   int
	closed bool
}

func (b *countedBody) Read(p []byte) (int, error) {
	n, err := b.Reader.Read(p)
	b.read += n
	return n, err
}

func (b *countedBody) Close() error { b.closed = true; return nil }

func TestRequestsHaveDeadline(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if _, ok := r.Context().Deadline(); !ok {
			t.Error("request has no deadline")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	if _, err := Get("https://api.example.test/"); err != nil {
		t.Fatal(err)
	}
	if _, err := PostJSON("https://api.example.test/", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRequestsRejectAdvertisedOversizedBody(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, ContentLength: 1 << 62, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	for _, call := range []func(string) ([]byte, error){Get, func(u string) ([]byte, error) { return PostJSON(u, nil, nil) }} {
		if body, err := call("https://api.example.test/"); err == nil || body != nil {
			t.Errorf("oversized response was accepted: body=%q err=%v", body, err)
		}
	}
}

func TestRequestSizeBoundaries(t *testing.T) {
	setRequestLimits(t, time.Minute, 1)
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		for _, size := range []int{999_999, 1_000_000, 1_000_001} {
			for _, knownLength := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%d/known=%v", method, size, knownLength), func(t *testing.T) {
					body := &countedBody{Reader: strings.NewReader(strings.Repeat("x", size))}
					var requestContext context.Context
					http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						requestContext = r.Context()
						length := int64(-1)
						if knownLength {
							length = int64(size)
						}
						return &http.Response{StatusCode: http.StatusOK, ContentLength: length, Body: body}, nil
					})}
					result, err := request(method, "https://api.example.test/?apiKey=fixture-secret", nil, nil)
					if size <= 1_000_000 {
						if err != nil || len(result) != size {
							t.Fatalf("valid body changed: length=%d err=%v", len(result), err)
						}
					} else if err == nil || result != nil || !strings.Contains(err.Error(), "response exceeds 1 MB") || strings.Contains(err.Error(), "fixture-secret") {
						t.Fatalf("oversized response not rejected cleanly: length=%d err=%v", len(result), err)
					}
					if !body.closed || body.read > 1_000_001 || requestContext.Err() != context.Canceled {
						t.Fatalf("response not stopped: closed=%v read=%d context=%v", body.closed, body.read, requestContext.Err())
					}
					if knownLength && size > 1_000_000 && body.read != 0 {
						t.Fatal("read a body already known to be too large")
					}
				})
			}
		}
	}
}

func TestRequestLimitAppliesAfterGzipDecoding(t *testing.T) {
	setRequestLimits(t, time.Minute, 1)
	for _, size := range []int{1_000_000, 1_000_001} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Encoding", "gzip")
				writer := gzip.NewWriter(w)
				io.WriteString(writer, strings.Repeat("x", size))
				writer.Close()
			}))
			defer server.Close()
			body, err := Get(server.URL)
			if size == 1_000_000 {
				if err != nil || len(body) != size {
					t.Fatalf("valid decoded body changed: length=%d err=%v", len(body), err)
				}
			} else if err == nil || body != nil || !strings.Contains(err.Error(), "response exceeds") {
				t.Fatalf("oversized decoded body accepted: length=%d err=%v", len(body), err)
			}
		})
	}
}

func TestRequestTimeoutStopsHeadersAndBody(t *testing.T) {
	setRequestLimits(t, 100*time.Millisecond, 1)
	for _, phase := range []string{"headers", "body", "continuous-body"} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			t.Run(method+"/"+phase, func(t *testing.T) {
				closed := make(chan struct{})
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer close(closed)
					if phase != "headers" {
						io.WriteString(w, "prefix")
						w.(http.Flusher).Flush()
					}
					ticker := time.NewTicker(10 * time.Millisecond)
					defer ticker.Stop()
					for {
						select {
						case <-r.Context().Done():
							return
						case <-ticker.C:
							if phase == "continuous-body" {
								io.WriteString(w, "x")
								w.(http.Flusher).Flush()
							}
						}
					}
				}))
				defer server.Close()
				done := make(chan struct{})
				var body []byte
				var err error
				go func() {
					body, err = request(method, server.URL+"?apiKey=fixture-secret", nil, nil)
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(3 * time.Second):
					server.CloseClientConnections()
					<-done
					t.Fatal("request did not stop at its deadline")
				}
				if !errors.Is(err, context.DeadlineExceeded) || body != nil || !strings.Contains(err.Error(), "request timed out") || strings.Contains(err.Error(), "fixture-secret") {
					t.Fatalf("unsafe timeout result: body=%q err=%v", body, err)
				}
				select {
				case <-closed:
				case <-time.After(3 * time.Second):
					t.Fatal("server request was not canceled")
				}
			})
		}
	}
}

func TestRedirectsKeepOneDeadline(t *testing.T) {
	setRequestLimits(t, 10*time.Minute, 512)
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	var deadline time.Time
	requests := 0
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		current, ok := r.Context().Deadline()
		if !ok {
			t.Fatal("missing deadline")
		}
		if requests == 0 {
			deadline = current
			if remaining := time.Until(current); remaining < 9*time.Minute || remaining > 10*time.Minute {
				t.Fatalf("timeout override not used: %s", remaining)
			}
		} else if current != deadline {
			t.Fatal("redirect changed the deadline")
		}
		requests++
		response := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{}`))}
		if requests < 3 {
			response.StatusCode = http.StatusFound
			response.Header.Set("Location", fmt.Sprintf("/step%d", requests))
		}
		return response, nil
	})}
	if _, err := Get("https://api.example.test/"); err != nil || requests != 3 {
		t.Fatalf("redirect handling changed: requests=%d err=%v", requests, err)
	}
}

func TestInvalidLimitsDoNotSendRequest(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Error("sent request with invalid limits")
		return nil, errors.New("unexpected request")
	})}
	for _, tc := range []struct {
		timeout time.Duration
		size    int64
	}{{0, 1}, {-1, 1}, {time.Second, 0}, {time.Second, -1}, {time.Second, math.MaxInt64/1_000_000 + 1}} {
		t.Run(fmt.Sprintf("%s/%d", tc.timeout, tc.size), func(t *testing.T) {
			setRequestLimits(t, tc.timeout, tc.size)
			for _, call := range []func(string) ([]byte, error){Get, func(u string) ([]byte, error) { return PostJSON(u, nil, nil) }} {
				if body, err := call("https://api.example.test/"); err == nil || body != nil {
					t.Fatalf("invalid limits accepted: %q %v", body, err)
				}
			}
		})
	}
	setRequestLimits(t, time.Duration(math.MaxInt64), math.MaxInt64/1_000_000)
	if err := ValidateRequestLimits(); err != nil {
		t.Fatalf("valid large limits rejected: %v", err)
	}
}
