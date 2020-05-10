package filesystem

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileDownloader method to download a file from an URI
func FileDownloader(url string, file string) error {
	// Create folders
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(file)
	defer out.Close()
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
