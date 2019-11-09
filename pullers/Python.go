package pullers

import (
	"automirror/configs"
	"automirror/utils"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Python struct {
	Url            string
	Folder         string
	DatabaseFile   string `toml:"database_file"`
	FileExtensions string `toml:"file_extensions"`
}

func BuildPython(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Python
	tomlFile, err := ioutil.ReadFile(pullerConfig.Config)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		return nil, err
	}

	config.Url = pullerConfig.Source
	config.Folder = pullerConfig.Destination
	return config, nil
}

func (p Python) Pull() (int, error) {
	err := utils.InitializeDatabase(p.DatabaseFile)
	if err != nil {
		return 0, err
	}
	return p.readRepository("/simple/")
}

// Private method to get archive list of artifact to clone
func (p Python) readRepository(subpath string) (int, error) {
	counter := 0
	resp, err := http.Get(p.Url + subpath)
	if err != nil {
		return counter, err
	}

	z := html.NewTokenizer(resp.Body)
	finished := false
	for !finished {
		next := z.Next()
		switch {
		case next == html.ErrorToken:
			// End of the document, we're done
			finished = true
			break
		case next == html.StartTagToken:
			token := z.Token()
			if token.Data == "a" && len(token.Attr) > 0 && token.Attr[0].Val != "../" {
				if strings.HasSuffix(token.Attr[0].Val, "/") {
					count, err := p.readRepository(token.Attr[0].Val)
					counter += count
					if err != nil {
						return counter, err
					}
				} else {
					match := p.match(token.Attr[0].Val)
					isExist, err := utils.ExistsInDatabase(p.DatabaseFile, match)
					if err != nil {
						return counter, err
					}
					if match != "" && !isExist {
						err := p.download(subpath, token.Attr[0].Val)
						if err != nil {
							return counter, err
						}
						counter++
						time.Sleep(time.Millisecond * 500)
					}
				}
			}
		}
	}
	return counter, nil
}

func (p Python) match(url string) string {
	re := regexp.MustCompile(fmt.Sprintf("^.*/(.+\\.(%s))#?.*$", p.FileExtensions))
	match := re.FindStringSubmatch(url)
	if match != nil {
		return match[1]
	}
	return ""
}

// Private method to clone artifacts
func (p Python) download(subpath string, url string) error {
	match := p.match(url)
	if match != "" {
		file := strings.Join([]string{p.Folder, strings.Replace(subpath, "/simple", "", 1), match}, "")

		if err := utils.FileDownloader(url, file); err != nil {
			return err
		}
		log.Infof("%s successfully pulled !\n", file)

		err := utils.InsertIntoDatabase(p.DatabaseFile, match)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("%s not matched", url)
	}
	return nil
}
