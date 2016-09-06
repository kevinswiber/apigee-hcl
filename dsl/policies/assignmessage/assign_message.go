package assignmessage

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// AssignMessage represents an <AssignMessage/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/assign-message-policy
type AssignMessage struct {
	XMLName                   string `xml:"AssignMessage" hcl:"-"`
	policy.Policy             `hcl:",squash"`
	DisplayName               string          `xml:",omitempty" hcl:"display_name"`
	Add                       *Add            `xml:",omitempty" hcl:"add"`
	Copy                      *Copy           `xml:",omitempty" hcl:"copy"`
	Remove                    *Remove         `xml:",omitempty" hcl:"remove"`
	Set                       *Set            `xml:",omitempty" hcl:"set"`
	AssignVariable            *assignVariable `xml:",omitempty" hcl:"assign_variable"`
	AssignTo                  *assignTo       `xml:",omitempty" hcl:"assign_to"`
	IgnoreUnresolvedVariables bool            `xml:",omitempty" hcl:"ignore_unresolved_variables"`
}

type assignVariable struct {
	Name  string `hcl:"name"`
	Ref   string `xml:",omitempty" hcl:"ref"`
	Value string `xml:",omitempty" hcl:"value"`
}

type assignTo struct {
	CreateNew bool   `xml:"createNew,attr" hcl:"create_new"`
	Transport string `xml:"transport,attr,omitempty" hcl:"transport"`
	Type      string `xml:"type,attr,omitempty" hcl:"type"`
	Value     string `xml:",chardata" hcl:"value"`
}

// DecodeHCL converts an HCL ast.ObjectItem into an AssignMessage object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p AssignMessage

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		errors = multierror.Append(errors, fmt.Errorf("assign message policy not an object"))
		return nil, errors
	}

	if addList := listVal.Filter("add"); len(addList.Items) > 0 {
		item := addList.Items[0]
		a, err := DecodeAddHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Add = a
		}
	} else {
		p.Add = nil
	}

	if copyList := listVal.Filter("copy"); len(copyList.Items) > 0 {
		item := copyList.Items[0]
		a, err := DecodeCopyHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Copy = a
		}
	} else {
		p.Copy = nil
	}

	if removeList := listVal.Filter("remove"); len(removeList.Items) > 0 {
		item := removeList.Items[0]
		a, err := DecodeRemoveHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Remove = a
		}
	} else {
		p.Remove = nil
	}

	if setList := listVal.Filter("set"); len(setList.Items) > 0 {
		item := setList.Items[0]
		a, err := DecodeSetHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Set = a
		}
	} else {
		p.Set = nil
	}

	if errors != nil {
		return nil, errors
	}

	return &p, nil
}
