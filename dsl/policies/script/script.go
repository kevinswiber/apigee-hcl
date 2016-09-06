package script

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// Script represents a <Script/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/python-script-policy
type Script struct {
	XMLName       string `xml:"Script" hcl:"-"`
	policy.Policy `hcl:",squash"`
	DisplayName   string   `xml:",omitempty" hcl:"display_name"`
	ResourceURL   string   `hcl:"resource_url"`
	IncludeURL    []string `xml:",omitempty" hcl:"include_url"`
	Content       string   `xml:"-" hcl:"content"`
}

// Resource represents an included file in a proxy bundle
func (s *Script) Resource() *policy.Resource {
	return &policy.Resource{
		URL:     s.ResourceURL,
		Content: s.Content,
	}
}

// DecodeHCL converts an HCL ast.ObjectItem into a Script object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p Script

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
