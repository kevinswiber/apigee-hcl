package properties

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

// Property represents the <Property/> element.
type Property struct {
	XMLName string `xml:"Property" hcl:"-" hcle:"omit"`
	Name    string `xml:"name,attr" hcl:",key"`
	Value   string `xml:",chardata" hcl:"-"`
}

// DecodeHCL converts an ast.ObjectItem into Property objects.
func DecodeHCL(item *ast.ObjectItem) ([]*Property, error) {
	var propsVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		propsVal = ot.List
	} else {
		return nil, fmt.Errorf("error decoding properties")
	}

	var newProps []*Property
	for _, p := range propsVal.Items {
		var val *ast.LiteralType
		if lt, ok := p.Val.(*ast.LiteralType); ok {
			val = lt
		}

		var newProp Property
		switch val.Token.Type {
		case token.NUMBER, token.FLOAT, token.BOOL:
			newProp = Property{Name: p.Keys[0].Token.Value().(string), Value: val.Token.Text}
		default:
			{
				var v string
				if err := hcl.DecodeObject(&v, p.Val); err != nil {
					return nil, err
				}

				newProp = Property{Name: p.Keys[0].Token.Value().(string), Value: v}
			}
		}
		newProps = append(newProps, &newProp)
	}

	return newProps, nil
}
