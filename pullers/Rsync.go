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

type Rsync struct {
	Url     string
	Folder  string
	Options string
}

func BuildRsync(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Rsync
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

func (r Rsync) Pull() (int, error) {
	err := r.mkdir()
	if err != nil {
		return -1, err
	}

	before, err := r.count()
	if err != nil {
		return before, err
	}

	err = r.download()
	if err != nil {
		return -1, err
	}

	after, err := r.count()
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (r Rsync) mkdir() error {
	_, err := os.Stat(r.Folder)
	if os.IsNotExist(err) {
		return os.Mkdir(r.Folder, 0755)
	}
	return err
}

func (r Rsync) count() (int, error) {
	files, err := ioutil.ReadDir(r.Folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}

// Private method to download artifacts
func (r Rsync) download() error {
	var args []string

	if len(r.Options) > 0 {
		args = append(args, strings.Fields(r.Options)...)
	}

	args = append(args, r.Url)
	args = append(args, r.Folder)

	cmd := exec.Command("rsync", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
