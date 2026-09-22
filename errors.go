package nem

import (
	"errors"

	"github.com/shanehull/go-nem/internal"
)

// APIError represents a non-success response from NEMWEB.
// It is a type alias, so callers use errors.As.
type APIError = internal.APIError

// ErrTableNotFound is returned when a typed decoder cannot find its table.
var ErrTableNotFound = errors.New("nem: table not found")
