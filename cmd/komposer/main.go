package main

import (
	"fmt"
	"github.com/Doublemine/komposer/core"
	"github.com/akamensky/argparse"
	"os"
)

func main() {

	parser := argparse.NewParser("komposer", "a terminal tools that let you multiple k8s config file merge as single one")
	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "show composer version"})
	force := parser.Flag("f", "force", &argparse.Options{Required: false, Help: "force overwrite kubeconfig file if exists"})
	special := parser.Flag("s", "special", &argparse.Options{Required: false, Help: "save on user dir after compose"})
	configList := parser.List("c", "config", &argparse.Options{Required: false, Help: "the kubeconfig file path, can be repeat. such as: -c /path/to/config -c /path/to/others "})
	duplicateSuffix := parser.String("r", "suffix", &argparse.Options{Required: false, Help: "the duplicates config name automation append suffix, default: -k6r"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Println(parser.Usage(err))
	}

	if *version {
		core.ShowVersion()
		return
	}

	if len(*configList) > 1 {
		err = core.KubeConfigVerifier(*configList)
		if err != nil {
			fmt.Println(parser.Usage(err))
			return
		}
		core.Compose(*configList, *force, *special, *duplicateSuffix)
		return
	}

	fmt.Println(parser.Usage(nil))
}
