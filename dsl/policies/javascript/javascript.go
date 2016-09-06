package javascript

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
	"github.com/kevinswiber/apigee-hcl/dsl/properties"
)

// JavaScript represents a <Javascript/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/javascript-policy
type JavaScript struct {
	XMLName       string `xml:"Javascript" hcl:"-"`
	policy.Policy `hcl:",squash"`
	TimeLimit     int                    `xml:"timeLimit,attr" hcl:"time_limit"`
	DisplayName   string                 `xml:",omitempty" hcl:"display_name"`
	ResourceURL   string                 `hcl:"resource_url"`
	IncludeURL    []string               `xml:",omitempty" hcl:"include_url"`
	Properties    []*properties.Property `xml:"Properties>Property" hcl:"properties"`
	Content       string                 `xml:"-" hcl:"content"`
}

// Resource represents an included file in a proxy bundle
func (j *JavaScript) Resource() *policy.Resource {
	return &policy.Resource{
		URL:     j.ResourceURL,
		Content: j.Content,
	}
}

// DecodeHCL converts HCL into an JavaScript object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p JavaScript

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("javascript policy not an object")
	}

	if propsList := listVal.Filter("properties"); len(propsList.Items) > 0 {
		props, err := properties.DecodeHCL(propsList.Items[0])
		if err != nil {
			return nil, err
		}

		p.Properties = props
	}

	return &p, nil
}
