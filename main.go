package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apg-hcl/config"
)

func main() {
	file := "./examples/helloworld.hcl"
	//file := "./examples/conditional_policy.hcl"

	d, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	hclRoot, err := hcl.Parse(string(d))
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	list, ok := hclRoot.Node.(*ast.ObjectList)
	if !ok {
		log.Fatalf("error parsing: file doesn't contain root object")
	}

	c, err := config.LoadConfigFromHCL(list)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	var output []byte
	if output, err = xml.MarshalIndent(c.Proxy, "", "    "); err != nil {
		log.Fatalf("err: %s", err)
	}
	fmt.Printf("%s\n", string(output))

	for _, proxyEndpoint := range c.ProxyEndpoints {
		var output []byte
		output, err := xml.MarshalIndent(proxyEndpoint, "", "    ")
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		fmt.Printf("%s\n", string(output))
	}

	for _, targetEndpoint := range c.TargetEndpoints {
		var output []byte
		output, err := xml.MarshalIndent(targetEndpoint, "", "    ")
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		fmt.Printf("%s\n", string(output))
	}

	for _, policy := range c.Policies {
		var output []byte
		output, err := xml.MarshalIndent(policy, "", "    ")
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		fmt.Printf("%s\n", string(output))
	}
}
