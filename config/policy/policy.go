package policy

import (
//"fmt"

//"github.com/hashicorp/hcl"
//"github.com/hashicorp/hcl/hcl/ast"
)

type Policy struct {
	Name            string `xml:"name,attr,omitempty" hcl:"-"`
	Enabled         bool   `xml:"enabled,attr,omitempty" hcl:"enabled"`
	ContinueOnError bool   `xml:"continueOnError,attr,omitempty" hcl:"continue_on_error"`
	Async           bool   `xml:"async,attr,omitempty" hcl:"async"`
}
