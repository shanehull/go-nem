// Package nem provides a client for AEMO National Electricity Market data
// published on NEMWEB.
//
// It discovers report files, downloads and caches them, and parses the two
// wire formats AEMO publishes: the nested record format used by MMS reports
// and the human-formatted text used by market notices.
//
// Usage:
//
//	client, err := nem.New(nem.WithCacheDir("/var/lib/nem"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	notices, err := client.FetchNotices(ctx, nem.Since(time.Now().Add(-24*time.Hour)))
package nem

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL   = "https://nemweb.com.au"
	defaultTimeout   = 60 * time.Second
	defaultUserAgent = "go-nem (+https://github.com/shanehull/go-nem)"
)

// Client is a NEMWEB client. Must be constructed via New.
// Safe for concurrent use by multiple goroutines.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	userAgent   string
	cache       *Cache
	minInterval time.Duration

	mu        sync.Mutex
	lastFetch time.Time
}

// New creates a NEMWEB client.
func New(opts ...ClientOption) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		userAgent:  defaultUserAgent,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// ClientOption configures a Client.
type ClientOption func(*Client) error

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("nem: http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		if baseURL == "" {
			return fmt.Errorf("nem: base URL cannot be empty")
		}
		c.baseURL = strings.TrimRight(baseURL, "/")
		return nil
	}
}

// WithCacheDir enables on-disk caching of downloaded files.
func WithCacheDir(dir string) ClientOption {
	return func(c *Client) error {
		cache, err := NewCache(dir)
		if err != nil {
			return err
		}
		c.cache = cache
		return nil
	}
}

// WithUserAgent sets the User-Agent header sent with each request.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) error {
		c.userAgent = ua
		return nil
	}
}

// WithMinFetchInterval sets the minimum spacing between network downloads.
// It is disabled by default.
func WithMinFetchInterval(d time.Duration) ClientOption {
	return func(c *Client) error {
		if d < 0 {
			return fmt.Errorf("nem: min fetch interval cannot be negative")
		}
		c.minInterval = d
		return nil
	}
}
