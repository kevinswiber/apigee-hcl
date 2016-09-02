package config

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config/hclerror"
)

// Proxy represents an <APIProxy/> element in an Apigee proxy bundle
//
// Documentation: http://docs.apigee.com/api-services/reference/api-proxy-configuration-reference#baseconfig
type Proxy struct {
	XMLName     string `xml:"APIProxy", hcl:"-"`
	Name        string `xml:"name,attr,omitempty" hcl:"-"`
	DisplayName string `xml:",omitempty" hcl:"display_name"`
	Description string `xml:",omitempty" hcl:"description"`
}

func loadProxyHCL(list *ast.ObjectList) (*Proxy, error) {
	var errors *multierror.Error

	var item = list.Items[0]
	if len(item.Keys) == 0 || item.Keys[0].Token.Value() == "" {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("proxy requires a name"),
		}
		errors = multierror.Append(errors, &newError)
		return nil, errors
	}

	n := item.Keys[0].Token.Value().(string)

	if _, ok := item.Val.(*ast.ObjectType); !ok {
		errors = multierror.Append(errors, fmt.Errorf("proxy not an object"))
		return nil, errors
	}

	var proxy Proxy
	if err := hcl.DecodeObject(&proxy, item.Val); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	proxy.Name = n

	return &proxy, nil
}
