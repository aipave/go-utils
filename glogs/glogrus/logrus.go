package glogrus

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/aipave/go-utils/gtime"
	"github.com/natefinch/lumberjack"

	"github.com/aipave/go-utils/ginfos"
	"github.com/aipave/go-utils/glogs/glogrotate"
	"github.com/aipave/go-utils/glogs/gpanic"

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
	TimestampFormat:       gtime.FormatDefaultMill,
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
		gpanic.Redirect("panic.log", cfg.alertUrl) // 重定向panic日志

		logrus.SetOutput(glogrotate.NewWriter(&lumberjack.Logger{
			Filename:   "log/" + ginfos.Runtime.Exec() + ".log",
			MaxSize:    256, // 256M
			MaxAge:     30,  // 30d
			MaxBackups: 300, // max 300 files
			LocalTime:  true,
			Compress:   true,
		}))

		logrus.AddHook(NewCtxHook())
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(defaultFormatter)
	})

	if len(cfg.filename) > 0 {
		logrus.SetOutput(glogrotate.NewWriter(&lumberjack.Logger{
			Filename:   cfg.filename,
			MaxSize:    256, // 256M
			MaxAge:     30,  // 30day
			MaxBackups: 300, // max 300files
			LocalTime:  true,
			Compress:   true,
		}))
	}

	if cfg.lumLogger != nil {
		logrus.SetOutput(glogrotate.NewWriter(cfg.lumLogger))
	}
}

func WithAlertUrl(url string) LogOption {
	return func(cfg *config) {
		cfg.alertUrl = url
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

// MustSetLevel set log level
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
