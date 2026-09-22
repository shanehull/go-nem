package nem

import (
	"fmt"
	"strings"
	"time"
)

// NEM is the NEM time zone, AEST (UTC+10), which does not observe daylight saving.
var NEM = time.FixedZone("NEM", 10*60*60)

const (
	dispatchInterval = 5 * time.Minute
	tradingInterval  = 30 * time.Minute
	tradingDayStart  = 4 * time.Hour
)

// ParseNEMTime parses an AEMO date-time. MMS reports use the year-first form
// "2006/01/02 15:04:05" and market notices use the day-first form
// "02/01/2006 15:04:05". The layout is detected from the first field, and runs
// of whitespace are collapsed because notice fields pad with spaces.
func ParseNEMTime(s string) (time.Time, error) {
	normalized := strings.Join(strings.Fields(s), " ")
	t, err := time.ParseInLocation(dateTimeLayout(normalized), normalized, NEM)
	if err != nil {
		return time.Time{}, fmt.Errorf("nem: parse NEM time %q: %w", s, err)
	}
	return t, nil
}

// ParseNEMDate parses an AEMO date, detecting the year-first "2006/01/02" and
// day-first "02/01/2006" forms.
func ParseNEMDate(s string) (time.Time, error) {
	trimmed := strings.TrimSpace(s)
	t, err := time.ParseInLocation(dateLayout(trimmed), trimmed, NEM)
	if err != nil {
		return time.Time{}, fmt.Errorf("nem: parse NEM date %q: %w", s, err)
	}
	return t, nil
}

// dateTimeLayout returns the date-time layout for s. A two-digit first field
// means day-first; otherwise year-first.
func dateTimeLayout(s string) string {
	if strings.IndexByte(s, '/') == 2 {
		return "02/01/2006 15:04:05"
	}
	return "2006/01/02 15:04:05"
}

// dateLayout returns the date layout for s. A two-digit first field means
// day-first; otherwise year-first.
func dateLayout(s string) string {
	if strings.IndexByte(s, '/') == 2 {
		return "02/01/2006"
	}
	return "2006/01/02"
}

// InNEM converts t to NEM time.
func InNEM(t time.Time) time.Time { return t.In(NEM) }

// TradingDay returns the NEM trading day, as "2006-01-02", that contains t.
// A trading day runs from 04:05 to 04:00 the following day.
func TradingDay(t time.Time) string {
	n := t.In(NEM)
	secs := n.Hour()*3600 + n.Minute()*60 + n.Second()
	if secs <= int(tradingDayStart.Seconds()) {
		n = n.AddDate(0, 0, -1)
	}
	return n.Format("2006-01-02")
}

// IntervalBeginning returns the start of the dispatch interval ending at t.
func IntervalBeginning(t time.Time) time.Time { return t.Add(-dispatchInterval) }

// TradingIntervalBeginning returns the start of the trading interval ending at t.
func TradingIntervalBeginning(t time.Time) time.Time { return t.Add(-tradingInterval) }
