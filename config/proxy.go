package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Proxy struct {
	XMLName     string `xml:"APIProxy", hcl:",-"`
	Name        string `xml:"name,attr,omitempty" hcl:",-"`
	DisplayName string `xml:",omitempty" hcl:"display_name"`
	Description string `xml:",omitempty" hcl:"description"`
}

func loadProxyHCL(list *ast.ObjectList) (*Proxy, error) {
	//TODO: Check if more than one proxy.  Report error.
	var item = list.Items[0]
	n := item.Keys[0].Token.Value().(string)

	if _, ok := item.Val.(*ast.ObjectType); !ok {
		return nil, fmt.Errorf("proxy not an object")
	}

	var proxy Proxy
	if err := hcl.DecodeObject(&proxy, item.Val); err != nil {
		return nil, err
	}

	proxy.Name = n

	return &proxy, nil
}
