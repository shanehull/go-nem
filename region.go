package nem

import (
	"errors"
	"fmt"
	"strings"
)

// Region is a NEM market region.
type Region string

// The NEM market regions.
const (
	RegionNSW1 Region = "NSW1"
	RegionQLD1 Region = "QLD1"
	RegionSA1  Region = "SA1"
	RegionTAS1 Region = "TAS1"
	RegionVIC1 Region = "VIC1"
)

// ErrRegionInvalid is returned for a region that is not a NEM market region.
var ErrRegionInvalid = errors.New("nem: region is invalid")

// Regions returns the NEM market regions.
func Regions() []Region {
	return []Region{RegionNSW1, RegionQLD1, RegionSA1, RegionTAS1, RegionVIC1}
}

// ParseRegion returns the region for a case-insensitive code.
func ParseRegion(s string) (Region, error) {
	r := Region(strings.ToUpper(strings.TrimSpace(s)))
	if !r.Valid() {
		return "", fmt.Errorf("nem: %q: %w", s, ErrRegionInvalid)
	}
	return r, nil
}

// Valid reports whether r is a NEM market region.
func (r Region) Valid() bool {
	switch r {
	case RegionNSW1, RegionQLD1, RegionSA1, RegionTAS1, RegionVIC1:
		return true
	default:
		return false
	}
}

// String returns the region code.
func (r Region) String() string { return string(r) }

// Name returns the full region name. It is empty for an unknown region.
func (r Region) Name() string {
	switch r {
	case RegionNSW1:
		return "New South Wales"
	case RegionQLD1:
		return "Queensland"
	case RegionSA1:
		return "South Australia"
	case RegionTAS1:
		return "Tasmania"
	case RegionVIC1:
		return "Victoria"
	default:
		return ""
	}
}
