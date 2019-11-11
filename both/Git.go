package both

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/utils"
	"os/exec"
	"strings"
)

type Git struct {
	Source      string
	Destination string
	Options     string
}

func NewGit(config configs.EngineConfig) (interface{}, error) {
	var git Git
	err := configs.Parse(&git, config.Config)
	if err != nil {
		return nil, err
	}
	return git, nil
}

func (g Git) Pull() (int, error) {
	err := utils.Mkdir(g.Destination)
	if err != nil {
		return -1, err
	}

	before, err := utils.Count(g.Destination)
	if err != nil {
		return before, err
	}

	err = g.clone()
	if err != nil {
		return -1, err
	}

	after, err := utils.Count(g.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (g Git) Push() error {
	return nil
}

// Private method to clone artifacts
func (g Git) clone() error {
	var args []string

	args = append(args, "clone")
	args = append(args, "--mirror")
	if len(g.Options) > 0 {
		args = append(args, strings.Fields(g.Options)...)
	}
	args = append(args, g.Source)
	args = append(args, g.Destination)

	cmd := exec.Command("git", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}

// Private method to push artifacts TODO
func (g Git) push() error {
	var args []string

	args = append(args, "push")
	args = append(args, "--mirror")
	if len(g.Options) > 0 {
		args = append(args, strings.Fields(g.Options)...)
	}
	args = append(args, g.Source)
	args = append(args, g.Destination)

	cmd := exec.Command("git", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}