package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/syberalexis/automirror/pkg/core"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

func main() {
	automirror := &core.Automirror{}

	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "Automirror is a software to download packages from internet into local.")
	app.Flag("config", "Config file").Default("/etc/automirror/config.toml").Short('c').StringVar(&automirror.Config)

	// Mirror commands
	mirrorCommand := app.Command("mirror", "Commands to manage a mirror")
	app.HelpFlag.Short('h')
	app.Version("0.0.1")

	// Mirrors list
	mirrorCommand.Command("list", "List mirrors process").
		Action(func(c *kingpin.ParseContext) error { return automirror.GetMirrors() })

	// Mirror start command
	mirrorStartCommand := mirrorCommand.Command("start", "Start a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.Start(*c.Elements[2].Value) })
	mirrorStartCommand.Arg("name", "mirror name").Required().String()

	// Mirror status command
	mirrorStatusCommand := mirrorCommand.Command("status", "Status of a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.Status(*c.Elements[2].Value) })
	mirrorStatusCommand.Arg("name", "mirror name").Required().String()

	// Mirror stop command
	mirrorStopCommand := mirrorCommand.Command("stop", "Stop a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.Stop(*c.Elements[2].Value) })
	mirrorStopCommand.Arg("name", "mirror name").Required().String()

	// Mirror restart command
	mirrorRestartCommand := mirrorCommand.Command("restart", "Restart a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.Restart(*c.Elements[2].Value) })
	mirrorRestartCommand.Arg("name", "mirror name").Required().String()

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
