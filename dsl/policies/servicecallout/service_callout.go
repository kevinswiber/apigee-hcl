package servicecallout

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/endpoints"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/assignmessage"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// ServiceCallout represents an <ServiceCallout/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/service-callout-policy
type ServiceCallout struct {
	XMLName               string `xml:"ServiceCallout" hcl:"-"`
	policy.Policy         `hcl:",squash"`
	DisplayName           string                           `xml:",omitempty" hcl:"display_name"`
	Request               *scRequest                       `hcl:"request"`
	HTTPTargetConnection  *endpoints.HTTPTargetConnection  `hcl:"http_target_connection"`
	LocalTargetConnection *endpoints.LocalTargetConnection `hcl:"local_target_connection"`
	Response              string                           `xml:",omitempty" hcl:"response"`
	Timeout               int                              `xml:",omitempty" hcl:"timeout"`
}

type scRequest struct {
	XMLName                   string                `xml:"Request" hcl:"-"`
	ClearPayload              bool                  `xml:"clearPayload,attr,omitempty" hcl:"clear_payload"`
	Variable                  string                `xml:"variable,attr,omitempty" hcl:"variable"`
	Add                       *assignmessage.Add    `xml:",omitempty" hcl:"add"`
	Copy                      *assignmessage.Copy   `xml:",omitempty" hcl:"copy"`
	Remove                    *assignmessage.Remove `xml:",omitempty" hcl:"remove"`
	Set                       *assignmessage.Set    `xml:",omitempty" hcl:"set"`
	IgnoreUnresolvedVariables bool                  `xml:",omitempty" hcl:"ignore_unresolved_variables"`
}

// DecodeHCL converts an HCL ast.ObjectItem into an ServiceCallout object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p ServiceCallout

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
		errors = multierror.Append(errors, fmt.Errorf("service_callout policy not an object"))
		return nil, errors
	}

	if requestList := listVal.Filter("request"); len(requestList.Items) > 0 {
		item := requestList.Items[0]
		r, err := decodeRequestHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Request = r
		}
	} else {
		p.Request = nil
	}

	if htcList := listVal.Filter("http_target_connection"); len(htcList.Items) > 0 {
		htc, err := endpoints.DecodeHTTPTargetConnectionHCL(htcList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.HTTPTargetConnection = htc
		}
	}

	if errors != nil {
		return nil, errors
	}

	return &p, nil
}

func decodeRequestHCL(item *ast.ObjectItem) (*scRequest, error) {
	var r scRequest
	var errors *multierror.Error
	var listVal *ast.ObjectList

	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		errors = multierror.Append(errors, fmt.Errorf("service_callout policy not an object"))
		return nil, errors
	}

	if err := hcl.DecodeObject(&r, item.Val.(*ast.ObjectType)); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	if addList := listVal.Filter("add"); len(addList.Items) > 0 {
		item := addList.Items[0]
		a, err := assignmessage.DecodeAddHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			r.Add = a
		}
	} else {
		r.Add = nil
	}

	if copyList := listVal.Filter("copy"); len(copyList.Items) > 0 {
		item := copyList.Items[0]
		a, err := assignmessage.DecodeCopyHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			r.Copy = a
		}
	} else {
		r.Copy = nil
	}

	if removeList := listVal.Filter("remove"); len(removeList.Items) > 0 {
		item := removeList.Items[0]
		a, err := assignmessage.DecodeRemoveHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			r.Remove = a
		}
	} else {
		r.Remove = nil
	}

	if setList := listVal.Filter("set"); len(setList.Items) > 0 {
		item := setList.Items[0]
		a, err := assignmessage.DecodeSetHCL(item)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			r.Set = a
		}
	} else {
		r.Set = nil
	}

	if errors != nil {
		return nil, errors
	}

	return &r, nil
}
