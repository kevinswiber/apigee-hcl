package raisefault

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/assignmessage"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// RaiseFault represents a <RaiseFault/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/raise-fault-policy
type RaiseFault struct {
	XMLName                   string `xml:"RaiseFault" hcl:"-"`
	policy.Policy             `hcl:",squash"`
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
	XMLName      string                   `xml:"Copy" hcl:"-"`
	Source       string                   `xml:"source,attr,omitempty" hcl:"-"`
	Headers      *[]*assignmessage.Header `xml:"Headers>Header" hcl:"header"`
	StatusCode   bool                     `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase bool                     `xml:",omitempty" hcl:"reason_phrase"`
}

type raiseFaultRemove struct {
	XMLName string                   `xml:"Remove" hcl:"-"`
	Headers *[]*assignmessage.Header `xml:"Headers>Header" hcl:"header"`
}

type raiseFaultSet struct {
	XMLName      string                   `xml:"Set" hcl:"-"`
	Headers      *[]*assignmessage.Header `xml:"Headers>Header" hcl:"header"`
	Payload      *assignmessage.Payload   `xml:",omitempty" hcl:"payload"`
	StatusCode   int                      `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase string                   `xml:",omitempty" hcl:"reason_phrase"`
}

// DecodeHCL converts an HCL ast.ObjectItem into a RaiseFault.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p RaiseFault

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
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
		a, err := decodeFaultResponse(item)
		if err != nil {
			return nil, err
		}
		p.FaultResponse = a
	} else {
		p.FaultResponse = nil
	}

	return &p, nil
}

func decodeFaultResponse(item *ast.ObjectItem) (*faultResponse, error) {
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
		a, err := decodeCopyHCL(item)
		if err != nil {
			return nil, err
		}

		result.Copy = a
	} else {
		result.Copy = nil
	}

	if removeList := listVal.Filter("remove"); len(removeList.Items) > 0 {
		item := removeList.Items[0]
		a, err := decodeRemoveHCL(item)
		if err != nil {
			return nil, err
		}

		result.Remove = a
	} else {
		result.Remove = nil
	}

	if setList := listVal.Filter("set"); len(setList.Items) > 0 {
		item := setList.Items[0]
		a, err := decodeSetHCL(item)
		if err != nil {
			return nil, err
		}

		result.Set = a
	} else {
		result.Set = nil
	}

	return result, nil
}

func decodeCopyHCL(item *ast.ObjectItem) (*raiseFaultCopy, error) {
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

	headers, err := assignmessage.DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}

func decodeRemoveHCL(item *ast.ObjectItem) (*raiseFaultRemove, error) {
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

	headers, err := assignmessage.DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}

func decodeSetHCL(item *ast.ObjectItem) (*raiseFaultSet, error) {
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

	headers, err := assignmessage.DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	return &result, nil
}
