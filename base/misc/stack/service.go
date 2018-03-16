package stack

import (
	"os"
	"path/filepath"
	"strings"
)

func GetServiceNameBy(sp string) string {
	args := []string{filepath.Base(os.Args[0])}
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			break
		}
		args = append(args, arg)
	}
	return strings.Join(args, sp)
}

func GetRootService() string {
	return filepath.Base(os.Args[0])
}

func GetServiceName() string {
	return GetServiceNameBy(".")
}
