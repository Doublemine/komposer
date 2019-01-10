package main

import (
	"fmt"
	"github.com/Doublemine/kubeconfig-mixer/core"
	"github.com/akamensky/argparse"
	"os"
)

func main() {

	parser := argparse.NewParser("composer", "a terminal tools that let you multiple k8s config file merge as single one")
	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "show composer version"})
	force := parser.Flag("f", "force", &argparse.Options{Required: false, Help: "force overwrite kubeconfig file if exists"})
	special := parser.Flag("s", "special", &argparse.Options{Required: false, Help: "save on user dir after compose"})
	configList := parser.List("c", "config", &argparse.Options{Required: true, Help: "the kubeconfig file path, can be repeat. such as: -c /path/to/config -c /path/to/others "})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Println(parser.Usage(err))
	}

	err = core.KubeConfigVerifier(*configList)
	if err != nil {
		fmt.Println(parser.Usage(err))
	}
	// Finally print the collected string
	fmt.Println(*version, *force, *special, *configList)
}
