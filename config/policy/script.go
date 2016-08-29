package policy

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type ScriptPolicy struct {
	XMLName     string `xml:"Script" hcl:"-"`
	Policy      `hcl:",squash"`
	DisplayName string `xml:",omitempty" hcl:"display_name"`
	ResourceURL string `hcl:"resource_url"`
	IncludeURL  string `xml:",omitempty" hcl:"include_url"`
}

func LoadScriptHCL(item *ast.ObjectItem) (interface{}, error) {
	var p ScriptPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
