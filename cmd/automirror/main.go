package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var (
	app    = kingpin.New(filepath.Base(os.Args[0]), "Automirror is a software to download packages from internet into local.")
	config = app.Flag("config", "Config file").Short('c').String()
)

func main() {
	app.HelpFlag.Short('h')
	app.Version("0.0.1")

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		fmt.Fprintln(os.Stdout, *config)
	}
}
