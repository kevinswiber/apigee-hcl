package policy

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type AssignMessagePolicy struct {
	Policy   `hcl:",squash"`
	AssignTo AssignMessageAssignTo `xml:",omitempty" hcl:"assign_to"`
}

type AssignMessageAssignTo struct {
	CreateNew bool   `xml:"createNew,attr" hcl:"create_new"`
	Transport string `xml:"transport,attr" hcl:"transport"`
	Type      string `xml:"type,attr" hcl:"type"`
}

func LoadAssignMessageHCL(item *ast.ObjectItem) (interface{}, error) {
	var p AssignMessagePolicy

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	p.Name = item.Keys[1].Token.Value().(string)

	return &p, nil
}
