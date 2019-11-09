package pullers

import (
	"automirror/configs"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	Url     string
	Folder  string
	Options string
}

func BuildGit(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Git
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

func (g Git) Pull() (int, error) {
	err := g.mkdir()
	if err != nil {
		return -1, err
	}

	before, err := g.count()
	if err != nil {
		return before, err
	}

	err = g.clone()
	if err != nil {
		return -1, err
	}

	after, err := g.count()
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (g Git) Push() error {
	return nil
}

func (g Git) mkdir() error {
	_, err := os.Stat(g.Folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(g.Folder, 0755)
	}
	return err
}

func (g Git) count() (int, error) {
	files, err := ioutil.ReadDir(g.Folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}

// Private method to clone artifacts
func (g Git) clone() error {
	var args []string

	args = append(args, "clone")
	args = append(args, "--mirror")
	if len(g.Options) > 0 {
		args = append(args, strings.Fields(g.Options)...)
	}
	args = append(args, g.Url)
	args = append(args, g.Folder)

	cmd := exec.Command("git", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
