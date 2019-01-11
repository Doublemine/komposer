package core

import (
	"fmt"
	"runtime"
)

const (
	version = "1.0"
)

func ShowVersion() {
	fmt.Printf("komposer version: %s, go version: %s %s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
