package cli

import (
	"encoding/xml"
	"fmt"
	"github.com/crufter/copyrecur"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hclParser "github.com/hashicorp/hcl/hcl/parser"
	"github.com/kevinswiber/apigee-hcl/dsl"
	hclEncoding "github.com/kevinswiber/apigee-hcl/dsl/encoding/hcl"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// InputValues is an array of input files
type InputValues []string

// String is part of an implementation of the flag.Value interface
func (v *InputValues) String() string {
	return ""
}

// Set is part of an implementation of the flag.Value interface
func (v *InputValues) Set(value string) error {
	*v = append(*v, value)
	return nil
}

// Options is an arguments container for running the CLI.
type Options struct {
	InputHCL      InputValues
	BuildPath     string
	ResourcesPath string
}

// Start runs the command line utility logic.
func Start(opts *Options) {
	var errors error
	l := log.New(os.Stderr, "", 0)

	buildPath := opts.BuildPath
	bundlePath := path.Join(buildPath, "apiproxy")
	proxyEndpointsPath := path.Join(bundlePath, "proxies")
	targetEndpointsPath := path.Join(bundlePath, "targets")
	policiesPath := path.Join(bundlePath, "policies")
	resourcesPath := path.Join(bundlePath, "resources")

	if err := os.RemoveAll(bundlePath); err != nil {
		errors = multierror.Append(errors, err)
		l.Fatal(errors)
	}

	var c dsl.Config

	for _, file := range opts.InputHCL {
		d, err := ioutil.ReadFile(file)
		if err != nil {
			errors = multierror.Append(err)
			l.Fatal(errors)
		}

		hclRoot, err := hcl.Parse(string(d))
		if err != nil {
			switch err.(type) {
			case *hclParser.PosError:
				e := err.(*hclParser.PosError)
				e2 := &hclEncoding.PosError{
					Pos: e.Pos,
					Err: e.Err,
				}
				e2.Pos.Filename = file
				errors = multierror.Append(errors, e2)
			default:
				errors = multierror.Append(errors, err)
			}
			l.Fatal(errors)
		}

		list, ok := hclRoot.Node.(*ast.ObjectList)
		if !ok {
			errors = multierror.Append(errors,
				fmt.Errorf("file doesn't contain root object"))
			l.Fatal(errors)
		}

		cfg, err := dsl.DecodeConfigHCL(list)
		if err != nil {
			if merr, ok := err.(*multierror.Error); ok {
				attachFilenameToPosErrors(file, merr)
			}
			errors = multierror.Append(errors, err)
			l.Fatal(errors)
		}

		if cfg.Proxy != nil && cfg.Proxy.Name != "" {
			c.Proxy = cfg.Proxy
		}

		c.ProxyEndpoints = append(c.ProxyEndpoints, cfg.ProxyEndpoints...)
		c.TargetEndpoints = append(c.TargetEndpoints, cfg.TargetEndpoints...)
		c.Policies = append(c.Policies, cfg.Policies...)

		if cfg.Resources != nil {
			if c.Resources == nil {
				c.Resources = make(map[string]string)
			}
			for k, v := range cfg.Resources {
				c.Resources[k] = v
			}
		}
	}

	// Validate

	if c.Proxy == nil {
		errors = multierror.Append(errors,
			fmt.Errorf("no proxy definition found"))
	}

	if len(c.ProxyEndpoints) == 0 {
		errors = multierror.Append(errors,
			fmt.Errorf("no proxy endpoint definition found"))
	}

	if errors != nil {
		l.Fatal(errors)
	}

	// Output

	output, err := xml.MarshalIndent(c.Proxy, "", "    ")
	if err != nil {
		errors = multierror.Append(errors, err)
		l.Fatal(errors)
	}

	if err := ensureDirectory(bundlePath); err != nil {
		errors = multierror.Append(errors, err)
		l.Fatal(errors)
	}

	apiProxyPath := path.Join(bundlePath, c.Proxy.Name+".xml")

	output = []byte(xml.Header + string(output))
	if err := ioutil.WriteFile(apiProxyPath, output, 0666); err != nil {
		errors = multierror.Append(errors, err)
		l.Fatal(errors)
	}

	if len(c.ProxyEndpoints) > 0 {
		if err := ensureDirectory(proxyEndpointsPath); err != nil {
			errors = multierror.Append(errors, err)
			l.Fatal(errors)
		}

		for _, proxyEndpoint := range c.ProxyEndpoints {
			var output []byte
			output, err := xml.MarshalIndent(proxyEndpoint, "", "    ")
			if err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}

			output = []byte(xml.Header + string(output))
			ePath := path.Join(proxyEndpointsPath, proxyEndpoint.Name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}
		}
	}

	if len(c.TargetEndpoints) > 0 {
		if err := ensureDirectory(targetEndpointsPath); err != nil {
			errors = multierror.Append(errors, err)
			l.Fatal(errors)
		}

		for _, targetEndpoint := range c.TargetEndpoints {
			var output []byte
			output, err := xml.MarshalIndent(targetEndpoint, "", "    ")
			if err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}

			output = []byte(xml.Header + string(output))
			ePath := path.Join(targetEndpointsPath, targetEndpoint.Name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}
		}
	}

	if len(c.Policies) > 0 {
		if err := ensureDirectory(policiesPath); err != nil {
			errors = multierror.Append(errors, err)
			l.Fatal(errors)
		}

		for _, policy := range c.Policies {
			var output []byte
			output, err := xml.MarshalIndent(policy, "", "    ")
			if err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}

			name := policy.Name()

			output = []byte(xml.Header + string(output))
			ePath := path.Join(policiesPath, name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}
		}
	}

	if stat, err := os.Stat(opts.ResourcesPath); err == nil {
		if stat.IsDir() {
			if err := copyrecur.CopyDir(opts.ResourcesPath, resourcesPath); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}
		}
	}

	if c.Resources != nil {
		if err := ensureDirectory(resourcesPath); err != nil {
			errors = multierror.Append(errors, err)
			l.Fatal(errors)
		}

		for fileName, content := range c.Resources {
			parts := strings.Split(fileName, "://")
			lang := parts[0]
			dir := path.Join(resourcesPath, lang)
			if err := ensureDirectory(dir); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}

			fileContents := []byte(content)
			if err := ioutil.WriteFile(path.Join(dir, parts[1]), fileContents, 0666); err != nil {
				errors = multierror.Append(errors, err)
				l.Fatal(errors)
			}
		}
	}
}

func ensureDirectory(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func attachFilenameToPosErrors(file string, errors *multierror.Error) {
	for _, e := range errors.Errors {
		switch e.(type) {
		case *hclEncoding.PosError:
		case *hclParser.PosError:
			e2 := e.(*hclEncoding.PosError)
			e2.Pos.Filename = file
		case *multierror.Error:
			attachFilenameToPosErrors(file, e.(*multierror.Error))
		}
	}
}
