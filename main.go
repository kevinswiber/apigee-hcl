package main

import (
	"flag"
	"github.com/kevinswiber/apg-hcl/cli"
	"path"
)

func main() {
	var options cli.CLIOptions

	flag.StringVar(&options.InputHCL, "i", "", "An HCL file to translate")
	flag.StringVar(&options.BuildPath, "o", path.Join(".", "build"), "A build path, default: ./build")
	flag.StringVar(&options.ResourcesPath, "r", path.Join(".", "resources"), "A path to resources, default: ./resources")
	flag.Parse()

	cli.StartCLI(&options)
}
