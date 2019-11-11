package utils

import (
	"io/ioutil"
	"os"
)

func Mkdir(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0755)
	}
	return err
}

func Count(folder string) (int, error) {
	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}
