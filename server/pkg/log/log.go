package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zap.NewAtomicLevel()

// Setup 日志设置
func Setup(lev string) {
	setLevel(lev)
	zap.ReplaceGlobals(zap.New(
		zapcore.NewTee(zapcore.NewCore(
			getConsoleEncoder(),
			zapcore.Lock(os.Stdout),
			level,
		)),
		zap.AddCaller(),
	))
}

var levels = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

// setLevel 日志级别设置
func setLevel(lev string) {
	if l, ok := levels[lev]; ok {
		level.SetLevel(l)
	} else {
		level.SetLevel(zapcore.DebugLevel)
	}
}

// getConsoleEncoder 控制台日志格式
func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Reset(lev string) {
	setLevel(lev)
}
