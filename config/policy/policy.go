package policy

import (
//"fmt"

//"github.com/hashicorp/hcl"
//"github.com/hashicorp/hcl/hcl/ast"
)

type Policy struct {
	Name            string `xml:"name,attr" hcl:"-"`
	Enabled         bool   `xml:"enabled,attr" hcl:"enabled"`
	ContinueOnError bool   `xml:"continueOnError,attr" hcl:"continue_on_error"`
	Async           bool   `xml:"async,attr" hcl:"async"`
}
