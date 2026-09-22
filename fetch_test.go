package nem_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/shanehull/go-nem"
)

func TestFetchTablesUsesCache(t *testing.T) {
	zipped := zipBytes(t, "PUBLIC_DISPATCHIS_202609210815_0000000539009632.CSV", readFixture(t, "dispatchis.csv"))
	srv, hits := newServer(t, dispatchISDir, readFixture(t, "listing.html"), zipped)
	client, err := nem.New(nem.WithBaseURL(srv.URL), nem.WithCacheDir(t.TempDir()))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	tables, err := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	if err != nil {
		t.Fatalf("FetchTables: %v", err)
	}
	prices, err := nem.DecodeDispatchPrice(tables)
	if err != nil {
		t.Fatalf("DecodeDispatchPrice: %v", err)
	}
	if len(prices) != 2 {
		t.Fatalf("got %d prices, want 2", len(prices))
	}
	if got := atomic.LoadInt32(hits); got != 1 {
		t.Fatalf("download count = %d, want 1", got)
	}

	if _, err := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1)); err != nil {
		t.Fatalf("FetchTables (cached): %v", err)
	}
	if got := atomic.LoadInt32(hits); got != 1 {
		t.Errorf("download count after cached fetch = %d, want 1", got)
	}
}

func TestFetchNotices(t *testing.T) {
	srv, _ := newServer(t, "/Reports/CURRENT/Market_Notice/", readFixture(t, "market_notice_listing.html"), readFixture(t, "market_notice.txt"))
	client, err := nem.New(nem.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	notices, err := client.FetchNotices(context.Background())
	if err != nil {
		t.Fatalf("FetchNotices: %v", err)
	}
	if len(notices) != 1 {
		t.Fatalf("got %d notices, want 1", len(notices))
	}
	if notices[0].NoticeID != 144924 {
		t.Errorf("NoticeID = %d, want 144924", notices[0].NoticeID)
	}
}
