package policy

import (
	"fmt"

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

func LoadAssignMessageHCL(item *ast.ObjectItem) (*AssignMessagePolicy, error) {
	var p AssignMessagePolicy

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	fmt.Printf("assign message: %+v\n", p)
	return &p, nil
}
