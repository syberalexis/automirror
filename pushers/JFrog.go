package pushers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type JFrog struct {
	Url    string
	ApiKey string
	Source string
}

func (j JFrog) Push() {
	files := j.getFiles()

	for _, file := range files {
		data, err := os.Open(j.Source + "/" + file)
		defer data.Close()
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("PUT", j.Url+"/"+file, data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("X-JFrog-Art-Api", j.ApiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			log.Print(file + " successfully pushed !")
		} else {
			log.Panic(resp.Status + " : Error when pushing " + file)
		}
	}
}

func (j JFrog) getFiles() []string {
	var files []string

	err := filepath.Walk(j.Source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, strings.TrimPrefix(path, j.Source+"/"))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}
