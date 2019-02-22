package core

import (
	"errors"
	"os"
)

func KubeConfigVerifier(args []string) error {

	if len(args) <= 1 {
		return errors.New("the config path must be more than single one")
	}

	for _, path := range args {
		if !FileExist(path) {
			return errors.New("the config file path: " + path + " not exist or it's a directory, please check it.")
		}
	}

	return nil
}

func FileExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if fi.Mode().IsRegular() {
		return true
	}
	return false
}
