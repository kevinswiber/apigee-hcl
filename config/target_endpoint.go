package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type TargetEndpoint struct {
	XMLName              string                `xml:"TargetEndpoint" hcl:"-"`
	Name                 string                `xml:"name,attr" hcl:"-"`
	PreFlow              *PreFlow              `hcl:"pre_flow"`
	Flows                []*Flow               `xml:"Flows,omitempty>Flow" hcl:"flows"`
	PostFlow             *PostFlow             `hcl:"post_flow"`
	FaultRules           []*FaultRule          `xml:"FaultRules,omitempty>FaultRule" hcl:"fault_rules"`
	DefaultFaultRule     *DefaultFaultRule     `hcl:"default_fault_rule"`
	HTTPTargetConnection *HTTPTargetConnection `hcl:"http_target_connection"`
}

type HTTPTargetConnection struct {
	XMLName      string        `xml:"HTTPTargetConnection" hcl:"-"`
	URL          string        `hcl:"url"`
	LoadBalancer *LoadBalancer `hcl:"load_balancer"`
	Properties   []*Property   `xml:"Properties,omitempty>Property" hcl:"properties"`
}

type LoadBalancer struct {
	XMLName      string                `xml:"LoadBalancer" hcl:"-"`
	Algorithm    string                `hcl:"algorithm"`
	Servers      []*LoadBalancerServer `xml:"Server" hcl:"server"`
	MaxFailures  int                   `xml:",omitempty" hcl:"max_failures"`
	RetryEnabled bool                  `xml:",omitempty" hcl:"retry_enabled"`
}

type LoadBalancerServer struct {
	XMLName    string `xml:"Server" hcl:"-"`
	Name       string `xml:"name,attr" hcl:"-"`
	Weight     int    `xml:",omitempty" hcl:"weight"`
	IsFallback bool   `xml:",omitempty" hcl:"is_fallback"`
}

func loadTargetEndpointsHCL(list *ast.ObjectList) ([]*TargetEndpoint, error) {
	var result []*TargetEndpoint
	for _, item := range list.Items {
		n := item.Keys[0].Token.Value().(string)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("target endpoint not an object")
		}

		var targetEndpoint TargetEndpoint

		if err := hcl.DecodeObject(&targetEndpoint, item.Val); err != nil {
			return nil, err
		}

		targetEndpoint.Name = n

		if preFlow := listVal.Filter("pre_flow"); len(preFlow.Items) > 0 {
			preFlow, err := loadPreFlowHCL(preFlow)
			if err != nil {
				return nil, err
			}

			targetEndpoint.PreFlow = preFlow
		}

		if flows := listVal.Filter("flow"); len(flows.Items) > 0 {
			flows, err := loadFlowsHCL(flows)
			if err != nil {
				return nil, err
			}

			targetEndpoint.Flows = flows
		}

		if postFlow := listVal.Filter("post_flow"); len(postFlow.Items) > 0 {
			postFlow, err := loadPostFlowHCL(postFlow)
			if err != nil {
				return nil, err
			}

			targetEndpoint.PostFlow = postFlow
		}

		if faultRulesList := listVal.Filter("fault_rule"); len(faultRulesList.Items) > 0 {
			faultRules, err := loadFaultRulesHCL(faultRulesList)
			if err != nil {
				return nil, err
			}

			targetEndpoint.FaultRules = faultRules
		}

		if defaultFaultRulesList := listVal.Filter("default_fault_rule"); len(defaultFaultRulesList.Items) > 0 {
			faultRule, err := loadDefaultFaultRuleHCL(defaultFaultRulesList.Items[0])
			if err != nil {
				return nil, err
			}

			targetEndpoint.DefaultFaultRule = faultRule
		}

		if htcList := listVal.Filter("http_target_connection"); len(htcList.Items) > 0 {
			htc, err := loadTargetEndpointHTTPTargetConnectionHCL(htcList.Items[0])
			if err != nil {
				return nil, err
			}

			targetEndpoint.HTTPTargetConnection = htc
		}

		result = append(result, &targetEndpoint)
	}

	return result, nil
}

func loadTargetEndpointHTTPTargetConnectionHCL(item *ast.ObjectItem) (*HTTPTargetConnection, error) {
	var htc HTTPTargetConnection

	if err := hcl.DecodeObject(&htc, item.Val); err != nil {
		return nil, fmt.Errorf("error decoding http target connection")
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("http proxy connection not an object")
	}

	if propsList := listVal.Filter("properties"); len(propsList.Items) > 0 {
		props, err := loadPropertiesHCL(propsList.Items[0])
		if err != nil {
			return nil, err
		}

		htc.Properties = props
	}

	if lbList := listVal.Filter("load_balancer"); len(lbList.Items) > 0 {
		var lb LoadBalancer
		if err := hcl.DecodeObject(&lb, lbList.Items[0]); err != nil {
			return nil, err
		}

		var lbListVal *ast.ObjectList
		if ot, ok := lbList.Items[0].Val.(*ast.ObjectType); ok {
			lbListVal = ot.List
		} else {
			return nil, fmt.Errorf("load balancer not an object")
		}

		var lbServers []*LoadBalancerServer
		if serversList := lbListVal.Filter("server"); len(serversList.Items) > 0 {
			fmt.Println("got servers")
			for _, item := range serversList.Items {
				var s LoadBalancerServer
				if err := hcl.DecodeObject(&s, item); err != nil {
					return nil, err
				}
				s.Name = item.Keys[0].Token.Value().(string)
				fmt.Printf("name = %s\n", s.Name)
				lbServers = append(lbServers, &s)
			}

			lb.Servers = lbServers
		}

		htc.LoadBalancer = &lb
	}

	return &htc, nil
}
