package endpoints

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hclEncoding "github.com/kevinswiber/apigee-hcl/dsl/encoding/hcl"
	"github.com/kevinswiber/apigee-hcl/dsl/properties"
)

// ProxyEndpoint represents a <ProxyEndpoint/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#proxyendpoint
type ProxyEndpoint struct {
	XMLName             string               `xml:"ProxyEndpoint" hcl:"-" hcle:"omit"`
	Name                string               `xml:"name,attr" hcl:",key"`
	PreFlow             *PreFlow             `hcl:"pre_flow"`
	Flows               []*Flow              `xml:"Flows>Flow" hcl:"flow"`
	PostFlow            *PostFlow            `hcl:"post_flow"`
	PostClientFlow      *PostClientFlow      `hcl:"post_client_flow"`
	FaultRules          []*FaultRule         `xml:"FaultRules>FaultRule" hcl:"fault_rule"`
	DefaultFaultRule    []*DefaultFaultRule  `hcl:"default_fault_rule"`
	HTTPProxyConnection *HTTPProxyConnection `hcl:"http_proxy_connection"`
	RouteRules          []*RouteRule         `xml:"RouteRule" hcl:"route_rule"`
}

// HTTPProxyConnection represents an <HTTPProxyConnection/> element
// in a ProxyEndpoint.
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#proxyendpoint-proxyendpointconfigurationelements
type HTTPProxyConnection struct {
	XMLName      string                 `xml:"HTTPProxyConnection" hcl:"-" hcle:"omit"`
	BasePath     string                 `hcl:"base_path"`
	VirtualHosts []string               `xml:"VirtualHost" hcl:"virtual_host"`
	Properties   []*properties.Property `xml:"Properties>Property" hcl:"properties"`
}

// RouteRule represents a <RouteRule/> element in an HTTPProxyConnection
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#proxyendpoint-proxyendpointconfigurationelements
type RouteRule struct {
	XMLName        string `xml:"RouteRule" hcl:"-" hcle:"omit"`
	Name           string `xml:"name,attr" hcl:",key"`
	Condition      string `xml:",omitempty" hcl:"condition" hcle:"omitempty"`
	TargetEndpoint string `xml:",omitempty" hcl:"target_endpoint" hcle:"omitempty"`
	URL            string `xml:",omitempty" hcl:"url" hcle:"omitempty"`
}

// DecodeProxyEndpointHCL converts an HCL ast.ObjectItem into a ProxyEndpoint.
func DecodeProxyEndpointHCL(item *ast.ObjectItem) (*ProxyEndpoint, error) {
	var errors *multierror.Error
	if len(item.Keys) == 0 || item.Keys[0].Token.Value() == "" {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("proxy endpoint requires a name"),
		}

		errors = multierror.Append(errors, &newError)
		return nil, errors
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("proxy endpoint is not an object"),
		}

		errors = multierror.Append(errors, &newError)
		return nil, errors
	}

	var proxyEndpoint ProxyEndpoint

	if err := hcl.DecodeObject(&proxyEndpoint, item); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	if hpcList := listVal.Filter("http_proxy_connection"); len(hpcList.Items) > 0 {
		hpc, err := decodeHTTPProxyConnectionHCL(hpcList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			proxyEndpoint.HTTPProxyConnection = hpc
		}
	}

	if errors != nil {
		return nil, errors
	}

	return &proxyEndpoint, nil
}

func decodeHTTPProxyConnectionHCL(item *ast.ObjectItem) (*HTTPProxyConnection, error) {
	var hpc HTTPProxyConnection

	if err := hcl.DecodeObject(&hpc, item); err != nil {
		return nil, fmt.Errorf("error decoding http proxy connection")
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("http proxy connection not an object")
	}

	if propsList := listVal.Filter("properties"); len(propsList.Items) > 0 {
		props, err := properties.DecodeHCL(propsList.Items[0])
		if err != nil {
			return nil, err
		}

		hpc.Properties = props
	}

	return &hpc, nil
}
