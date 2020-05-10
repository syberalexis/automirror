package engine

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/pkg/model"
	"os/exec"
	"strings"
)

type Wget struct {
	AbstractPuller
	Destination string `yaml:"destination"`
	Options     string `yaml:"options"`
}

func NewWget(config configs.EngineConfig, logger *log.Logger) (wget *Wget, err error) {
	err = configs.Parse(&wget, config.Config)
	wget.Puller = wget
	wget.Archives = config.Archives
	wget.Logger = logger
	return
}

func (wget *Wget) Pull(archive model.Archive) error {
	var args []string

	if len(wget.Options) > 0 {
		args = append(args, strings.Fields(wget.Options)...)
	}

	args = append(args, archive.Name)
	args = append(args, "-P")
	args = append(args, wget.Destination)

	cmd := exec.Command("wget", args...)
	cmd.Stdout = wget.Logger.Writer()
	cmd.Stderr = wget.Logger.Writer()
	return cmd.Run()
}
