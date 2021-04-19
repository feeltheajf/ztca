package fs

import (
	"io/ioutil"
	"os"
)

const (
	defaultPermissionsFile      = 0600
	defaultPermissionsDirectory = 0700
)

func UserConfigDir() string {
	path, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return path
}

func WriteFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, defaultPermissionsFile)
}

func Mkdir(path string) error {
	err := os.MkdirAll(path, defaultPermissionsDirectory)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}
