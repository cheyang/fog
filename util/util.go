package util

import (
	"os"
	"path/filepath"
)

func GetStorePath(name string) (storePath string, err error) {
	var pwd string
	if pwd, err = os.Getwd(); err != nil {
		return
	} else {
		storePath = filepath.Join(pwd, ".fog", name)
	}

	return
}
