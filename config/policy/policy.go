package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Policy struct {
	Name            string `xml:"name,attr,omitempty" hcl:"-"`
	Enabled         bool   `xml:"enabled,attr" hcl:"enabled"`
	ContinueOnError bool   `xml:"continueOnError,attr,omitempty" hcl:"continue_on_error"`
	Async           bool   `xml:"async,attr,omitempty" hcl:"async"`
}

func LoadCommonPolicyHCL(item *ast.ObjectItem, p *Policy) error {

	if err := hcl.DecodeObject(p, item.Val.(*ast.ObjectType)); err != nil {
		return err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("policy not an object")
	}

	if enabledList := listVal.Filter("enabled"); len(enabledList.Items) == 0 {
		p.Enabled = true
	}

	p.Name = item.Keys[1].Token.Value().(string)

	return nil
}
