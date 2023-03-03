package ylogrus

import (
	"github.com/alyu01/go-utils/glogs/ylogrotate"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// NewLogger 创建单独日志文件
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
    l.SetOutput(ylogrotate.NewWriter(&lumberjack.Logger{
        Filename:   c.filename,
        MaxSize:    200, // 单个文件最大200M
        MaxAge:     30,  // 最长30天
        MaxBackups: 300, // 最大300个文件
        LocalTime:  true,
        Compress:   true,
    }))

    for _, hook := range c.hooks {
        l.AddHook(hook)
    }
    if c.lumLogger != nil {
        l.SetOutput(ylogrotate.NewWriter(c.lumLogger))
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
