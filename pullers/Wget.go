package pullers

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/utils"
	"os/exec"
	"strings"
)

type Wget struct {
	Source      string
	Destination string
	Options     string
}

func NewWget(config configs.EngineConfig) (interface{}, error) {
	var wget Wget
	err := configs.Parse(&wget, config.Config)
	if err != nil {
		return nil, err
	}
	return wget, nil
}

func (w Wget) Pull() (int, error) {
	err := utils.Mkdir(w.Destination)
	if err != nil {
		return -1, err
	}

	before, err := utils.Count(w.Destination)
	if err != nil {
		return before, err
	}

	err = w.download()
	if err != nil {
		return -1, err
	}

	after, err := utils.Count(w.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Private method to clone artifacts
func (w Wget) download() error {
	var args []string

	if len(w.Options) > 0 {
		args = append(args, strings.Fields(w.Options)...)
	}

	args = append(args, w.Source)
	args = append(args, "-P")
	args = append(args, w.Destination)

	cmd := exec.Command("wget", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
