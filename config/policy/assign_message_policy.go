package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type AssignMessagePolicy struct {
	XMLName     string `xml:"AssignMessage" hcl:"-"`
	Policy      `hcl:",squash"`
	DisplayName string   `xml:",omitempty" hcl:"display_name"`
	Add         add      `xml:",omitempty" hcl:"add"`
	AssignTo    assignTo `xml:",omitempty" hcl:"assign_to"`
}

type add struct {
	XMLName string    `xml:"Add" hcl:"-"`
	Headers []*header `xml:"Headers>Header" hcl:"header"`
}

type header struct {
	XMLName string `xml:"Header" hcl:"-"`
	Name    string `xml:"name,attr" hcl"-"`
	Value   string `xml:",omitempty" hcl:"value"`
}

type assignTo struct {
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

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("assign message policy not an object")
	}

	if addList := listVal.Filter("add"); len(addList.Items) > 0 {
		headers, err := loadHeadersHCL(addList.Items[0])
		if err != nil {
			return nil, err
		}
		p.Add.Headers = headers
	}

	return &p, nil
}

func loadHeadersHCL(item *ast.ObjectItem) ([]*header, error) {
	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("header parent not an object")
	}

	var headers []*header
	if headerList := listVal.Filter("header"); len(headerList.Items) > 0 {
		for _, h := range headerList.Items {
			var hdr header
			if err := hcl.DecodeObject(&hdr, h.Val.(*ast.ObjectType)); err != nil {
				return nil, err
			}
			hdr.Name = h.Keys[0].Token.Value().(string)

			headers = append(headers, &hdr)
		}
	}

	return headers, nil
}
