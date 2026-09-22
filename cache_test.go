package nem_test

import (
	"testing"

	"github.com/shanehull/go-nem"
)

func TestCachePutGet(t *testing.T) {
	cache, err := nem.NewCache(t.TempDir())
	if err != nil {
		t.Fatalf("NewCache: %v", err)
	}
	ref := nem.FileRef{Report: nem.ReportDispatchIS, Name: "PUBLIC_DISPATCHIS_202609210810_0000000539009107.zip"}

	if _, ok := cache.Get(ref); ok {
		t.Fatal("expected a cache miss")
	}
	if err := cache.Put(ref, []byte("payload")); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, ok := cache.Get(ref)
	if !ok {
		t.Fatal("expected a cache hit")
	}
	if string(got) != "payload" {
		t.Errorf("got %q, want payload", got)
	}
}

func TestCacheEmptyDir(t *testing.T) {
	if _, err := nem.NewCache(""); err == nil {
		t.Fatal("expected error for empty cache dir")
	}
}
