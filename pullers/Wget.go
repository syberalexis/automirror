package pullers

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Wget struct {
	Url     string
	Folder  string
	Options string
}

func BuildWget(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Wget
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

func (w Wget) Pull() (int, error) {
	err := w.mkdir()
	if err != nil {
		return -1, err
	}

	before, err := w.count()
	if err != nil {
		return before, err
	}

	err = w.download()
	if err != nil {
		return -1, err
	}

	after, err := w.count()
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (w Wget) mkdir() error {
	_, err := os.Stat(w.Folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(w.Folder, 0755)
	}
	return err
}

func (w Wget) count() (int, error) {
	files, err := ioutil.ReadDir(w.Folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}

// Private method to clone artifacts
func (w Wget) download() error {
	var args []string

	if len(w.Options) > 0 {
		args = append(args, strings.Fields(w.Options)...)
	}

	args = append(args, w.Url)
	args = append(args, "-P")
	args = append(args, w.Folder)

	cmd := exec.Command("wget", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
