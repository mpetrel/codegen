package tool

import "os"

func PathExist(pathName string) (bool, error) {
	_, err := os.Stat(pathName)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
