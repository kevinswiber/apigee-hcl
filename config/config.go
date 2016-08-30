package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apg-hcl/config/policy"
)

type Config struct {
	Proxy           *Proxy
	ProxyEndpoints  []*ProxyEndpoint
	TargetEndpoints []*TargetEndpoint
	Policies        []interface{}
	Resources       map[string]string
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

				switch p.(type) {
				case policy.ScriptPolicy:
					script := p.(policy.ScriptPolicy)
					if len(script.ResourceURL) > 0 && len(script.Content) > 0 {
						if c.Resources == nil {
							c.Resources = make(map[string]string)
						}
						c.Resources[script.ResourceURL] = script.Content
					}
				case policy.JavaScriptPolicy:
					script := p.(policy.JavaScriptPolicy)
					if len(script.ResourceURL) > 0 && len(script.Content) > 0 {
						if c.Resources == nil {
							c.Resources = make(map[string]string)
						}
						c.Resources[script.ResourceURL] = script.Content
					}
				}
				ps = append(ps, p)
			}
		}

		c.Policies = ps
	}
	return &c, nil
}
