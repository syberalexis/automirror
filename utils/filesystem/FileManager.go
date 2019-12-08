package filesystem

import (
	"io/ioutil"
	"os"
	"strings"
)

// Mkdir utils method to create a directory structure
func Mkdir(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0755)
	}
	return err
}

// Count utils method to count a directory structure
func Count(folder string) (int, error) {
	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}

// Combine directory and filename
func Combine(directory string, filename string) string {
	if directory != "" {
		if strings.HasSuffix(directory, "/") {
			return directory + filename
		}
		return directory + "/" + filename
	}
	return filename
}
