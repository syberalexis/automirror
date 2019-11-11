package pullers

import (
	"encoding/json"
	"fmt"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/utils"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Docker struct {
	Source      string
	Destination string
	AuthUri     string  `toml:"auth_uri"`
	Images      []Image `toml:"images"`
}

type Image struct {
	Name string `toml:"name"`
}

func NewDocker(config configs.EngineConfig) (interface{}, error) {
	var docker Docker
	err := configs.Parse(&docker, config.Config)
	if err != nil {
		return nil, err
	}
	return docker, nil
}

func (d Docker) Pull() (int, error) {
	err := utils.Mkdir(d.Destination)
	if err != nil {
		return -1, err
	}

	before, err := utils.Count(d.Destination)
	if err != nil {
		return before, err
	}

	for _, image := range d.Images {
		i, err := d.getTags(image.Name)
		if err != nil {
			return 0, err
		}

		for _, tag := range i.Tags {
			err := utils.Mkdir(fmt.Sprintf("%s/%s", d.Destination, image.Name))
			if err != nil {
				return -1, err
			}
			err = d.pull(image.Name, tag)
			if err != nil {
				return 0, err
			}
			err = d.save(image.Name, tag)
			if err != nil {
				return 0, err
			}
		}
	}

	after, err := utils.Count(d.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}
func (d Docker) authenticate(image string) (authentication, error) {
	var authentication authentication

	resp, err := http.Get(strings.Join([]string{d.AuthUri, fmt.Sprintf("scope=repository:library/%s:pull", image)}, "&"))
	if err != nil {
		return authentication, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authentication, err
	}

	err = json.Unmarshal(body, &authentication)
	return authentication, err
}

func (d Docker) getTags(name string) (image, error) {
	var image image

	authentication, err := d.authenticate(name)
	if err != nil {
		return image, err
	}

	req, err := http.NewRequest("GET", strings.Join([]string{d.Source, "library", name, "tags", "list"}, "/"), nil)
	if err != nil {
		return image, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authentication.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return image, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &image)
	return image, err
}

func (d Docker) pull(name string, tag string) error {
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", name, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Docker) save(name string, tag string) error {
	cmd := exec.Command("docker", "save", name, "--output", fmt.Sprintf("%s/%s/%s", d.Destination, name, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type authentication struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IssuedAt    string `json:"issued_at"`
}

type image struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
