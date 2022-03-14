package pkg

import (
	"ipashare/pkg/conf"
	"ipashare/pkg/log"
)

// Reset all
func Reset() error {
	if err := conf.Reset(); err != nil {
		return err
	}
	log.Reset(conf.Log.Level)
	return nil
}
