package pushers

import (
	"automirror/configs"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type JFrog struct {
	Url           string `toml:"url"`
	ApiKey        string `toml:"api_key"`
	Source        string `toml:"source"`
	ExcludeRegexp string `toml:"exclude_regexp"`
}

func BuildJFrog(pusherConfig configs.PusherConfig) Pusher {
	var config JFrog
	tomlFile, err := ioutil.ReadFile(pusherConfig.Config)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func (j JFrog) Push() {
	files := j.getFiles()

	for _, file := range files {
		if !j.fileExists(file) {
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
				fmt.Printf("%s successfully pushed !\n", file)
			} else {
				log.Panic(resp.Status + " : Error when pushing " + file)
			}
		}
	}
}

func (j JFrog) fileExists(file string) bool {
	req, err := http.NewRequest("GET", j.Url+"/"+file, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-JFrog-Art-Api", j.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode == 200
}

func (j JFrog) getFiles() []string {
	var files []string

	err := filepath.Walk(j.Source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matched, _ := regexp.MatchString(j.ExcludeRegexp, info.Name())
			if !matched {
				files = append(files, strings.TrimPrefix(path, j.Source+"/"))
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}
