package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apg-hcl/config/common"
)

type JavaScriptPolicy struct {
	XMLName     string `xml:"Javascript" hcl:"-"`
	Policy      `hcl:",squash"`
	TimeLimit   int                `xml:"timeLimit,attr" hcl:"time_limit"`
	DisplayName string             `xml:",omitempty" hcl:"display_name"`
	ResourceURL string             `hcl:"resource_url"`
	IncludeURL  string             `xml:",omitempty" hcl:"include_url"`
	Properties  []*common.Property `xml:"Properties>Property" hcl:"properties"`
	Content     string             `xml:"-" hcl:"content"`
}

func LoadJavaScriptHCL(item *ast.ObjectItem) (interface{}, error) {
	var p JavaScriptPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
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
		props, err := common.LoadPropertiesHCL(propsList.Items[0])
		if err != nil {
			return nil, err
		}

		p.Properties = props
	}

	return p, nil
}
