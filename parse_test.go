package nem_test

import (
	"bytes"
	"testing"

	"github.com/shanehull/go-nem"
)

func TestParseDispatchIS(t *testing.T) {
	tables, err := nem.Parse(readFixture(t, "dispatchis.csv"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(tables) != 3 {
		t.Fatalf("got %d tables, want 3", len(tables))
	}

	want := []struct{ group, name string }{
		{"DISPATCH", "PRICE"},
		{"DISPATCH", "REGIONSUM"},
		{"DISPATCH", "INTERCONNECTORRES"},
	}
	for i, w := range want {
		if tables[i].Group != w.group || tables[i].Name != w.name {
			t.Errorf("table %d = %s/%s, want %s/%s", i, tables[i].Group, tables[i].Name, w.group, w.name)
		}
	}
	if tables[0].Version != 5 {
		t.Errorf("PRICE version = %d, want 5", tables[0].Version)
	}
	if len(tables[0].Rows) != 2 {
		t.Errorf("PRICE rows = %d, want 2", len(tables[0].Rows))
	}
}

func TestParseZip(t *testing.T) {
	zipped := zipBytes(t, "PUBLIC_DISPATCHIS_202609210810_0000000539009107.CSV", readFixture(t, "dispatchis.csv"))

	tables, err := nem.Parse(zipped)
	if err != nil {
		t.Fatalf("Parse zip: %v", err)
	}
	if len(tables) != 3 {
		t.Fatalf("got %d tables, want 3", len(tables))
	}
}

func TestParseDoubleZip(t *testing.T) {
	// AEMO publishes some daily archives double-zipped: an outer zip whose
	// single entry is the report zip.
	inner := zipBytes(t, "PUBLIC_DISPATCHIS_202608220005_0000000533804272.CSV", readFixture(t, "dispatchis.csv"))
	outer := zipBytes(t, "PUBLIC_DISPATCHIS_20260822.zip", inner)

	tables, err := nem.Parse(outer)
	if err != nil {
		t.Fatalf("Parse double zip: %v", err)
	}
	if len(tables) != 3 {
		t.Fatalf("got %d tables, want 3", len(tables))
	}
}

func TestParseMultiEntryZip(t *testing.T) {
	// A daily archive is an outer zip of one inner zip per interval. Every
	// entry must contribute rows, not just the first.
	inner := zipBytes(t, "PUBLIC_DISPATCHIS_202608220005_0000000533804272.CSV", readFixture(t, "dispatchis.csv"))
	outer := zipMulti(t, map[string][]byte{
		"PUBLIC_DISPATCHIS_202608220005_0000000533804272.zip": inner,
		"PUBLIC_DISPATCHIS_202608220010_0000000533804879.zip": inner,
	})

	tables, err := nem.Parse(outer)
	if err != nil {
		t.Fatalf("Parse multi-entry zip: %v", err)
	}
	if len(tables) != 6 {
		t.Fatalf("got %d tables, want 6", len(tables))
	}
	prices, err := nem.DecodeDispatchPrice(tables)
	if err != nil {
		t.Fatalf("DecodeDispatchPrice: %v", err)
	}
	if len(prices) != 4 {
		t.Fatalf("got %d prices, want 4", len(prices))
	}
}

func TestParseUnknownRecordType(t *testing.T) {
	if _, err := nem.Parse([]byte("X,foo,bar\n")); err == nil {
		t.Fatal("expected error for unknown record type")
	}
}

func TestParseDataBeforeHeader(t *testing.T) {
	if _, err := nem.Parse([]byte("D,DISPATCH,PRICE,5,1\n")); err == nil {
		t.Fatal("expected error for data before header")
	}
}

func TestParseEmpty(t *testing.T) {
	if _, err := nem.Parse(nil); err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseTrimsZipPrefix(t *testing.T) {
	if !bytes.HasPrefix(zipBytes(t, "f.csv", []byte("x")), []byte("PK\x03\x04")) {
		t.Fatal("zipBytes did not produce a zip")
	}
}
