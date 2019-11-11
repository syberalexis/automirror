package both

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/utils"
	"os/exec"
	"strings"
)

type Rsync struct {
	Source      string
	Destination string
	Options     string
}

func NewRsync(config configs.EngineConfig) (interface{}, error) {
	var rsync Rsync
	err := configs.Parse(&rsync, config.Config)
	if err != nil {
		return nil, err
	}
	return rsync, nil
}

func (r Rsync) Pull() (int, error) {
	err := utils.Mkdir(r.Destination)
	if err != nil {
		return -1, err
	}

	before, err := utils.Count(r.Destination)
	if err != nil {
		return before, err
	}

	err = r.synchronize()
	if err != nil {
		return -1, err
	}

	after, err := utils.Count(r.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (r Rsync) Push() error {
	return r.synchronize()
}

// Private method to clone artifacts
func (r Rsync) synchronize() error {
	var args []string

	if len(r.Options) > 0 {
		args = append(args, strings.Fields(r.Options)...)
	}

	args = append(args, r.Source)
	args = append(args, r.Destination)

	cmd := exec.Command("rsync", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
