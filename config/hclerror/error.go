package hclerror

import (
	"fmt"
	"github.com/hashicorp/hcl/hcl/token"
)

type PosError struct {
	Pos token.Pos
	Err error
}

func (e *PosError) Error() string {
	return fmt.Sprintf("%s (at %s, line %d, col %d)",
		e.Err, e.Pos.Filename, e.Pos.Line, e.Pos.Column)
}
