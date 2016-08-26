package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

type Config struct {
	Proxy           *Proxy
	ProxyEndpoints  []*ProxyEndpoint
	TargetEndpoints []*TargetEndpoint
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

	return &config, nil
}
