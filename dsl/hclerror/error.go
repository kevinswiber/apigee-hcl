package hclerror

import (
	"fmt"
	"github.com/hashicorp/hcl/hcl/token"
)

// PosError contains an error along with its position in an HCL file
type PosError struct {
	Pos token.Pos
	Err error
}

// Error implements the error interface
func (e *PosError) Error() string {
	return fmt.Sprintf("%s (at %s, line %d, col %d)",
		e.Err, e.Pos.Filename, e.Pos.Line, e.Pos.Column)
}
