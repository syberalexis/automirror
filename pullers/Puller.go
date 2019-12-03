package pullers

import (
	log "github.com/sirupsen/logrus"
)

// Puller interface to expose methods for pulling processes
type Puller interface {
	Pull(log *log.Logger) (int, error)
}
