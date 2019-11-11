package pullers

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/syberalexis/automirror/configs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Docker struct {
	Url     string
	Folder  string
	AuthUri string  `toml:"auth_uri"`
	Images  []Image `toml:"images"`
}

type Image struct {
	Name string `toml:"name"`
}

func BuildDocker(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Docker
	tomlFile, err := ioutil.ReadFile(pullerConfig.Config)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}

	config.Url = pullerConfig.Source
	config.Folder = pullerConfig.Destination
	return config, nil
}

func (d Docker) Pull() (int, error) {
	err := d.mkdir(d.Folder)
	if err != nil {
		return -1, err
	}

	before, err := d.count()
	if err != nil {
		return before, err
	}

	for _, image := range d.Images {
		i, err := d.getTags(image.Name)
		if err != nil {
			return 0, err
		}

		for _, tag := range i.Tags {
			err := d.mkdir(fmt.Sprintf("%s/%s", d.Folder, image.Name))
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

	after, err := d.count()
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

	req, err := http.NewRequest("GET", strings.Join([]string{d.Url, "library", name, "tags", "list"}, "/"), nil)
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
	cmd := exec.Command("docker", "save", name, "--output", fmt.Sprintf("%s/%s/%s", d.Folder, name, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Docker) mkdir(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0755)
	}
	return err
}

func (d Docker) count() (int, error) {
	files, err := ioutil.ReadDir(d.Folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
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
