package main

import (
	"flag"
	"github.com/kevinswiber/apigee-hcl/cli"
	"path"
)

func main() {
	var options cli.CLIOptions

	flag.StringVar(&options.InputHCL, "i", "", "Required. An HCL file to translate")
	flag.StringVar(&options.BuildPath, "o", path.Join(".", "build"), "Optional. A build path")
	flag.StringVar(&options.ResourcesPath, "r", path.Join(".", "resources"), "Optional. A path to resources")
	flag.Parse()

	if options.InputHCL == "" {
		flag.Usage()
		return
	}

	cli.StartCLI(&options)
}
