package policy

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// RaiseFaultPolicy represents a <RaiseFaultPolicy/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/raise-fault-policy
type RaiseFaultPolicy struct {
	XMLName                   string `xml:"RaiseFault" hcl:"-"`
	Policy                    `hcl:",squash"`
	DisplayName               string         `xml:",omitempty" hcl:"display_name"`
	FaultResponse             *faultResponse `xml:"FaultResponse" hcl:"fault_response"`
	IgnoreUnresolvedVariables bool           `xml:"IgnoreUnresolvedVariables" hcl:"ignore_unresolved_variables"`
}

type faultResponse struct {
	Copy   *raiseFaultCopy   `xml:",omitempty" hcl:"copy"`
	Remove *raiseFaultRemove `xml:",omitempty" hcl:"remove"`
	Set    *raiseFaultSet    `xml:",omitempty" hcl:"set"`
}

type raiseFaultCopy struct {
	XMLName      string    `xml:"Copy" hcl:"-"`
	Source       string    `xml:"source,attr,omitempty" hcl:"-"`
	Headers      []*header `xml:"Headers>Header" hcl:"header"`
	StatusCode   bool      `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase bool      `xml:",omitempty" hcl:"reason_phrase"`
}

type raiseFaultRemove struct {
	XMLName string    `xml:"Remove" hcl:"-"`
	Headers []*header `xml:"Headers>Header" hcl:"header"`
}

type raiseFaultSet struct {
	XMLName      string    `xml:"Set" hcl:"-"`
	Headers      []*header `xml:"Headers>Header" hcl:"header"`
	Payload      *payload  `xml:",omitempty" hcl:"payload"`
	StatusCode   int       `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase string    `xml:",omitempty" hcl:"reason_phrase"`
}

// LoadRaiseFaultHCL converts an HCL ast.ObjectItem into a RaiseFaultPolicy.
func LoadRaiseFaultHCL(item *ast.ObjectItem) (interface{}, error) {
	var p RaiseFaultPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("raise fault policy not an object")
	}

	if ignoreUnresolved := listVal.Filter("ignore_unresolved_variables"); len(ignoreUnresolved.Items) == 0 {
		p.IgnoreUnresolvedVariables = true
	}

	if faultResponseList := listVal.Filter("fault_response"); len(faultResponseList.Items) > 0 {
		item := faultResponseList.Items[0]
		a, err := loadRaiseFaultFaultResponse(item)
		if err != nil {
			return nil, err
		}
		p.FaultResponse = a
	} else {
		p.FaultResponse = nil
	}

	return p, nil
}

func loadRaiseFaultFaultResponse(item *ast.ObjectItem) (*faultResponse, error) {
	var result *faultResponse

	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList

	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("fault response not an object")
	}

	if copyList := listVal.Filter("copy"); len(copyList.Items) > 0 {
		item := copyList.Items[0]
		a, err := loadRaiseFaultCopyHCL(item)
		if err != nil {
			return nil, err
		}

		result.Copy = a
	} else {
		result.Copy = nil
	}

	if removeList := listVal.Filter("remove"); len(removeList.Items) > 0 {
		item := removeList.Items[0]
		a, err := loadRaiseFaultRemoveHCL(item)
		if err != nil {
			return nil, err
		}

		result.Remove = a
	} else {
		result.Remove = nil
	}

	if setList := listVal.Filter("set"); len(setList.Items) > 0 {
		item := setList.Items[0]
		a, err := loadRaiseFaultSetHCL(item)
		if err != nil {
			return nil, err
		}

		result.Set = a
	} else {
		result.Set = nil
	}

	return result, nil
}

func loadRaiseFaultCopyHCL(item *ast.ObjectItem) (*raiseFaultCopy, error) {
	var result raiseFaultCopy
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("copy not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}

func loadRaiseFaultRemoveHCL(item *ast.ObjectItem) (*raiseFaultRemove, error) {
	var result raiseFaultRemove
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("remove not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}

func loadRaiseFaultSetHCL(item *ast.ObjectItem) (*raiseFaultSet, error) {
	var result raiseFaultSet
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("set not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}
