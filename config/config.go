package config

import (
	//	"fmt"
	//	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apg-hcl/config/policy"
)

type Config struct {
	Proxy           *Proxy
	ProxyEndpoints  []*ProxyEndpoint
	TargetEndpoints []*TargetEndpoint
	Policies        []interface{}
}

func LoadConfigFromHCL(list *ast.ObjectList) (*Config, error) {
	var config Config

	if proxies := list.Filter("proxy"); len(proxies.Items) > 0 {
		result, err := loadProxyHCL(proxies)

		if err != nil {
			return nil, err
		}

		config.Proxy = result
	}

	if proxyEndpoints := list.Filter("proxy_endpoint"); len(proxyEndpoints.Items) > 0 {
		result, err := loadProxyEndpointsHCL(proxyEndpoints)
		if err != nil {
			return nil, err
		}

		config.ProxyEndpoints = result
	}

	if targetEndpoints := list.Filter("target_endpoint"); len(targetEndpoints.Items) > 0 {
		result, err := loadTargetEndpointsHCL(targetEndpoints)
		if err != nil {
			return nil, err
		}

		config.TargetEndpoints = result
	}

	if policies := list.Filter("policy"); len(policies.Items) > 0 {
		for _, item := range policies.Items {
			policyType := item.Keys[0].Token.Value().(string)

			if policyType == "assign_message" {
				_, err := policy.LoadAssignMessageHCL(item)

				if err != nil {
					return nil, err
				}

			}
		}
	}
	return &config, nil
}
