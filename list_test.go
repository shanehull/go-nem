package nem_test

import (
	"context"
	"testing"

	"github.com/shanehull/go-nem"
)

const dispatchISDir = "/Reports/CURRENT/DispatchIS_Reports/"

func TestList(t *testing.T) {
	srv, _ := newServer(t, dispatchISDir, readFixture(t, "listing.html"), nil)
	client, err := nem.New(nem.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	refs, err := client.List(context.Background(), nem.ReportDispatchIS)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(refs) != 3 {
		t.Fatalf("got %d refs, want 3", len(refs))
	}
	if refs[0].Name != "PUBLIC_DISPATCHIS_202609210805_0000000539009108.zip" {
		t.Errorf("first ref = %s", refs[0].Name)
	}
	if refs[0].Size != 20629 {
		t.Errorf("first size = %d, want 20629", refs[0].Size)
	}
	if !refs[0].Modified.Equal(mustTime(t, "2026-09-21T08:05:00+10:00")) {
		t.Errorf("first modified = %s", refs[0].Modified)
	}
	if refs[2].Name != "PUBLIC_DISPATCHIS_202609210815_0000000539009632.zip" {
		t.Errorf("last ref = %s", refs[2].Name)
	}
}

func TestListLimit(t *testing.T) {
	srv, _ := newServer(t, dispatchISDir, readFixture(t, "listing.html"), nil)
	client, _ := nem.New(nem.WithBaseURL(srv.URL))

	refs, err := client.List(context.Background(), nem.ReportDispatchIS, nem.Limit(2))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("got %d refs, want 2", len(refs))
	}
	if refs[0].Name != "PUBLIC_DISPATCHIS_202609210810_0000000539009107.zip" {
		t.Errorf("first ref = %s, want the 0810 file", refs[0].Name)
	}
}

func TestListOnName(t *testing.T) {
	srv, _ := newServer(t, dispatchISDir, readFixture(t, "listing.html"), nil)
	client, _ := nem.New(nem.WithBaseURL(srv.URL))

	match, err := client.List(context.Background(), nem.ReportDispatchIS,
		nem.OnName(mustTime(t, "2026-09-21T00:00:00+10:00")))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(match) != 3 {
		t.Fatalf("got %d refs, want 3", len(match))
	}

	none, err := client.List(context.Background(), nem.ReportDispatchIS,
		nem.OnName(mustTime(t, "2026-09-22T00:00:00+10:00")))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("got %d refs, want 0", len(none))
	}
}

func TestListSince(t *testing.T) {
	srv, _ := newServer(t, dispatchISDir, readFixture(t, "listing.html"), nil)
	client, _ := nem.New(nem.WithBaseURL(srv.URL))

	refs, err := client.List(context.Background(), nem.ReportDispatchIS,
		nem.Since(mustTime(t, "2026-09-21T08:10:00+10:00")))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("got %d refs, want 2", len(refs))
	}
}
