package policy

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type VerifyAPIKeyPolicy struct {
	XMLName     string `xml:"VerifyAPIKey" hcl:"-"`
	Policy      `hcl:",squash"`
	DisplayName string   `xml:",omitempty" hcl:"display_name"`
  APIKey      string    *apikey `hcl:"apikey"`
}

type apikey struct {
	XMLName string `xml:"APIKey" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

func LoadVerifyAPIKeyHCL(item *ast.ObjectItem) (interface{}, error) {
	var p VerifyAPIKeyPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return p, nil
}
