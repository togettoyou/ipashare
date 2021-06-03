package logger

import (
	"github.com/natefinch/lumberjack"
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
		getLevel(),
	)
	cores = append(cores, consoleCore)

	if conf.Config.LogConfig.IsFile {
		// 创建低于Error级别的文件Core记录器接口
		lowPriorityFileCore := zapcore.NewCore(
			getFileEncoder(),
			getFileWriter(
				conf.Config.LogConfig.FilePath,
				conf.Config.LogConfig.MaxSize,
				conf.Config.LogConfig.MaxBackups,
				conf.Config.LogConfig.MaxAge,
			),
			getLowPriorityLevel(),
		)

		// 创建高于等于Error级别的文件Core记录器接口
		highPriorityFileCore := zapcore.NewCore(
			getFileEncoder(),
			getFileWriter(
				conf.Config.LogConfig.ErrFilePath,
				conf.Config.LogConfig.MaxSize,
				conf.Config.LogConfig.MaxBackups,
				conf.Config.LogConfig.MaxAge,
			),
			getHighPriorityLevel(),
		)

		cores = append(cores, lowPriorityFileCore, highPriorityFileCore)
	}

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

var levels = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLevel 日志模式
func getLevel() zapcore.Level {
	if level, ok := levels[conf.Config.LogConfig.Level]; ok {
		return level
	} else {
		return zapcore.DebugLevel
	}
}

func getHighPriorityLevel() zap.LevelEnablerFunc {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= getLevel()
	})
	return highPriority
}

func getLowPriorityLevel() zap.LevelEnablerFunc {
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= getLevel()
	})
	return lowPriority
}

// getFileEncoder 文件日志格式
func getFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getConsoleEncoder 控制台日志格式
func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Level大写且带颜色的输出
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      // 简短的Call显示
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getFileWriter 文件日志保存配置
// filename 文件保存路径
// maxSize 最大空间，单位MB
// maxBackup 保留的最大旧日志文件数
// maxAge 保留旧日志文件的最大天数
func getFileWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Reset() {
	Setup()
}
