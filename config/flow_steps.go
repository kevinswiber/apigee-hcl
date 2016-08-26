package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type FlowStep struct {
	XMLName   string `xml:"Step"`
	Name      string
	Condition string `xml:",omitempty" hcl:"condition"`
}

type FlowRequest struct {
	XMLName string      `xml:"Request" hcl:",-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

type FlowResponse struct {
	XMLName string      `xml:"Response" hcl:",-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
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
			stepName := step.Keys[0].Token.Value().(string)
			var s FlowStep
			if err := hcl.DecodeObject(&s, step.Val); err != nil {
				return nil, fmt.Errorf("error decoding step object")
			}
			s.Name = stepName

			flowSteps = append(flowSteps, &s)
		}
	}

	return flowSteps, nil
}
