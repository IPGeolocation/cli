package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type brokenBody struct{ closed bool }

func (*brokenBody) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (b *brokenBody) Close() error           { b.closed = true; return nil }

func TestPostJSONDoesNotExposeAPIKey(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})}
	_, err := PostJSON("https://api.example.test/?apiKey=private-test-key", nil, nil)
	if err == nil || strings.Contains(err.Error(), "private-test-key") {
		t.Fatalf("request error exposes the key: %v", err)
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("request error lost its cause: %v", err)
	}
}

func TestPostJSONRejectsIncompleteResponse(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	body := &brokenBody{}
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})}
	data, err := PostJSON("https://api.example.test/", nil, nil)
	if !errors.Is(err, io.ErrUnexpectedEOF) || data != nil || !body.closed {
		t.Fatalf("incomplete response accepted or body left open: data=%q err=%v closed=%v", data, err, body.closed)
	}
}

func TestPostJSONPreservesPayloadAndHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != `{"uaString":"quoted \"agent\""}` || r.Method != http.MethodPost ||
			r.Header.Get("Content-Type") != "application/json" || r.Header.Get("X-Test") != "value" {
			t.Errorf("unexpected request: %s %s headers=%v body=%s err=%v", r.Method, r.URL, r.Header, body, err)
		}
		io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()
	body, err := PostJSON(server.URL, map[string]string{"uaString": `quoted "agent"`}, map[string]string{"X-Test": "value"})
	if err != nil || string(body) != `{"ok":true}` {
		t.Fatalf("unexpected result: %q %v", body, err)
	}
}

func TestGetTLSVerification(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()
	if _, err := Get(server.URL); err == nil {
		t.Fatal("accepted an untrusted TLS certificate")
	}
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = server.Client()
	if body, err := Get(server.URL); err != nil || string(body) != `{"ok":true}` {
		t.Fatalf("trusted TLS request failed: %s %v", body, err)
	}
}

func TestRequestErrorsRedactPlainAndEncodedKeys(t *testing.T) {
	key := "private +&key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, key+" "+url.QueryEscape(key)+" "+url.PathEscape(key))
		json.NewEncoder(w).Encode(key)
	}))
	defer server.Close()
	for _, call := range []func(string) ([]byte, error){Get, func(u string) ([]byte, error) { return PostJSON(u, nil, nil) }} {
		_, err := call(server.URL + "?apiKey=" + url.QueryEscape(key))
		encoded, _ := json.Marshal(key)
		if err == nil || !strings.Contains(err.Error(), "server error") || strings.Contains(err.Error(), key) || strings.Contains(err.Error(), url.QueryEscape(key)) || strings.Contains(err.Error(), url.PathEscape(key)) || strings.Contains(err.Error(), string(encoded[1:len(encoded)-1])) {
			t.Fatalf("unsafe status error: %v", err)
		}
		if _, err := call("https://example.test/%invalid?apiKey=" + key); err == nil || strings.Contains(err.Error(), key) {
			t.Fatalf("unsafe URL error: %v", err)
		}
	}
}
