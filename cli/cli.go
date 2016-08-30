package cli

import (
	"encoding/xml"
	"github.com/crufter/copyrecur"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

type CLIOptions struct {
	InputHCL      string
	BuildPath     string
	ResourcesPath string
}

func StartCLI(opts *CLIOptions) {
	buildPath := opts.BuildPath
	bundlePath := path.Join(buildPath, "apiproxy")
	proxyEndpointsPath := path.Join(bundlePath, "proxies")
	targetEndpointsPath := path.Join(bundlePath, "targets")
	policiesPath := path.Join(bundlePath, "policies")
	resourcesPath := path.Join(bundlePath, "resources")

	file := opts.InputHCL

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

	if err := ensureDirectory(bundlePath); err != nil {
		log.Fatalf("err: %s", err)
	}

	apiProxyPath := path.Join(bundlePath, c.Proxy.Name+".xml")

	output = []byte(xml.Header + string(output))
	if err := ioutil.WriteFile(apiProxyPath, output, 0666); err != nil {
		log.Fatalf("err: %s", err)
	}

	if len(c.ProxyEndpoints) > 0 {
		if err := ensureDirectory(proxyEndpointsPath); err != nil {
			log.Fatalf("err: %s", err)
		}

		for _, proxyEndpoint := range c.ProxyEndpoints {
			var output []byte
			output, err := xml.MarshalIndent(proxyEndpoint, "", "    ")
			if err != nil {
				log.Fatalf("err: %s", err)
			}

			output = []byte(xml.Header + string(output))
			ePath := path.Join(proxyEndpointsPath, proxyEndpoint.Name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				log.Fatalf("err: %s", err)
			}
		}
	}

	if len(c.TargetEndpoints) > 0 {
		if err := ensureDirectory(targetEndpointsPath); err != nil {
			log.Fatalf("err: %s", err)
		}

		for _, targetEndpoint := range c.TargetEndpoints {
			var output []byte
			output, err := xml.MarshalIndent(targetEndpoint, "", "    ")
			if err != nil {
				log.Fatalf("err: %s", err)
			}

			output = []byte(xml.Header + string(output))
			ePath := path.Join(targetEndpointsPath, targetEndpoint.Name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				log.Fatalf("err: %s", err)
			}
		}
	}

	if len(c.Policies) > 0 {
		if err := ensureDirectory(policiesPath); err != nil {
			log.Fatalf("err: %s", err)
		}

		for _, policy := range c.Policies {
			var output []byte
			output, err := xml.MarshalIndent(policy, "", "    ")
			if err != nil {
				log.Fatalf("err: %s", err)
			}

			val := reflect.ValueOf(policy)
			name := val.FieldByName("Name").String()

			output = []byte(xml.Header + string(output))
			ePath := path.Join(policiesPath, name+".xml")
			if err := ioutil.WriteFile(ePath, output, 0666); err != nil {
				log.Fatalf("err: %s", err)
			}
		}
	}

	if stat, err := os.Stat(opts.ResourcesPath); err == nil {
		if stat.IsDir() {
			if err := copyrecur.CopyDir(opts.ResourcesPath, resourcesPath); err != nil {
				log.Fatalf("err :%s", err)
			}
		}
	}

	if c.Resources != nil {
		if err := ensureDirectory(resourcesPath); err != nil {
			log.Fatalf("err: %s", err)
		}

		for fileName, content := range c.Resources {
			parts := strings.Split(fileName, "://")
			lang := parts[0]
			dir := path.Join(resourcesPath, lang)
			if err := ensureDirectory(dir); err != nil {
				log.Fatalf("err: %s", err)
			}

			fileContents := []byte(content)
			if err := ioutil.WriteFile(path.Join(dir, parts[1]), fileContents, 0666); err != nil {
				log.Fatalf("err: %s", err)
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
