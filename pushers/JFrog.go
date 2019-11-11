package pushers

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type JFrog struct {
	Source        string `toml:"source"`
	Destination   string `toml:"destination"`
	ApiKey        string `toml:"api_key"`
	ExcludeRegexp string `toml:"exclude_regexp"`
}

func NewJFrog(config configs.EngineConfig) (interface{}, error) {
	var jFrog JFrog
	err := configs.Parse(&jFrog, config.Config)
	if err != nil {
		return nil, err
	}
	return jFrog, nil
}

func (j JFrog) Push() error {
	files, err := j.getFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		isExist, err := j.fileExists(file)
		if err != nil {
			return err
		}
		if !isExist {
			data, err := os.Open(j.Source + "/" + file)
			defer data.Close()
			if err != nil {
				return err
			}

			req, err := http.NewRequest("PUT", j.Destination+"/"+file, data)
			if err != nil {
				return err
			}
			req.Header.Set("X-JFrog-Art-Api", j.ApiKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				log.Infof("%s successfully pushed !\n", file)
			} else {
				log.Errorf("%s : Error when pushing %s", resp.Status, file)
			}
		}
	}
	return nil
}

func (j JFrog) fileExists(file string) (bool, error) {
	req, err := http.NewRequest("GET", j.Destination+"/"+file, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("X-JFrog-Art-Api", j.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200, nil
}

func (j JFrog) getFiles() ([]string, error) {
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
		return nil, err
	}

	return files, nil
}
