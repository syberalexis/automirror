package git

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/utils/filesystem"
	"os/exec"
	"strings"
)

// Git object to pull and push with git unix command
type Git struct {
	Source      string
	Destination string
	Options     string
}

// NewGit method to construct Git
func NewGit(config configs.EngineConfig) (interface{}, error) {
	var git Git
	err := configs.Parse(&git, config.Config)
	if err != nil {
		return nil, err
	}
	return git, nil
}

// Pull pull a git repo
// Inherits public method to launch pulling process
// Return number of downloaded artifacts and error
func (g Git) Pull(log *log.Logger) (int, error) {
	err := filesystem.Mkdir(g.Destination)
	if err != nil {
		return -1, err
	}

	before, err := filesystem.Count(g.Destination)
	if err != nil {
		return before, err
	}

	err = g.clone(log)
	if err != nil {
		return -1, err
	}

	after, err := filesystem.Count(g.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Push a git repo
// Inherits public method to launch pushing process
// Return error
func (g Git) Push() error {
	return nil
}

// Private method to clone artifacts
func (g Git) clone(log *log.Logger) error {
	var args []string

	args = append(args, "clone")
	args = append(args, "--mirror")
	if len(g.Options) > 0 {
		args = append(args, strings.Fields(g.Options)...)
	}
	args = append(args, g.Source)
	args = append(args, g.Destination)

	cmd := exec.Command("git", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}

// Private method to push artifacts TODO
func (g Git) push(log *log.Logger) error {
	var args []string

	args = append(args, "push")
	args = append(args, "--mirror")
	if len(g.Options) > 0 {
		args = append(args, strings.Fields(g.Options)...)
	}
	args = append(args, g.Source)
	args = append(args, g.Destination)

	cmd := exec.Command("git", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
