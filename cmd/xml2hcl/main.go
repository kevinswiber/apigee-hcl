package main

import (
	"encoding/xml"
	"fmt"
	"github.com/kevinswiber/apigee-hcl/dsl/endpoints"
	"github.com/rodaine/hclencoder"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	cli := &cli{
		Args:   os.Args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	os.Exit(cli.Run())
}

type cli struct {
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *cli) Run() int {
	if 2 != len(c.Args) {
		return c.Usage()
	}

	file := c.Args[1]

	var d []byte
	var err error

	if "-h" == file {
		return c.Usage()
	} else if "-" == file {
		d, err = ioutil.ReadAll(c.Stdin)
	} else {
		d, err = ioutil.ReadFile(file)
	}

	if err != nil {
		log.Fatalf("err: %s", err)
	}

	fmt.Println(string(d))
	type config struct {
		ProxyEndpoints  []*endpoints.ProxyEndpoint  `hcl:"proxy_endpoint"`
		TargetEndpoints []*endpoints.TargetEndpoint `hcl:"target_endpoint"`
	}
	var pe config
	if err := xml.Unmarshal(d, &pe.ProxyEndpoints); err != nil {
		log.Fatalf("unmarshal err: %s", err)
	}

	hcl, err := hclencoder.Encode(pe)
	if err != nil {
		log.Fatal("unable to encode: ", err)
	}

	fmt.Fprintln(c.Stdout, string(hcl))

	return 0
}

func (c *cli) Usage() int {
	fmt.Fprintln(c.Stderr, "usage: hcl2json <file>")
	return 1
}
