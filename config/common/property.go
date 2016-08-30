package common

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Property struct {
	XMLName string      `xml:"Property"`
	Name    string      `xml:"name,attr" hcl:",key"`
	Value   interface{} `xml:",chardata" hcl:"-"`
}

func LoadPropertiesHCL(item *ast.ObjectItem) ([]*Property, error) {
	var propsVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		propsVal = ot.List
	} else {
		return nil, fmt.Errorf("error decoding properties")
	}

	var newProps []*Property
	for _, p := range propsVal.Items {
		var val interface{}
		if err := hcl.DecodeObject(&val, p.Val); err != nil {
			return nil, fmt.Errorf("can't decode property object")
		}

		newProp := Property{Name: p.Keys[0].Token.Value().(string), Value: val}
		newProps = append(newProps, &newProp)
	}

	return newProps, nil
}
