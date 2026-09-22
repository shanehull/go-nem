package nem

import "time"

// Kind describes how a report file is encoded.
type Kind int

const (
	// KindNested is the AEMO nested record format, with I, D, and C rows.
	KindNested Kind = iota
	// KindText is a human-formatted text report, such as a market notice.
	KindText
)

// Tier is a NEMWEB publication tier.
type Tier int

const (
	// TierCurrent is the live publication tier.
	TierCurrent Tier = iota
	// TierArchive is the historical publication tier.
	TierArchive
)

func (t Tier) String() string {
	if t == TierArchive {
		return "ARCHIVE"
	}
	return "CURRENT"
}

// Report identifies a NEMWEB report directory and its file naming.
type Report struct {
	// Dir is the report directory under /Reports/{CURRENT,ARCHIVE}/.
	Dir string
	// Prefix is the file name prefix shared by every file in the report.
	Prefix string
	// Kind is how the report files are encoded.
	Kind Kind
}

// FileRef references a single file published under a report directory.
type FileRef struct {
	Report   Report
	Name     string
	URL      string
	Size     int64
	Modified time.Time
	Tier     Tier
}

// Table is one table parsed from a nested record report. Columns are the
// header defined by the I row, and Rows hold the D row values aligned to them.
type Table struct {
	Group   string
	Name    string
	Version int
	Columns []string
	Rows    [][]string
}
