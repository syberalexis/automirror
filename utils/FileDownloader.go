package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileDownloader(url string, file string) error {
	// Create folders
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(file)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
