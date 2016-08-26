package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type PreFlow struct {
	XMLName  string       `xml:"PreFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type Flow struct {
	XMLName   string       `xml:"Flow" hcl:"-"`
	Name      string       `xml:"name,attr" hcl:"-"`
	Condition string       `xml:",omitempty" hcl:"condition"`
	Request   FlowRequest  `hcl:"request"`
	Response  FlowResponse `hcl:"response"`
}

type PostFlow struct {
	XMLName  string       `xml:"PostFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type PostClientFlow struct {
	XMLName  string       `xml:"PostClientFlow" hcl:"-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type FlowStep struct {
	XMLName   string `xml:"Step"`
	Name      string
	Condition string `xml:",omitempty" hcl:"condition"`
}

type FlowRequest struct {
	XMLName string      `xml:"Request" hcl:"-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

type FlowResponse struct {
	XMLName string      `xml:"Response" hcl:"-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

func loadPreFlowHCL(list *ast.ObjectList) (*PreFlow, error) {
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

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func loadFlowsHCL(list *ast.ObjectList) ([]*Flow, error) {
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

			steps, err := loadFlowSteps(item)
			if err != nil {
				return nil, err
			}

			flow.Request.Steps = steps
		}

		if response := listVal.Filter("response"); len(response.Items) > 0 {
			item := response.Items[0]

			steps, err := loadFlowSteps(item)
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

func loadPostFlowHCL(list *ast.ObjectList) (*PostFlow, error) {
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

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func loadPostClientFlowHCL(list *ast.ObjectList) (*PostClientFlow, error) {
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

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Request.Steps = steps
	}

	if response := listVal.Filter("response"); len(response.Items) > 0 {
		item := response.Items[0]

		steps, err := loadFlowSteps(item)
		if err != nil {
			return nil, err
		}

		result.Response.Steps = steps
	}

	return &result, nil
}

func loadFlowSteps(list *ast.ObjectItem) ([]*FlowStep, error) {
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
