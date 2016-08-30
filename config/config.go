package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

type Config struct {
	Proxy           *Proxy
	ProxyEndpoints  []*ProxyEndpoint
	TargetEndpoints []*TargetEndpoint
	Policies        []interface{}
}

func LoadConfigFromHCL(list *ast.ObjectList) (*Config, error) {
	var c Config

	if proxies := list.Filter("proxy"); len(proxies.Items) > 0 {
		result, err := loadProxyHCL(proxies)

		if err != nil {
			return nil, err
		}

		c.Proxy = result
	}

	if proxyEndpoints := list.Filter("proxy_endpoint"); len(proxyEndpoints.Items) > 0 {
		result, err := loadProxyEndpointsHCL(proxyEndpoints)
		if err != nil {
			return nil, err
		}

		c.ProxyEndpoints = result
	}

	if targetEndpoints := list.Filter("target_endpoint"); len(targetEndpoints.Items) > 0 {
		result, err := loadTargetEndpointsHCL(targetEndpoints)
		if err != nil {
			return nil, err
		}

		c.TargetEndpoints = result
	}

	if policies := list.Filter("policy"); len(policies.Items) > 0 {
		var ps []interface{}

		for _, item := range policies.Items {
			policyType := item.Keys[0].Token.Value().(string)

			if f, ok := PolicyList[policyType]; ok {
				p, err := f(item)
				if err != nil {
					return nil, err
				}

				ps = append(ps, p)
			}
		}

		c.Policies = ps
	}
	return &c, nil
}
