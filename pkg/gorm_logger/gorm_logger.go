package gorm_logger

import (
	"LearnJapan.com/configs"
	"LearnJapan.com/pkg/logger"
	"context"
	"fmt"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	cfg    *configs.Configs
	logger *logger.Logger
	level  gormlogger.LogLevel
}

func NewGormLogger(cfg *configs.Configs) *GormLogger {
	return &GormLogger{
		cfg:    cfg,
		logger: logger.NewLogger(cfg),
		level:  gormlogger.LogLevel(1),
	}
}

func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.level = level

	return l
}

func (l GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Info(msg)
}

func (l GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Warn(msg)
}

func (l GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Error(msg)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	sql, rows := fc()
	l.logger.Info(fmt.Sprintf("[%d ms, %d rows] sql -> %s", elapsed.Milliseconds(), rows, sql))
}
