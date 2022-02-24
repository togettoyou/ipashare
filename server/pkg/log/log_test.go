package log

import (
	"testing"
)

func TestZap(t *testing.T) {
	Setup("debug")

	logger := New("user")

	logger.Named("api").L().Info("日志")
	logger.Named("svc").L().Info("日志")

}
