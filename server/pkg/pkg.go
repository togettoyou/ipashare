package pkg

import (
	"supersign/pkg/conf"
	"supersign/pkg/log"
)

// Reset all
func Reset() error {
	if err := conf.Reset(); err != nil {
		return err
	}
	log.Reset(conf.Log.Level)
	return nil
}
