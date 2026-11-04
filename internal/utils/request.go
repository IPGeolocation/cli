package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var RequestTimeout = 5 * time.Minute
var MaxResponseMB int64 = 256

func ValidateRequestLimits() error {
	if RequestTimeout <= 0 {
		return fmt.Errorf("request timeout must be greater than zero")
	}
	if MaxResponseMB <= 0 {
		return fmt.Errorf("response size limit must be greater than zero")
	}
	if MaxResponseMB > math.MaxInt64/1_000_000 {
		return fmt.Errorf("response size limit is too large")
	}
	return nil
}

func Get(apiURL string) ([]byte, error) {
	return request(http.MethodGet, apiURL, nil, nil)
}

func PostJSON(apiURL string, payload interface{}, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	return request(http.MethodPost, apiURL, bytes.NewReader(data), headers)
}

type requestError struct {
	cause   error
	message string
}

func (e *requestError) Error() string { return e.message }
func (e *requestError) Unwrap() error { return e.cause }

func request(method, apiURL string, data io.Reader, headers map[string]string) (body []byte, err error) {
	if err := ValidateRequestLimits(); err != nil {
		return nil, err
	}
	maxBytes := MaxResponseMB * 1_000_000
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: invalid URL")
	}

	// API keys can appear in transport errors and in error responses from the server.
	defer func() {
		if err != nil && ctx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("request timed out; increase --timeout for longer requests: %w", err)
		}
		if key := req.URL.Query().Get("apiKey"); err != nil && key != "" {
			encoded, _ := json.Marshal(key)
			message := strings.NewReplacer(
				key, "[REDACTED]", url.QueryEscape(key), "[REDACTED]",
				url.PathEscape(key), "[REDACTED]", string(encoded[1:len(encoded)-1]), "[REDACTED]",
			).Replace(err.Error())
			err = &requestError{cause: err, message: message}
		}
	}()

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.ContentLength > maxBytes {
		return nil, fmt.Errorf("response exceeds %d MB; increase --max-response-mb for larger results", MaxResponseMB)
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("response exceeds %d MB; increase --max-response-mb for larger results", MaxResponseMB)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	return body, nil
}
