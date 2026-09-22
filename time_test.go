package nem_test

import (
	"testing"
	"time"

	"github.com/shanehull/go-nem"
)

func TestParseNEMTime(t *testing.T) {
	got, err := nem.ParseNEMTime("26/08/2026     11:18:33")
	if err != nil {
		t.Fatalf("ParseNEMTime: %v", err)
	}
	if !got.Equal(mustTime(t, "2026-08-26T11:18:33+10:00")) {
		t.Errorf("got %s", got)
	}
}

func TestParseNEMDate(t *testing.T) {
	got, err := nem.ParseNEMDate(" 26/08/2026 ")
	if err != nil {
		t.Fatalf("ParseNEMDate: %v", err)
	}
	if !got.Equal(mustTime(t, "2026-08-26T00:00:00+10:00")) {
		t.Errorf("got %s", got)
	}
}

func TestTradingDay(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want string
	}{
		{"start boundary", time.Date(2026, 9, 21, 4, 0, 0, 0, nem.NEM), "2026-09-20"},
		{"first interval", time.Date(2026, 9, 21, 4, 5, 0, 0, nem.NEM), "2026-09-21"},
		{"late morning", time.Date(2026, 9, 21, 11, 0, 0, 0, nem.NEM), "2026-09-21"},
		{"evening peak", time.Date(2026, 9, 21, 18, 30, 0, 0, nem.NEM), "2026-09-21"},
		{"just before boundary", time.Date(2026, 9, 21, 3, 59, 0, 0, nem.NEM), "2026-09-20"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nem.TradingDay(tt.in); got != tt.want {
				t.Errorf("TradingDay(%s) = %s, want %s", tt.in.Format(time.RFC3339), got, tt.want)
			}
		})
	}
}

func TestIntervalBeginning(t *testing.T) {
	ending := time.Date(2026, 9, 21, 8, 10, 0, 0, nem.NEM)
	want := time.Date(2026, 9, 21, 8, 5, 0, 0, nem.NEM)
	if got := nem.IntervalBeginning(ending); !got.Equal(want) {
		t.Errorf("IntervalBeginning = %s, want %s", got, want)
	}
}
