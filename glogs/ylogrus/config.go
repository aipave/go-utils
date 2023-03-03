package ylogrus

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type config struct {
	lumLogger    *lumberjack.Logger
	logLevel     *logrus.Level
	reportCaller *bool
	formatter    logrus.Formatter
	hooks        []logrus.Hook
	maxAge       int64
	filename     string
}

type LogOption func(*config)

// WithLumLogger 设置日志轮转
func WithLumLogger(lum *lumberjack.Logger) LogOption {
	return func(c *config) {
		c.lumLogger = lum
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(level logrus.Level) LogOption {
	return func(c *config) {
		c.logLevel = &level
	}
}

// WithReportCaller 打印调用者
func WithReportCaller(v bool) LogOption {
	return func(c *config) {
		c.reportCaller = &v
	}
}

// WithFormatter 设置日志格式
func WithFormatter(format logrus.Formatter) LogOption {
	return func(c *config) {
		c.formatter = format
	}
}

func WithHooks(hook ...logrus.Hook) LogOption {
	return func(c *config) {
		c.hooks = hook
	}
}

// WithMaxAge 设置日志最长存储天数
func WithMaxAge(days int64) LogOption {
	return func(c *config) {
		c.maxAge = days
	}
}

func WithFileName(filename string) LogOption {
	return func(c *config) {
		c.filename = filename
	}
}
