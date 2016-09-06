package verifyapikey

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// VerifyAPIKey represents a <VerifyAPIKey/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/verify-api-key-policy
type VerifyAPIKey struct {
	XMLName       string `xml:"VerifyAPIKey" hcl:"-"`
	policy.Policy `hcl:",squash"`
	DisplayName   string  `xml:",omitempty" hcl:"display_name"`
	APIKey        *apikey `hcl:"apikey"`
}

type apikey struct {
	XMLName string `xml:"APIKey" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeHCL converts an HCL ast.ObjectItem into a VerifyAPIKey object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p VerifyAPIKey

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
