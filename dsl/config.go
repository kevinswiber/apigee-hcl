package dsl

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/endpoints"
	"github.com/kevinswiber/apigee-hcl/dsl/hclerror"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// Config is a container for holding the contents of an exported Apigee proxy bundle
type Config struct {
	Proxy           *Proxy
	ProxyEndpoints  []*endpoints.ProxyEndpoint
	TargetEndpoints []*endpoints.TargetEndpoint
	Policies        []policy.Namer
	Resources       map[string]string
}

// DecodeConfigHCL converts an HCL ast.ObjectList into a Config object
func DecodeConfigHCL(list *ast.ObjectList) (*Config, error) {
	var errors *multierror.Error

	var c Config

	if proxies := list.Filter("proxy"); len(proxies.Items) > 0 {
		result, err := decodeProxyHCL(proxies)
		if err != nil {
			errors = multierror.Append(errors, err)
			return nil, errors
		}

		c.Proxy = result
	}

	if proxyEndpoints := list.Filter("proxy_endpoint"); len(proxyEndpoints.Items) > 0 {
		var result []*endpoints.ProxyEndpoint
		for _, item := range proxyEndpoints.Items {
			proxyEndpoint, err := endpoints.DecodeProxyEndpointHCL(item)
			if err != nil {
				errors = multierror.Append(errors, err)
				return nil, errors
			}
			result = append(result, proxyEndpoint)
		}

		c.ProxyEndpoints = result
	}

	if targetEndpoints := list.Filter("target_endpoint"); len(targetEndpoints.Items) > 0 {
		var result []*endpoints.TargetEndpoint
		for _, item := range targetEndpoints.Items {
			targetEndpoint, err := endpoints.DecodeTargetEndpointHCL(item)
			if err != nil {
				errors = multierror.Append(errors, err)
				return nil, errors
			}
			result = append(result, targetEndpoint)
		}

		c.TargetEndpoints = result
	}

	if policyList := list.Filter("policy"); len(policyList.Items) > 0 {
		var ps []policy.Namer

		for _, item := range policyList.Items {
			if len(item.Keys) < 2 ||
				item.Keys[0].Token.Value() == "" ||
				item.Keys[1].Token.Value() == "" {
				pos := item.Val.Pos()
				newError := hclerror.PosError{
					Pos: pos,
					Err: fmt.Errorf("policy requires a type and name"),
				}

				errors = multierror.Append(errors, &newError)
				continue
			}
			policyType := item.Keys[0].Token.Value().(string)

			if f, ok := PolicyList[policyType]; ok {
				p, err := f(item)
				if err != nil {
					errors = multierror.Append(errors, err)
				}

				switch p.(type) {
				case policy.Resourcer:
					resourcePolicy := p.(policy.Resourcer)
					r := resourcePolicy.Resource()
					if len(r.URL) > 0 && len(r.Content) > 0 {
						if c.Resources == nil {
							c.Resources = make(map[string]string)
						}
						c.Resources[r.URL] = r.Content
					}
				}
				ps = append(ps, p.(policy.Namer))
			}
		}

		if errors != nil {
			return nil, errors
		}

		c.Policies = ps
	}
	return &c, nil
}
