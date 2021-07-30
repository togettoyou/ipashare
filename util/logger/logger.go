package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"super-signature/util/conf"
)

// Setup 设置日志
func Setup() {
	// 所有的Core记录器接口
	cores := make([]zapcore.Core, 0)

	// 创建控制台Core记录器接口
	consoleCore := zapcore.NewCore(
		getConsoleEncoder(),
		zapcore.Lock(os.Stdout),
		getLevel(conf.Config.Model),
	)
	cores = append(cores, consoleCore)

	// Option可选配置
	options := []zap.Option{
		zap.AddCaller(),                   // Caller调用显示
		zap.AddStacktrace(zap.ErrorLevel), // 堆栈跟踪级别
	}

	// 构建一个 zap 实例
	logger := zap.New(
		zapcore.NewTee(cores...),
		options...,
	)

	zap.ReplaceGlobals(logger) // 设为全局zap实例
}

// getLevel 日志模式
func getLevel(mode string) zapcore.Level {
	if mode == "release" {
		return zapcore.InfoLevel
	} else {
		return zapcore.DebugLevel
	}
}

// getConsoleEncoder 控制台日志格式
func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Level大写且带颜色的输出
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      // 简短的Call显示
	return zapcore.NewConsoleEncoder(encoderConfig)
}
