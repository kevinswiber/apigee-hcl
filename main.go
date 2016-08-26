package main

import (
	//"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	//"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"github.com/kevinswiber/apigee-hcl/config"
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

	config, err := config.LoadConfigFromHCL(list)
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	//fmt.Printf("%#v\n", config.Proxy)
	var output []byte
	if output, err = xml.MarshalIndent(config.Proxy, "", "    "); err != nil {
		log.Fatalf("err: %s", err)
	}
	fmt.Printf("%s\n", string(output))

	for _, proxyEndpoint := range config.ProxyEndpoints {
		//fmt.Printf("%#v\n", proxyEndpoint)
		var output []byte
		output, err := xml.MarshalIndent(proxyEndpoint, "", "    ")
		if err != nil {
			log.Fatalf("err: %s", err)
		}
		fmt.Printf("%s\n", string(output))
	}
}
