package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Policy Represents a base Policy element. Each policy type should embed a Policy.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#policies
type Policy struct {
	InternalName    string `xml:"name,attr,omitempty" hcl:"-"`
	Enabled         bool   `xml:"enabled,attr" hcl:"enabled"`
	ContinueOnError bool   `xml:"continueOnError,attr,omitempty" hcl:"continue_on_error"`
	Async           bool   `xml:"async,attr,omitempty" hcl:"async"`
}

// Namer is used to set and retrieve a policy name
type Namer interface {
	Name() string
	SetName(string)
}

// Name returns the name of the policy.
func (p *Policy) Name() string {
	return p.InternalName
}

// SetName sets the name of the policy.
func (p *Policy) SetName(name string) {
	p.InternalName = name
}

// Resourcer is used for policies with resources
type Resourcer interface {
	Resource() *Resource
}

// Resource represents an included file in a proxy bundle
type Resource struct {
	URL     string
	Content string
}

// DecodeHCL converts an HCL ast.ObjectItem into a Policy object.
func DecodeHCL(item *ast.ObjectItem, p *Policy) error {
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

	p.SetName(item.Keys[1].Token.Value().(string))

	return nil
}
