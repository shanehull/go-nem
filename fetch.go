package nem

import (
	"context"
	"fmt"
	"time"

	"github.com/shanehull/go-nem/internal"
)

// Download returns the raw bytes of a file, using the cache when one is configured.
func (c *Client) Download(ctx context.Context, ref FileRef) ([]byte, error) {
	if c.cache != nil {
		if data, ok := c.cache.Get(ref); ok {
			return data, nil
		}
	}

	if err := c.pace(ctx); err != nil {
		return nil, err
	}

	data, err := internal.Get(ctx, c.httpClient, ref.URL, c.userAgent)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if err := c.cache.Put(ref, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// FetchTables lists and downloads every file matching the options, and returns
// the parsed tables from all of them.
func (c *Client) FetchTables(ctx context.Context, report Report, opts ...ListOption) ([]Table, error) {
	if report.Kind != KindNested {
		return nil, fmt.Errorf("nem: %s is not a nested record report", report.Dir)
	}

	refs, err := c.List(ctx, report, opts...)
	if err != nil {
		return nil, err
	}

	var tables []Table
	for _, ref := range refs {
		data, err := c.Download(ctx, ref)
		if err != nil {
			return nil, err
		}
		parsed, err := Parse(data)
		if err != nil {
			return nil, fmt.Errorf("nem: parse %s: %w", ref.Name, err)
		}
		tables = append(tables, parsed...)
	}
	return tables, nil
}

// FetchNotices lists and downloads every market notice matching the options.
func (c *Client) FetchNotices(ctx context.Context, opts ...ListOption) ([]MarketNotice, error) {
	refs, err := c.List(ctx, ReportMarketNotice, opts...)
	if err != nil {
		return nil, err
	}

	var notices []MarketNotice
	for _, ref := range refs {
		data, err := c.Download(ctx, ref)
		if err != nil {
			return nil, err
		}
		notice, err := DecodeMarketNotice(data)
		if err != nil {
			return nil, fmt.Errorf("nem: decode notice %s: %w", ref.Name, err)
		}
		notices = append(notices, notice)
	}
	return notices, nil
}

func (c *Client) pace(ctx context.Context) error {
	if c.minInterval <= 0 {
		return nil
	}

	c.mu.Lock()
	now := time.Now()
	wait := time.Duration(0)
	if !c.lastFetch.IsZero() {
		wait = time.Until(c.lastFetch.Add(c.minInterval))
	}
	if wait < 0 {
		wait = 0
	}
	c.lastFetch = now.Add(wait)
	c.mu.Unlock()

	if wait == 0 {
		return nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
