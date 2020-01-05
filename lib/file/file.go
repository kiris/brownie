package file

import (
	"os"
)

//func IsExists(path string) bool {
//	_, err := os.Stat(path)
//	return err == nil
//}

func IsExistsDir(path string) bool {
	src, err := os.Stat(path)
	return err == nil && src.IsDir()
}

func IsExistsFile(path string) bool {
	src, err := os.Stat(path)
	return err == nil && !src.IsDir()
}