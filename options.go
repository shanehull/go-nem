package nem

import "time"

type listParams struct {
	since time.Time
	until time.Time
	limit int
	tier  Tier
}

// ListOption configures a List call.
type ListOption func(*listParams)

// Since filters to files modified at or after t.
func Since(t time.Time) ListOption {
	return func(p *listParams) { p.since = t }
}

// Until filters to files modified strictly before t.
func Until(t time.Time) ListOption {
	return func(p *listParams) { p.until = t }
}

// On filters to files modified on the given day, in NEM time.
func On(day time.Time) ListOption {
	return func(p *listParams) {
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, NEM)
		p.since = start
		p.until = start.AddDate(0, 0, 1)
	}
}

// Limit caps the number of files returned, keeping the newest.
func Limit(n int) ListOption {
	return func(p *listParams) { p.limit = n }
}

// FromArchive lists the archive tier instead of the current tier.
func FromArchive() ListOption {
	return func(p *listParams) { p.tier = TierArchive }
}
