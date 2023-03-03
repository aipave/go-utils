package ylogrus

import (
    "context"
    "fmt"
    "path/filepath"
    "runtime"
    "sync"
    "time"

    "github.com/natefinch/lumberjack"

    "github.com/alyu01/go-utils/ginfos"
    "github.com/alyu01/go-utils/glogs/ylogrotate"
    "github.com/alyu01/go-utils/glogs/ypanic"

    nested "github.com/antonfisher/nested-logrus-formatter"

    "github.com/sirupsen/logrus"
)

type timeCost struct{}

func fullPathCallerFormatter(frame *runtime.Frame) string {
    return fmt.Sprintf(" [%v:%v %v]", frame.File, frame.Line, filepath.Base(frame.Function))
}

func shortPathCallerFormatter(frame *runtime.Frame) string {
    return fmt.Sprintf(" [%v:%v %v]", filepath.Base(frame.File), frame.Line, filepath.Base(frame.Function))
}

var once sync.Once

var defaultFormatter = &nested.Formatter{
    TimestampFormat:       "2006-01-02 15:04:05.000",
    ShowFullLevel:         true,
    CallerFirst:           true,
    CustomCallerFormatter: fullPathCallerFormatter,
}

func init() {
    Init()
}

// Init 初始化log rus
func Init(opts ...LogOption) {
    var cfg config
    for _, fn := range opts {
        fn(&cfg)
    }

    once.Do(func() {
        ypanic.Redirect("panic.log") // 重定向panic日志

        logrus.SetOutput(ylogrotate.NewWriter(&lumberjack.Logger{
            Filename:   "log/" + ginfos.Runtime.Exec() + ".log",
            MaxSize:    200, // 单个文件最大200M
            MaxAge:     30,  // 最长30天
            MaxBackups: 300, // 最大300个文件
            LocalTime:  true,
            Compress:   true,
        }))

        logrus.AddHook(NewCtxHook())
        logrus.SetLevel(logrus.DebugLevel)
        logrus.SetReportCaller(true)
        logrus.SetFormatter(defaultFormatter)
    })

    if len(cfg.filename) > 0 {
        logrus.SetOutput(ylogrotate.NewWriter(&lumberjack.Logger{
            Filename:   cfg.filename,
            MaxSize:    200, // 单个文件最大200M
            MaxAge:     30,  // 最长30天
            MaxBackups: 300, // 最大300个文件
            LocalTime:  true,
            Compress:   true,
        }))
    }

    if cfg.lumLogger != nil {
        logrus.SetOutput(ylogrotate.NewWriter(cfg.lumLogger))
    }
}

func WithFields(fields logrus.Fields) *logrus.Entry {
    ctx := context.WithValue(context.Background(), timeCost{}, time.Now())
    return logrus.WithContext(ctx).WithFields(fields)
}

func WithCtx(ctx context.Context) *logrus.Entry {
    ctx = context.WithValue(ctx, timeCost{}, time.Now())
    return logrus.WithContext(ctx)
}

// MustSetLevel 设置日志级别
func MustSetLevel(level string) {
    logLevel, err := logrus.ParseLevel(level)
    if err != nil {
        panic(err)
    }
    logrus.SetLevel(logLevel)
}

func SetShortLogPath() {
    defaultFormatter.CustomCallerFormatter = shortPathCallerFormatter
    logrus.SetFormatter(defaultFormatter)
}

type ctxHook struct{}

func NewCtxHook() logrus.Hook {
    return ctxHook{}
}

func (c ctxHook) Levels() []logrus.Level {
    return []logrus.Level{
        logrus.PanicLevel,
        logrus.FatalLevel,
        logrus.ErrorLevel,
        logrus.WarnLevel,
        logrus.InfoLevel,
        logrus.DebugLevel,
        logrus.TraceLevel,
    }
}

func (c ctxHook) Fire(entry *logrus.Entry) error {
    if entry.Context == nil {
        return nil
    }

    if uid := entry.Context.Value("uid"); uid != nil {
        entry.Data["uid"] = uid
    }

    if cost := entry.Context.Value(timeCost{}); cost != nil {
        begin, ok := cost.(time.Time)
        if ok {
            entry.Data["cost"] = time.Since(begin)
        }
    }

    return nil
}
