package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
	req, err := http.NewRequest(method, apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: invalid URL")
	}

	// API keys can appear in transport errors and in error responses from the server.
	defer func() {
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

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	return body, nil
}
