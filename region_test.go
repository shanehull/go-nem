package nem_test

import (
	"errors"
	"testing"

	"github.com/shanehull/go-nem"
)

func TestRegions(t *testing.T) {
	regions := nem.Regions()
	if len(regions) != 5 {
		t.Fatalf("got %d regions, want 5", len(regions))
	}
	for _, region := range regions {
		if !region.Valid() {
			t.Errorf("%q is not valid", region)
		}
		if region.Name() == "" {
			t.Errorf("%q has no name", region)
		}
	}
}

func TestParseRegion(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    nem.Region
		wantErr bool
	}{
		{"exact", "NSW1", nem.RegionNSW1, false},
		{"lowercase", "vic1", nem.RegionVIC1, false},
		{"padded", "  qld1  ", nem.RegionQLD1, false},
		{"unknown", "WA1", "", true},
		{"empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nem.ParseRegion(tt.in)
			if tt.wantErr {
				if !errors.Is(err, nem.ErrRegionInvalid) {
					t.Fatalf("error = %v, want ErrRegionInvalid", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRegion(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseRegion(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRegionValid(t *testing.T) {
	if nem.Region("NSW").Valid() {
		t.Error("NSW is not a NEM market region")
	}
	if !nem.RegionNSW1.Valid() {
		t.Error("NSW1 should be valid")
	}
}

func TestRegionName(t *testing.T) {
	if got := nem.RegionVIC1.Name(); got != "Victoria" {
		t.Errorf("Name = %q, want Victoria", got)
	}
	if got := nem.Region("BAD").Name(); got != "" {
		t.Errorf("Name for unknown region = %q, want empty", got)
	}
}
