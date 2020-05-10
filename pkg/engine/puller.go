package engine

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/model"
)

type Puller interface {
	Engine
	Pull(archive model.Archive) error
}

type AbstractPuller struct {
	Puller
	Archives []model.Archive
	Logger   *log.Logger
}

func (puller *AbstractPuller) Run() error {
	for _, archive := range puller.Archives {
		puller.Pull(archive)
	}
	return nil
}
