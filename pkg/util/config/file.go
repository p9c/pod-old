// Copyright (c) 2015 The btcsuite developers

package cfgutil

import "os"

// FileExists reports whether the named file or directory exists.
func FileExists(

	filePath string) (bool, error) {

	_, err := os.Stat(filePath)

	if err != nil {

		if os.IsNotExist(err) {

			return false, nil
		}
		return false, err
	}
	return true, nil
}
