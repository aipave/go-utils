package glogrus

import (
    "github.com/alyu01/go-utils/glogs/glogrotate"
    "github.com/natefinch/lumberjack"
    "github.com/sirupsen/logrus"
)

// NewLogger
func NewLogger(opts ...LogOption) (l *logrus.Logger) {
    var c config
    for _, fn := range opts {
        fn(&c)
    }

    l = logrus.New()
    l.AddHook(NewCtxHook())
    l.SetLevel(logrus.DebugLevel)
    l.SetReportCaller(true)
    l.SetFormatter(defaultFormatter)
    l.SetOutput(glogrotate.NewWriter(&lumberjack.Logger{
        Filename:   c.filename,
        MaxSize:    256, // 256M
        MaxAge:     30,  //
        MaxBackups: 300, //
        LocalTime:  true,
        Compress:   true,
    }))

    for _, hook := range c.hooks {
        l.AddHook(hook)
    }
    if c.lumLogger != nil {
        l.SetOutput(glogrotate.NewWriter(c.lumLogger))
    }

    if c.logLevel != nil {
        l.SetLevel(*c.logLevel)
    }

    if c.reportCaller != nil {
        l.SetReportCaller(*c.reportCaller)
    }

    if c.formatter != nil {
        l.SetFormatter(c.formatter)
    }

    return
}
