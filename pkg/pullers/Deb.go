package pullers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/utils/filesystem"
	"os/exec"
	"strings"
)

// Deb object to pull a Debian based mirror with debmirror unix command
type Deb struct {
	Host        string
	Destination string
	Dist        string
	Arch        string
	Section     string
	Root        string
	Method      string
	Keyring     string
	Cleanup     bool
	Source      bool
	I18N        bool
	Options     string
}

// NewDeb method to construct Deb
func NewDeb(config configs.EngineConfig) (interface{}, error) {
	var deb Deb
	err := configs.Parse(&deb, config.Config)
	if err != nil {
		return nil, err
	}
	return deb, nil
}

// Pull a Debian based repo
// Inherits public method to launch pulling process
// Return number of downloaded artifacts and error
func (d Deb) Pull(log *log.Logger) (int, error) {
	err := filesystem.Mkdir(d.Destination)
	if err != nil {
		return -1, err
	}

	before, err := filesystem.Count(d.Destination)
	if err != nil {
		return before, err
	}

	err = d.download(log)
	if err != nil {
		return -1, err
	}

	after, err := filesystem.Count(d.Destination)
	if err != nil {
		return after, err
	}

	return after - before, nil
}

// Private method to clone artifacts
func (d Deb) download(log *log.Logger) error {
	var args []string

	if len(d.Host) > 0 {
		args = append(args, fmt.Sprintf("--host=%s", d.Host))
	}

	if len(d.Arch) > 0 {
		args = append(args, fmt.Sprintf("--arch=%s", d.Arch))
	}

	if len(d.Dist) > 0 {
		args = append(args, fmt.Sprintf("--dist=%s", d.Dist))
	}

	if len(d.Section) > 0 {
		args = append(args, fmt.Sprintf("--section=%s", d.Section))
	}

	if len(d.Root) > 0 {
		args = append(args, fmt.Sprintf("--root=%s", d.Root))
	}

	if len(d.Method) > 0 {
		args = append(args, fmt.Sprintf("--method=%s", d.Method))
	}

	if len(d.Keyring) > 0 {
		args = append(args, fmt.Sprintf("--keyring=%s", d.Keyring))
	}

	if d.Source {
		args = append(args, "--source")
	} else {
		args = append(args, "--nosource")
	}

	if d.I18N {
		args = append(args, "--i18n")
	}

	if !d.Cleanup {
		args = append(args, "--nocleanup")
	}

	if len(d.Options) > 0 {
		args = append(args, strings.Fields(d.Options)...)
	}

	args = append(args, d.Destination)

	cmd := exec.Command("debmirror", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
