package endpoints

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// PreFlow represents a <PreFlow/> element for
// ProxyEndpoint and TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#flows
type PreFlow struct {
	XMLName  string       `xml:"PreFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

// Flow represents a <Flow/> element for
// ProxyEndpoint and TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#flows
type Flow struct {
	XMLName   string       `xml:"Flow" hcl:"-"`
	Name      string       `xml:"name,attr" hcl:"-"`
	Condition string       `xml:",omitempty" hcl:"condition"`
	Request   FlowRequest  `hcl:"request"`
	Response  FlowResponse `hcl:"response"`
}

// PostFlow represents a <PostFlow/> element for
// ProxyEndpoint and TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#flows
type PostFlow struct {
	XMLName  string       `xml:"PostFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

// PostClientFlow represents a <PostClientFlow/> element for
// ProxyEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#flows
type PostClientFlow struct {
	XMLName  string       `xml:"PostClientFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

// FaultRule represents a <FaultRule/> element for
// ProxyEndpoint and TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/content/fault-handling
type FaultRule struct {
	XMLName   string      `xml:"FaultRule" hcl:"-"`
	Name      string      `xml:"name,attr" hcl:"-"`
	Condition string      `xml:",omitempty" hcl:"condition"`
	Steps     []*FlowStep `xml:",innerxml" hcl:"step"`
}

// DefaultFaultRule represents a <DefaultFaultRule/> element for
// ProxyEndpoint and TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/content/fault-handling
type DefaultFaultRule struct {
	XMLName       string      `xml:"DefaultFaultRule" hcl:"-"`
	Name          string      `xml:"name,attr" hcl:"-"`
	Condition     string      `xml:",omitempty" hcl:"condition"`
	Steps         []*FlowStep `xml:",innerxml" hcl:"step"`
	AlwaysEnforce bool        `xml:",omitempty" hcl:"always_enforce"`
}

// FlowStep represents a <Step/> element for
// Request and Response flows in ProxyEndpoint and
// TargetEndpoint definitions.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#policies-policyattachment
type FlowStep struct {
	XMLName   string `xml:"Step"`
	Name      string
	Condition string `xml:",omitempty" hcl:"condition"`
}

// FlowRequest represents a <Request/> element for
// PreFlow, Flow, PostFlow, and PostClientFlow definitions
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#watchaquickhowtovideo-flowconfigurationelements
type FlowRequest struct {
	XMLName string      `xml:"Request" hcl:"-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

// FlowResponse represents a <Response/> element for
// PreFlow, Flow, PostFlow, and PostClientFlow definitions
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#watchaquickhowtovideo-flowconfigurationelements
type FlowResponse struct {
	XMLName string      `xml:"Response" hcl:"-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

func decodePreFlowHCL(list *ast.ObjectList) (*PreFlow, error) {
	var result PreFlow
	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("pre flow item not an object")
	}

	if request := listVal.Filter("request"); len(request.Items) > 0 {
		item := request.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func decodeFlowsHCL(list *ast.ObjectList) ([]*Flow, error) {
	var result []*Flow

	for _, item := range list.Items {
		var flow Flow
		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("pre flow item not an object")
		}

		if request := listVal.Filter("request"); len(request.Items) > 0 {
			item := request.Items[0]

			steps, err := decodeFlowStepsHCL(item)
			if err != nil {
				return nil, err
			}

			flow.Request.Steps = steps
		}

		if response := listVal.Filter("response"); len(response.Items) > 0 {
			item := response.Items[0]

			steps, err := decodeFlowStepsHCL(item)
			if err != nil {
				return nil, err
			}

			flow.Response.Steps = steps
		}

		flow.Name = item.Keys[0].Token.Value().(string)
		result = append(result, &flow)
	}

	return result, nil
}

func decodePostFlowHCL(list *ast.ObjectList) (*PostFlow, error) {
	var result PostFlow
	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("post flow item not an object")
	}

	if request := listVal.Filter("request"); len(request.Items) > 0 {
		item := request.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func decodePostClientFlowHCL(list *ast.ObjectList) (*PostClientFlow, error) {
	var result PostClientFlow
	item := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("post client flow item not an object")
	}

	if request := listVal.Filter("request"); len(request.Items) > 0 {
		item := request.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func decodeFaultRulesHCL(list *ast.ObjectList) ([]*FaultRule, error) {
	var result []*FaultRule

	for _, item := range list.Items {
		var faultRule FaultRule

		steps, err := decodeFlowStepsHCL(item)
		if err != nil {
			return nil, err
		}

		faultRule.Steps = steps

		faultRule.Name = item.Keys[0].Token.Value().(string)
		result = append(result, &faultRule)
	}

	return result, nil
}

func decodeDefaultFaultRuleHCL(item *ast.ObjectItem) (*DefaultFaultRule, error) {
	var faultRule DefaultFaultRule

	if err := hcl.DecodeObject(&faultRule, item.Val); err != nil {
		return nil, fmt.Errorf("error decoding step object")
	}

	steps, err := decodeFlowStepsHCL(item)
	if err != nil {
		return nil, err
	}

	faultRule.Steps = steps

	faultRule.Name = item.Keys[0].Token.Value().(string)

	return &faultRule, nil
}

func decodeFlowStepsHCL(list *ast.ObjectItem) ([]*FlowStep, error) {
	var listVal *ast.ObjectList
	if ot, ok := list.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("request item not an object")
	}

	var flowSteps []*FlowStep
	if steps := listVal.Filter("step"); len(steps.Items) > 0 {
		for _, step := range steps.Items {
			var s FlowStep
			if err := hcl.DecodeObject(&s, step.Val); err != nil {
				return nil, fmt.Errorf("error decoding step object")
			}
			s.Name = step.Keys[0].Token.Value().(string)

			flowSteps = append(flowSteps, &s)
		}
	}

	return flowSteps, nil
}
