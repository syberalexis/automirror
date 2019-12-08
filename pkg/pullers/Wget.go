package pullers

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/utils/filesystem"
	"os/exec"
	"strings"
)

// Wget object to pull web resource with wget unix command
type Wget struct {
	Source      string
	Destination string
	Options     string
}

// NewWget method to construct Wget
func NewWget(config configs.EngineConfig) (interface{}, error) {
	var wget Wget
	err := configs.Parse(&wget, config.Config)
	if err != nil {
		return nil, err
	}
	return wget, nil
}

// Pull resource from url
// Inherits public method to launch pulling process
// Return number of downloaded artifacts and error
func (w Wget) Pull(log *log.Logger) (int, error) {
	err := filesystem.Mkdir(w.Destination)
	if err != nil {
		return -1, err
	}

	before, err := filesystem.Count(w.Destination)
	if err != nil {
		return before, err
	}

	err = w.download(log)
	if err != nil {
		return -1, err
	}

	after, err := filesystem.Count(w.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Private method to clone artifacts
func (w Wget) download(log *log.Logger) error {
	var args []string

	if len(w.Options) > 0 {
		args = append(args, strings.Fields(w.Options)...)
	}

	args = append(args, w.Source)
	args = append(args, "-P")
	args = append(args, w.Destination)

	cmd := exec.Command("wget", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
