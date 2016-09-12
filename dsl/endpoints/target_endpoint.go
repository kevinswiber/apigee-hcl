package endpoints

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	hclEncoding "github.com/kevinswiber/apigee-hcl/dsl/encoding/hcl"
	"github.com/kevinswiber/apigee-hcl/dsl/properties"
)

// TargetEndpoint represents a <TargetEndpoint/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint
type TargetEndpoint struct {
	XMLName               string                 `xml:"TargetEndpoint" hcl:"-" hcle:"omit"`
	Name                  string                 `xml:"name,attr" hcl:",key"`
	PreFlow               *PreFlow               `hcl:"pre_flow"`
	Flows                 []*Flow                `xml:"Flows,omitempty>Flow" hcl:"flow"`
	PostFlow              *PostFlow              `hcl:"post_flow"`
	FaultRules            []*FaultRule           `xml:"FaultRules,omitempty>FaultRule" hcl:"fault_rule"`
	DefaultFaultRule      []*DefaultFaultRule    `hcl:"default_fault_rule"`
	HTTPTargetConnection  *HTTPTargetConnection  `hcl:"http_target_connection"`
	LocalTargetConnection *LocalTargetConnection `xml:",omitempty" hcl:"local_target_connection" hcle:"omitempty"`
	ScriptTarget          *ScriptTarget          `xml:",omitempty" hcl:"script_target" hcle:"omitempty"`
	SSLInfo               *SSLInfo               `xml:",omitempty" hcl:"ssl_info" hcle:"omitempty"`
}

// HTTPTargetConnection represents an <HTTPTargetConnection/> element
// in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type HTTPTargetConnection struct {
	XMLName      string                 `xml:"HTTPTargetConnection" hcl:"-" hcle:"omit"`
	URL          string                 `hcl:"url"`
	LoadBalancer *LoadBalancer          `hcl:"load_balancer"`
	Properties   []*properties.Property `xml:"Properties>Property" hcl:"properties"`
}

// LoadBalancer represents a <LoadBalancer/> element in an
// HTTPTargetConnection.
//
// Documentation: http://docs.apigee.com/api-platform/content/load-balance-api-traffic-across-multiple-backend-servers#configuringatargetendpointtoloadbalanceacrossnamedtargetservers
type LoadBalancer struct {
	XMLName      string                `xml:"LoadBalancer" hcl:"-" hcle:"omit"`
	Algorithm    string                `hcl:"algorithm"`
	Servers      []*LoadBalancerServer `xml:"Server" hcl:"server"`
	MaxFailures  int                   `xml:",omitempty" hcl:"max_failures" hcle:"omitempty"`
	RetryEnabled bool                  `xml:",omitempty" hcl:"retry_enabled" hcle:"omitempty"`
}

// LocalTargetConnection represents a <LocalTargetConnection/> element
// in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type LocalTargetConnection struct {
	XMLName       string `xml:"LocalTargetConnection" hcl:"-" hcle:"omit"`
	APIProxy      string `xml:",omitempty" hcl:"api_proxy" hcle:"omitempty"`
	ProxyEndpoint string `xml:",omitempty" hcl:"proxy_endpoint" hcle:"omitempty"`
	Path          string `xml:",omitempty" hcl:"path" hcle:"omitempty"`
}

// ScriptTarget represents a <ScriptTarget/> element in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type ScriptTarget struct {
	XMLName              string                 `xml:"ScriptTarget" hcl:"-" hcle:"omit"`
	ResourceURL          string                 `hcl:"resource_url"`
	EnvironmentVariables []*EnvironmentVariable `xml:"EnvironmentVariables>EnvironmentVariable" hcl:"environment_variables"`
	Arguments            []string               `xml:"Arguments>Argument" hcl:"arguments"`
}

// SSLInfo represents an <SSLInfo/> element in a TargetEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#tlsssltargetendpointconfiguration-tlsssltargetendpointconfigurationelements
type SSLInfo struct {
	XMLName           string   `xml:"SSLInfo" hcl:"-" hcle:"omit"`
	Enabled           bool     `xml:",omitempty" hcl:"enabled" hcle:"omitempty"`
	TrustStore        string   `xml:",omitempty" hcl:"trust_store" hcle:"omitempty"`
	ClientAuthEnabled bool     `xml:",omitempty" hcl:"client_auth_enabled" hcle:"omitempty"`
	KeyStore          string   `xml:",omitempty" hcl:"key_store" hcle:"omitempty"`
	KeyAlias          string   `xml:",omitempty" hcl:"key_alias" hcle:"omitempty"`
	Ciphers           []string `xml:"Ciphers>Cipher" hcl:"ciphers" hcle:"omitempty"`
	Protocols         []string `xml:"Protocols>Protocol" hcl:"protocols" hcle:"omitempty"`
}

// EnvironmentVariable represents an <EnvironmentVariable/> element
// in a ScriptTarget.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#targetendpoint-targetendpointconfigurationelements
type EnvironmentVariable struct {
	XMLName string `xml:"EnvironmentVariable" hcl:"-" hcle:"omit"`
	Name    string `xml:"name,attr" hcl:",key"`
	Value   string `xml:",chardata" hcl:"-"`
}

// LoadBalancerServer represents a <LoadBalancerServer/> element
// in a LoadBalancer.
//
// Documentation: http://docs.apigee.com/api-platform/content/load-balance-api-traffic-across-multiple-backend-servers#configuringatargetendpointtoloadbalanceacrossnamedtargetservers
type LoadBalancerServer struct {
	XMLName    string `xml:"Server" hcl:"-" hcle:"omit"`
	Name       string `xml:"name,attr" hcl:",key"`
	Weight     int    `xml:",omitempty" hcl:"weight" hcle:"omitempty"`
	IsFallback bool   `xml:",omitempty" hcl:"is_fallback" hcle:"omitempty"`
}

// DecodeTargetEndpointsHCL converts an HCL ast.ObjectItem into a TargetEndpoint object.
func DecodeTargetEndpointsHCL(list *ast.ObjectList) ([]*TargetEndpoint, error) {
	var errors *multierror.Error

	var targetEndpoints []*TargetEndpoint
	if err := hcl.DecodeObject(&targetEndpoints, list); err != nil {
		return nil, err
	}

	for i, item := range list.Items {
		if len(item.Keys) == 0 || item.Keys[0].Token.Value() == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("target endpoint requires a name"),
			}

			errors = multierror.Append(errors, &newError)
			continue
		}

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("http proxy connection not an object")
		}

		if htcList := listVal.Filter("http_target_connection"); len(htcList.Items) > 0 {
			htc, err := DecodeHTTPTargetConnectionHCL(htcList.Items[0])
			if err != nil {
				errors = multierror.Append(errors, err)
			} else {
				targetEndpoints[i].HTTPTargetConnection = htc
			}
		}

		if scriptTargetList := listVal.Filter("script_target"); len(scriptTargetList.Items) > 0 {
			st, err := decodeTargetEndpointScriptTargetHCL(scriptTargetList.Items[0])
			if err != nil {
				errors = multierror.Append(errors, err)
			} else {
				targetEndpoints[i].ScriptTarget = st
			}
		}
	}

	if errors != nil {
		return nil, errors
	}

	return targetEndpoints, nil
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
		var val *ast.LiteralType
		if lt, ok := p.Val.(*ast.LiteralType); ok {
			val = lt
		}

		var newEnv EnvironmentVariable
		switch val.Token.Type {
		case token.NUMBER, token.FLOAT, token.BOOL:
			newEnv = EnvironmentVariable{Name: p.Keys[0].Token.Value().(string), Value: val.Token.Text}
		default:
			{
				var v string
				if err := hcl.DecodeObject(&v, p.Val); err != nil {
					return nil, err
				}

				newEnv = EnvironmentVariable{Name: p.Keys[0].Token.Value().(string), Value: v}
			}
		}

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

	return &htc, nil
}
