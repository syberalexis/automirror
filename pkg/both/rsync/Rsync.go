package rsync

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/utils/filesystem"
	"os/exec"
	"strings"
)

// Rsync object to pull and push with rsync unix command
type Rsync struct {
	Source      string
	Destination string
	Options     string
}

// NewRsync method to construct Rsync
func NewRsync(config configs.EngineConfig) (interface{}, error) {
	var rsync Rsync
	err := configs.Parse(&rsync, config.Config)
	if err != nil {
		return nil, err
	}
	return rsync, nil
}

// Pull a remote folder to a local
// Inherits public method to launch pulling process
// Return number of downloaded artifacts and error
func (r Rsync) Pull(log *log.Logger) (int, error) {
	err := filesystem.Mkdir(r.Destination)
	if err != nil {
		return -1, err
	}

	before, err := filesystem.Count(r.Destination)
	if err != nil {
		return before, err
	}

	err = r.synchronize(log)
	if err != nil {
		return -1, err
	}

	after, err := filesystem.Count(r.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Push a local folder to a remote
// Inherits public method to launch pushing process
// Return error
func (r Rsync) Push(log *log.Logger) error {
	return r.synchronize(log)
}

// Private method to clone artifacts
func (r Rsync) synchronize(log *log.Logger) error {
	var args []string

	if len(r.Options) > 0 {
		args = append(args, strings.Fields(r.Options)...)
	}

	args = append(args, r.Source)
	args = append(args, r.Destination)

	cmd := exec.Command("rsync", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
