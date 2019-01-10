package core

import (
	"errors"
	"fmt"
	"os"
)

func KubeConfigVerifier(args []string) error {
	fmt.Println(args)

	if len(args) <= 1 {
		return errors.New("the config path must be more than single one")
	}

	for _, path := range args {
		if !PathExist(path) {
			return errors.New("the config file path: " + path + " not exist, please check it.")
		}
	}

	return nil
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
