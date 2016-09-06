package endpoints

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/hclerror"
	"github.com/kevinswiber/apigee-hcl/dsl/properties"
)

// TargetEndpoint represents a <TargetEndpoint/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint
type TargetEndpoint struct {
	XMLName               string                 `xml:"TargetEndpoint" hcl:"-"`
	Name                  string                 `xml:"name,attr" hcl:"-"`
	PreFlow               *PreFlow               `hcl:"pre_flow"`
	Flows                 []*Flow                `xml:"Flows,omitempty>Flow" hcl:"flows"`
	PostFlow              *PostFlow              `hcl:"post_flow"`
	FaultRules            []*FaultRule           `xml:"FaultRules,omitempty>FaultRule" hcl:"fault_rules"`
	DefaultFaultRule      *DefaultFaultRule      `hcl:"default_fault_rule"`
	HTTPTargetConnection  *HTTPTargetConnection  `hcl:"http_target_connection"`
	LocalTargetConnection *LocalTargetConnection `xml:",omitempty" hcl:"local_target_connection"`
	ScriptTarget          *ScriptTarget          `xml:",omitempty" hcl:"script_target"`
	SSLInfo               *SSLInfo               `xml:",omitempty" hcl:"ssl_info"`
}

// HTTPTargetConnection represents an <HTTPTargetConnection/> element
// in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type HTTPTargetConnection struct {
	XMLName      string                 `xml:"HTTPTargetConnection" hcl:"-"`
	URL          string                 `hcl:"url"`
	LoadBalancer *LoadBalancer          `hcl:"load_balancer"`
	Properties   []*properties.Property `xml:"Properties>Property" hcl:"properties"`
}

// LoadBalancer represents a <LoadBalancer/> element in an
// HTTPTargetConnection.
//
// Documentation: http://docs.apigee.com/api-platform/content/load-balance-api-traffic-across-multiple-backend-servers#configuringatargetendpointtoloadbalanceacrossnamedtargetservers
type LoadBalancer struct {
	XMLName      string                `xml:"LoadBalancer" hcl:"-"`
	Algorithm    string                `hcl:"algorithm"`
	Servers      []*LoadBalancerServer `xml:"Server" hcl:"server"`
	MaxFailures  int                   `xml:",omitempty" hcl:"max_failures"`
	RetryEnabled bool                  `xml:",omitempty" hcl:"retry_enabled"`
}

// LocalTargetConnection represents a <LocalTargetConnection/> element
// in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type LocalTargetConnection struct {
	XMLName       string `xml:"LocalTargetConnection" hcl:"-"`
	APIProxy      string `xml:",omitempty" hcl:"api_proxy"`
	ProxyEndpoint string `xml:",omitempty" hcl:"proxy_endpoint"`
	Path          string `xml:",omitempty" hcl:"path"`
}

// ScriptTarget represents a <ScriptTarget/> element in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type ScriptTarget struct {
	XMLName              string                 `xml:"ScriptTarget" hcl:"-"`
	ResourceURL          string                 `hcl:"resource_url"`
	EnvironmentVariables []*EnvironmentVariable `xml:"EnvironmentVariables>EnvironmentVariable" hcl:"environment_variables"`
	Arguments            []string               `xml:"Arguments>Argument" hcl:"arguments"`
}

// SSLInfo represents an <SSLInfo/> element in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#tlsssltargetendpointconfiguration-tlsssltargetendpointconfigurationelements
type SSLInfo struct {
	XMLName           string   `xml:"SSLInfo" hcl:"-"`
	Enabled           bool     `xml:",omitempty" hcl:"enabled"`
	TrustStore        string   `xml:",omitempty" hcl:"trust_store"`
	ClientAuthEnabled bool     `xml:",omitempty" hcl:"client_auth_enabled"`
	KeyStore          string   `xml:",omitempty" hcl:"key_store"`
	KeyAlias          string   `xml:",omitempty" hcl:"key_alias"`
	Ciphers           []string `xml:"Ciphers>Cipher" hcl:"ciphers"`
	Protocols         []string `xml:"Protocols>Protocol" hcl:"protocols"`
}

// EnvironmentVariable represents an <EnvironmentVariable/> element
// in a ScriptTarget.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type EnvironmentVariable struct {
	XMLName string      `xml:"EnvironmentVariable" hcl:"-"`
	Name    string      `xml:"name,attr" hcl:",key"`
	Value   interface{} `xml:",chardata" hcl:"-"`
}

// LoadBalancerServer represents a <LoadBalancerServer/> element
// in a LoadBalancer.
//
// Documentation: http://docs.apigee.com/api-platform/content/load-balance-api-traffic-across-multiple-backend-servers#configuringatargetendpointtoloadbalanceacrossnamedtargetservers
type LoadBalancerServer struct {
	XMLName    string `xml:"Server" hcl:"-"`
	Name       string `xml:"name,attr" hcl:"-"`
	Weight     int    `xml:",omitempty" hcl:"weight"`
	IsFallback bool   `xml:",omitempty" hcl:"is_fallback"`
}

// DecodeTargetEndpointHCL converts an HCL ast.ObjectItem into a TargetEndpoint object.
func DecodeTargetEndpointHCL(item *ast.ObjectItem) (*TargetEndpoint, error) {
	var errors *multierror.Error
	if len(item.Keys) == 0 || item.Keys[0].Token.Value() == "" {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("target endpoint requires a name"),
		}

		errors = multierror.Append(errors, &newError)
		return nil, errors
	}
	n := item.Keys[0].Token.Value().(string)

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		errors = multierror.Append(errors, fmt.Errorf("target endpoint is not an object"))
		return nil, errors
	}

	var targetEndpoint TargetEndpoint

	if err := hcl.DecodeObject(&targetEndpoint, item.Val); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	targetEndpoint.Name = n

	if preFlow := listVal.Filter("pre_flow"); len(preFlow.Items) > 0 {
		preFlow, err := decodePreFlowHCL(preFlow)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.PreFlow = preFlow
		}
	}

	if flows := listVal.Filter("flow"); len(flows.Items) > 0 {
		flows, err := decodeFlowsHCL(flows)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.Flows = flows
		}
	}

	if postFlow := listVal.Filter("post_flow"); len(postFlow.Items) > 0 {
		postFlow, err := decodePostFlowHCL(postFlow)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.PostFlow = postFlow
		}
	}

	if faultRulesList := listVal.Filter("fault_rule"); len(faultRulesList.Items) > 0 {
		faultRules, err := decodeFaultRulesHCL(faultRulesList)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.FaultRules = faultRules
		}
	}

	if defaultFaultRulesList := listVal.Filter("default_fault_rule"); len(defaultFaultRulesList.Items) > 0 {
		faultRule, err := decodeDefaultFaultRuleHCL(defaultFaultRulesList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.DefaultFaultRule = faultRule
		}
	}

	if htcList := listVal.Filter("http_target_connection"); len(htcList.Items) > 0 {
		htc, err := DecodeHTTPTargetConnectionHCL(htcList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.HTTPTargetConnection = htc
		}
	}

	if scriptTargetList := listVal.Filter("script_target"); len(scriptTargetList.Items) > 0 {
		st, err := decodeTargetEndpointScriptTargetHCL(scriptTargetList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			targetEndpoint.ScriptTarget = st
		}
	}

	if errors != nil {
		return nil, errors
	}

	return &targetEndpoint, nil
}

func decodeTargetEndpointScriptTargetHCL(item *ast.ObjectItem) (*ScriptTarget, error) {
	var st ScriptTarget

	if err := hcl.DecodeObject(&st, item.Val); err != nil {
		return nil, fmt.Errorf("error decoding http target connection")
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("http proxy connection not an object")
	}

	if envsList := listVal.Filter("environment_variables"); len(envsList.Items) > 0 {
		envs, err := decodeTargetEndpointScriptTargetEnvironmentVariablesHCL(envsList.Items[0])
		if err != nil {
			return nil, err
		}

		st.EnvironmentVariables = envs
	}

	return &st, nil
}

func decodeTargetEndpointScriptTargetEnvironmentVariablesHCL(item *ast.ObjectItem) ([]*EnvironmentVariable, error) {
	var envsVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		envsVal = ot.List
	} else {
		return nil, fmt.Errorf("error decoding enverties")
	}

	var newEnvs []*EnvironmentVariable
	for _, p := range envsVal.Items {
		var val interface{}
		if err := hcl.DecodeObject(&val, p.Val); err != nil {
			return nil, fmt.Errorf("can't decode environment variable object")
		}

		newEnv := EnvironmentVariable{Name: p.Keys[0].Token.Value().(string), Value: val}
		newEnvs = append(newEnvs, &newEnv)
	}

	return newEnvs, nil
}

// DecodeHTTPTargetConnectionHCL converts an HCL ast.ObjectItem into an HTTPTargetConnection object.
func DecodeHTTPTargetConnectionHCL(item *ast.ObjectItem) (*HTTPTargetConnection, error) {
	var htc HTTPTargetConnection

	if err := hcl.DecodeObject(&htc, item.Val); err != nil {
		return nil, fmt.Errorf("error decoding http target connection")
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("http target connection not an object")
	}

	if propsList := listVal.Filter("properties"); len(propsList.Items) > 0 {
		props, err := properties.DecodeHCL(propsList.Items[0])
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
			for _, item := range serversList.Items {
				var s LoadBalancerServer
				if err := hcl.DecodeObject(&s, item); err != nil {
					return nil, err
				}
				s.Name = item.Keys[0].Token.Value().(string)
				lbServers = append(lbServers, &s)
			}

			lb.Servers = lbServers
		}

		htc.LoadBalancer = &lb
	}

	return &htc, nil
}
