package nem_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-nem"
)

func TestLiveMarketNotices(t *testing.T) {
	if os.Getenv("NEM_LIVE") == "" {
		t.Skip("set NEM_LIVE=1 to run live NEMWEB tests")
	}

	client, err := nem.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	refs, err := client.List(context.Background(), nem.ReportMarketNotice, nem.Limit(5))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(refs) == 0 {
		t.Fatal("no market notices returned")
	}

	notices, err := client.FetchNotices(context.Background(), nem.Limit(1))
	if err != nil {
		t.Fatalf("FetchNotices: %v", err)
	}
	if len(notices) == 0 || notices[0].NoticeID == 0 {
		t.Fatalf("unexpected notices: %+v", notices)
	}
}
