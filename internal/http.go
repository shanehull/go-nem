// Package internal provides the HTTP transport shared by the client.
package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// APIError represents a non-success response from NEMWEB.
type APIError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("nem: HTTP %d for %s", e.StatusCode, e.URL)
}

// Get performs a GET request and returns the response body.
func Get(ctx context.Context, client *http.Client, url, userAgent string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("nem: request creation failed: %w", err)
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nem: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nem: reading response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, URL: url, Body: string(body)}
	}

	return body, nil
}
