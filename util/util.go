package util

import "path/filepath"

func GetStorePath(name string) (storePath string, err error) {
	var etc string = "/etc"
	// if etc, err = os.Getwd(); err != nil {
	// 	return storePath, err
	// } else {
	storePath = filepath.Join(etc, ".fog", name)
	// }

	return storePath, err
}
