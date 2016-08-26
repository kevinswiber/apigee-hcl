package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Config struct {
	Proxy          *Proxy
	ProxyEndpoints []*ProxyEndpoint
}

type ProxyEndpoint struct {
	XMLName string  `xml:"ProxyEndpoint" hcl:",-"`
	Name    string  `xml:"name,attr" hcl:",-"`
	PreFlow PreFlow `hcl:"pre_flow"`
	//HTTPProxyConnection interface{} `hcl:"http_proxy_connection"`
	//RouteRule           interface{} `hcl:"route_rule"`
}

type PreFlow struct {
	XMLName  string       `xml:"PreFlow" hcl:",-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type PostFlow struct {
	XMLName  string       `xml:"PostFlow" hcl:",-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type FlowRequest struct {
	XMLName string      `xml:"Request" hcl:",-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

type FlowResponse struct {
	XMLName string      `xml:"Response" hcl:",-"`
	Steps   []*FlowStep `xml:",innerxml" hcl:"step"`
}

type FlowStep struct {
	XMLName   string `xml:"Step"`
	Name      string
	Condition string `xml:",omitempty" hcl:"condition"`
}

func LoadConfigFromHCL(list *ast.ObjectList) (*Config, error) {
	var config Config

	if proxies := list.Filter("proxy"); len(proxies.Items) > 0 {
		result, err := loadProxyHcl(proxies)

		if err != nil {
			return nil, err
		}

		config.Proxy = result
	}

	if proxyEndpoints := list.Filter("proxy_endpoint"); len(proxyEndpoints.Items) > 0 {
		result, err := loadProxyEndpointsHcl(proxyEndpoints)
		if err != nil {
			return nil, err
		}

		config.ProxyEndpoints = result
	}

	return &config, nil
}

func loadProxyEndpointsHcl(list *ast.ObjectList) ([]*ProxyEndpoint, error) {
	var result []*ProxyEndpoint
	for _, item := range list.Items {
		n := item.Keys[0].Token.Value().(string)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("proxy endpoint not an object")
		}

		var proxyEndpoint ProxyEndpoint

		if err := hcl.DecodeObject(&proxyEndpoint, item.Val); err != nil {
			return nil, err
		}

		proxyEndpoint.Name = n

		if preFlow := listVal.Filter("pre_flow"); len(preFlow.Items) > 0 {
			proxyEndpoint.PreFlow = PreFlow{}
			preFlowItem := preFlow.Items[0]

			var preFlowItemVal *ast.ObjectList
			if preFlowOT, ok := preFlowItem.Val.(*ast.ObjectType); ok {
				preFlowItemVal = preFlowOT.List
			} else {
				return nil, fmt.Errorf("pre flow item not an object")
			}

			if request := preFlowItemVal.Filter("request"); len(request.Items) > 0 {
				item := request.Items[0]

				steps, err := loadFlowSteps(item)
				if err != nil {
					return nil, err
				}

				proxyEndpoint.PreFlow.Request.Steps = steps
			}

			if response := preFlowItemVal.Filter("response"); len(response.Items) > 0 {
				item := response.Items[0]

				steps, err := loadFlowSteps(item)
				if err != nil {
					return nil, err
				}

				proxyEndpoint.PreFlow.Response.Steps = steps
			}

		}

		result = append(result, &proxyEndpoint)
	}

	return result, nil
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
