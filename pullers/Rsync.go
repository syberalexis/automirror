package pullers

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/pushers"
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

func BuildRsyncPusher(pusherConfig configs.PusherConfig) (pushers.Pusher, error) {
	var config Rsync
	tomlFile, err := ioutil.ReadFile(pusherConfig.Config)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		return nil, err
	}

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

	err = r.synchronize()
	if err != nil {
		return -1, err
	}

	after, err := r.count()
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (r Rsync) Push() error {
	return r.synchronize()
}

func (r Rsync) mkdir() error {
	_, err := os.Stat(r.Folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(r.Folder, 0755)
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

// Private method to clone artifacts
func (r Rsync) synchronize() error {
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
