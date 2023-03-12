package glogrus

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
	alertUrl     string
}

type LogOption func(*config)

// WithLumLogger Set Up Log Rotation
func WithLumLogger(lum *lumberjack.Logger) LogOption {
	return func(c *config) {
		c.lumLogger = lum
	}
}

// WithAlerUrl
func WithAlerUrl(url string) LogOption {
	return func(c *config) {
		c.alertUrl = url
	}
}

// WithLogLevel
func WithLogLevel(level logrus.Level) LogOption {
	return func(c *config) {
		c.logLevel = &level
	}
}

// WithReportCaller
func WithReportCaller(v bool) LogOption {
	return func(c *config) {
		c.reportCaller = &v
	}
}

// WithFormatter
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

// WithMaxAge function specifies the maximum length of time is the log file storage,
// rather than the biggest number of log files.
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
