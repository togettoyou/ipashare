package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
	ZapLogger     *zap.Logger
	LogLevel      gormlogger.LogLevel
}

func NewGormLogger(zapLogger *zap.Logger) gormlogger.Interface {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		ZapLogger:     zapLogger,
		LogLevel:      gormlogger.Info,
	}
}

func (gl *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *gl
	newLogger.LogLevel = level
	return &newLogger
}

func (gl GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= gormlogger.Info {
		gl.ZapLogger.Sugar().Infof(msg, data...)
	}
}

func (gl GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= gormlogger.Warn {
		gl.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

func (gl GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= gormlogger.Error {
		gl.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

func (gl GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.LogLevel > gormlogger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && gl.LogLevel >= gormlogger.Error:
			sql, rows := fc()
			if err == gorm.ErrRecordNotFound {
				gl.ZapLogger.Warn(sql, zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows))
			} else {
				gl.ZapLogger.Error(sql, zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows))
			}
		case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.LogLevel >= gormlogger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.SlowThreshold)
			gl.ZapLogger.Warn(sql, zap.String("tips", slowLog), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows))
		case gl.LogLevel == gormlogger.Info:
			sql, rows := fc()
			gl.ZapLogger.Debug(sql, zap.Duration("elapsed", elapsed), zap.Int64("rows", rows))
		}
	}
}
