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
	startMirrorName   string
	statusMirrorName  string
	stopMirrorName    string
	restartMirrorName string
)

func main() {
	automirror := &core.Automirror{
		ConfigFile: defaultConfigFile,
	}

	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "Automirror is a software to download packages from internet into local.")
	app.HelpFlag.Short('h')
	app.Version("0.0.1")
	app.Action(func(c *kingpin.ParseContext) error { automirror.Init(); return nil })

	// Flags
	app.Flag("config", "Config file").Default(defaultConfigFile).Short('c').StringVar(&automirror.ConfigFile)

	// Start command
	app.Command("start", "Commands to start mirrors").
		Action(func(c *kingpin.ParseContext) error { automirror.Start(); return nil })

	// Status command
	app.Command("status", "Print Status message").
		Action(func(c *kingpin.ParseContext) error { automirror.Status(); return nil })

	// Stop command
	app.Command("stop", "Commands to stop mirrors").
		Action(func(c *kingpin.ParseContext) error { automirror.Stop(); return nil })

	// Restart command
	app.Command("restart", "Commands to restart mirrors").
		Action(func(c *kingpin.ParseContext) error { automirror.Restart(); return nil })

	// Mirror commands
	mirrorCommand := app.Command("mirror", "Commands to manage a mirror")

	// Mirrors list
	mirrorCommand.Command("list", "List mirrors process").
		Action(func(c *kingpin.ParseContext) error { return automirror.GetMirrors() })

	// Mirror start command
	mirrorStartCommand := mirrorCommand.Command("start", "Start a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.StartMirror(startMirrorName) })
	mirrorStartCommand.Arg("name", "mirror name").Required().StringVar(&startMirrorName)

	// Mirror status command
	mirrorStatusCommand := mirrorCommand.Command("status", "Status of a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.StatusMirror(statusMirrorName) })
	mirrorStatusCommand.Arg("name", "mirror name").Required().StringVar(&statusMirrorName)

	// Mirror stop command
	mirrorStopCommand := mirrorCommand.Command("stop", "Stop a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.StopMirror(stopMirrorName) })
	mirrorStopCommand.Arg("name", "mirror name").Required().StringVar(&stopMirrorName)

	// Mirror restart command
	mirrorRestartCommand := mirrorCommand.Command("restart", "Restart a mirror process").
		Action(func(c *kingpin.ParseContext) error { return automirror.RestartMirror(restartMirrorName) })
	mirrorRestartCommand.Arg("name", "mirror name").Required().StringVar(&restartMirrorName)

	// Parsing
	args, err := app.Parse(os.Args[1:])
	defer automirror.Destroy()
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}
