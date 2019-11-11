package pullers

import (
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Deb struct {
	Url     string
	Folder  string
	Dist    string
	Arch    string
	Section string
	Root    string
	Method  string
	Keyring string
	Cleanup bool
	Source  bool
	I18N    bool
	Options string
}

func BuildDeb(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Deb
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

func (d Deb) Pull() (int, error) {
	err := d.mkdir()
	if err != nil {
		return -1, err
	}

	before, err := d.count()
	if err != nil {
		return before, err
	}

	err = d.download()
	if err != nil {
		return -1, err
	}

	after, err := d.count()
	if err != nil {
		return after, err
	}

	return after - before, nil
}

func (d Deb) mkdir() error {
	_, err := os.Stat(d.Folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(d.Folder, 0755)
	}
	return err
}

func (d Deb) count() (int, error) {
	files, err := ioutil.ReadDir(d.Folder)

	if err != nil {
		return -1, err
	}

	return len(files), nil
}

// Private method to clone artifacts
func (d Deb) download() error {
	var args []string

	if len(d.Url) > 0 {
		args = append(args, fmt.Sprintf("--host=%s", d.Url))
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

	args = append(args, d.Folder)

	cmd := exec.Command("debmirror", args...)
	cmd.Stdout = log.StandardLogger().Writer()
	cmd.Stderr = log.StandardLogger().Writer()
	return cmd.Run()
}
