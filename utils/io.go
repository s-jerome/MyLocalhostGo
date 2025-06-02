package utils

import "os"

func FileExists(filePath string) (bool, error) {
	var stats, err = os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	_ = stats
	if os.IsNotExist(err) {
		return false, nil
	} else {
		return true, err
	}
}
