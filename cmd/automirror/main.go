package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/syberalexis/automirror/pkg/core"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var (
	defaultConfigFile = "/etc/automirror/config.yaml"
	configFile        string
)

func main() {
	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "Automirror is a software to download packages from internet into local.")
	app.HelpFlag.Short('h')
	app.Version("0.0.1")
	app.Action(func(c *kingpin.ParseContext) error { core.NewAutomirror(configFile).Run(); return nil })

	// Flags
	app.Flag("config", "Config file").Default(defaultConfigFile).Short('c').StringVar(&configFile)

	// Parsing
	args, err := app.Parse(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}
