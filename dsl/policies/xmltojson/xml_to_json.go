package xmltojson

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/hclerror"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// XMLToJSON represents an <XMLToJSON/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/xml-json-policy
type XMLToJSON struct {
	XMLName        string `xml:"XMLToJSON" hcl:"-"`
	policy.Policy  `hcl:",squash"`
	DisplayName    string          `xml:",omitempty" hcl:"display_name"`
	Source         string          `xml:",omitempty" hcl:"source"`
	OutputVariable string          `xml:",omitempty" hcl:"output_variable"`
	Options        *xmlJSONOptions `xml:",omitempty" hcl:"options"`
	Format         string          `xml:",omitempty" hcl:"format"`
}

type xmlJSONOptions struct {
	XMLName                  string               `xml:"Options" hcl:"-"`
	RecognizeNumber          bool                 `xml:",omitempty" hcl:"recognize_number"`
	RecognizeBoolean         bool                 `xml:",omitempty" hcl:"recognize_boolean"`
	RecognizeNull            bool                 `xml:",omitempty" hcl:"recognize_null"`
	NullValue                string               `xml:",omitempty" hcl:"null_value"`
	NamespaceBlockName       string               `xml:",omitempty" hcl:"namespace_block_name"`
	DefaultNamespaceNodeName string               `xml:",omitempty" hcl:"default_namespace_node_name"`
	NamespaceSeparator       string               `xml:",omitempty" hcl:"namespace_separator"`
	TextAlwaysAsProperty     bool                 `xml:",omitempty" hcl:"text_always_as_property"`
	TextNodeName             string               `xml:",omitempty" hcl:"text_node_name"`
	AttributeBlockName       string               `xml:",omitempty" hcl:"attribute_block_name"`
	OutputPrefix             string               `xml:",omitempty" hcl:"output_prefix"`
	OutputSuffix             string               `xml:",omitempty" hcl:"output_suffix"`
	StripLevels              int                  `xml:",omitempty" hcl:"strip_levels"`
	TreatAsArray             *xmlJSONTreatAsArray `xml:",omitempty" hcl:"treat_as_array"`
}

type xmlJSONTreatAsArray struct {
	XMLName string          `xml:"TreatAsArray" hcl:"-"`
	Paths   *[]*xmlJSONPath `xml:"Path,omitempty" hcl:"path"`
}

type xmlJSONPath struct {
	XMLName string `xml:"Path" hcl:"-"`
	Unwrap  bool   `xml:"unwrap,attr,omitempty" hcl:"unwrap"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeHCL converts an HCL ast.ObjectItem into a XMLToJSON object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p XMLToJSON

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if p.Options != nil && p.Format != "" {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("xml_to_json must specify either options or format, not both"),
		}
		return nil, &newError
	}

	if p.Options == nil && p.Format == "" {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("xml_to_json must specify either options or format"),
		}
		return nil, &newError
	}

	if p.Options != nil && p.Options.NamespaceBlockName != "" {
		if p.Options.DefaultNamespaceNodeName == "" ||
			p.Options.NamespaceSeparator == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("xml_to_json must specify default_namespace_node_name and " +
					"namespace_separator when namespace_block_name is used"),
			}
			return nil, &newError
		}
	}

	return &p, nil
}
