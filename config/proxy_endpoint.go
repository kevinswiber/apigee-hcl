package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type ProxyEndpoint struct {
	XMLName             string               `xml:"ProxyEndpoint" hcl:",-"`
	Name                string               `xml:"name,attr" hcl:",-"`
	PreFlow             *PreFlow             `hcl:"pre_flow"`
	Flows               *Flows               `hcl:"flows"`
	PostFlow            *PostFlow            `hcl:"post_flow"`
	HTTPProxyConnection *HTTPProxyConnection `hcl:"http_proxy_connection"`
	RouteRules          []*RouteRule         `xml:"RouteRule" hcl:"route_rule"`
}

type PreFlow struct {
	XMLName  string       `xml:"PreFlow" hcl:",-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type Flows struct {
	XMLName  string       `xml:"Flows" hcl:",-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type PostFlow struct {
	XMLName  string       `xml:"PostFlow" hcl:",-"`
	Request  FlowRequest  `hcl:"request"`
	Response FlowResponse `hcl:"response"`
}

type HTTPProxyConnection struct {
	XMLName     string   `xml:"HTTPProxyConnection", hcl:",-"`
	BasePath    string   `hcl:"base_path"`
	VirtualHost []string `xml:",innerxml" hcl:"virtual_host"`
}

type RouteRule struct {
	XMLName        string `xml:"RouteRule"`
	Name           string `xml:",attr", hcl:",-"`
	Condition      string `xml:",omitempty" hcl:"condition"`
	TargetEndpoint string `xml:",omitempty" hcl:"target_endpoint"`
	URL            string `xml:",omitempty" hcl:"url"`
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
			preFlow, err := loadProxyEndpointPreFlowHcl(preFlow)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.PreFlow = preFlow
		}

		result = append(result, &proxyEndpoint)
	}

	return result, nil
}

func loadProxyEndpointPreFlowHcl(list *ast.ObjectList) (*PreFlow, error) {
	var result PreFlow
	preFlowItem := list.Items[0]

	var listVal *ast.ObjectList
	if ot, ok := preFlowItem.Val.(*ast.ObjectType); ok {
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
