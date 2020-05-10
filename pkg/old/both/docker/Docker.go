package docker

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/old"
	"github.com/syberalexis/automirror/utils/filesystem"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Docker object to pull and push Docker images with docker unix command
type Docker struct {
	Source      string
	Destination string
	AuthURI     string  `toml:"auth_uri"`
	Images      []Image `toml:"images"`
}

// Image structure from specific configuration file
type Image struct {
	Name string `toml:"name"`
}

// NewDocker method to construct Docker
func NewDocker(config old.EngineConfig) (interface{}, error) {
	var docker Docker
	err := old.Parse(&docker, config.Config)
	if err != nil {
		return nil, err
	}
	return docker, nil
}

// Pull Docker images
// Inherits public method to launch pulling process
// Return number of downloaded artifacts and error
func (d Docker) Pull(log *log.Logger) (int, error) {
	err := filesystem.Mkdir(d.Destination)
	if err != nil {
		return -1, err
	}

	before, err := filesystem.Count(d.Destination)
	if err != nil {
		return before, err
	}

	for _, image := range d.Images {
		i, err := d.getTags(image.Name)
		if err != nil {
			return 0, err
		}

		for _, tag := range i.Tags {
			err := filesystem.Mkdir(fmt.Sprintf("%s/%s", d.Destination, image.Name))
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

	after, err := filesystem.Count(d.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Push Docker images
// Inherits public method to launch pushing process
// Return error
func (d Docker) Push() error {
	return nil
}

// private method to authenticate with anonymous on Docker Registry
func (d Docker) authenticate(image string) (authentication, error) {
	var authentication authentication

	resp, err := http.Get(strings.Join([]string{d.AuthURI, fmt.Sprintf("scope=repository:library/%s:pull", image)}, "&"))
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

// private method to get Tags from a Docker image
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
	if err != nil {
		return image, err
	}

	err = json.Unmarshal(body, &image)
	return image, err
}

// private method to run docker pull unix command
func (d Docker) pull(name string, tag string) error {
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", name, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// private method to run docker save unix command
func (d Docker) save(name string, tag string) error {
	cmd := exec.Command("docker", "save", name, "--output", fmt.Sprintf("%s/%s/%s", d.Destination, name, tag))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Structure returned by authentication request
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
