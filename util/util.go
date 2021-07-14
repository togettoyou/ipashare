package util

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"super-signature/model"
	"super-signature/util/conf"
	"super-signature/util/logger"
)

// Reset 重设配置
func Reset() {
	zap.L().Info("Hot reload config.")
	conf.Reset()
	logger.Reset()
	model.Reset()
	resetGinMode()
}

const (
	debugMode   string = "debug"
	releaseMode string = "release"
	testMode    string = "test"
)

var modes = map[string]string{
	"debug":   debugMode,
	"release": releaseMode,
	"test":    testMode,
}

func resetGinMode() {
	if mode, ok := modes[conf.Config.Server.RunMode]; ok {
		gin.SetMode(mode)
	} else {
		gin.SetMode(debugMode)
	}
}
