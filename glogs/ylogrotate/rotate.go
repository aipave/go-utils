package ylogrotate

import (
    "sync"

    "github.com/robfig/cron"

    "github.com/natefinch/lumberjack"
)

var m sync.Mutex
var loggerMap = make(map[string]*lumberjack.Logger)

// Divide the log every day at 0:00
func init() {
    c := cron.New()
    _ = c.AddFunc("0 0 0 * * *", rotate)
    c.Start()
}

// NewWriter
func NewWriter(logger *lumberjack.Logger) *lumberjack.Logger {
    m.Lock()
    defer m.Unlock()

    if logfile, exist := loggerMap[logger.Filename]; exist {
        _ = logfile.Close()
    }

    loggerMap[logger.Filename] = logger

    return logger
}

func GetWriter(filename string) *lumberjack.Logger {
    m.Lock()
    defer m.Unlock()

    // For the same log file, only one copy is initialized
    if logger, ok := loggerMap[filename]; ok {
        return logger
    }

    // init
    logger := &lumberjack.Logger{
        Filename:   filename,
        MaxSize:    200, // max 200M
        MaxAge:     30,  // max 30day
        MaxBackups: 30,  // max files
        Compress:   true,
    }
    loggerMap[filename] = logger

    return logger
}

func rotate() {
    if len(loggerMap) == 0 {
        return // nothing to do
    }

    for _, logger := range loggerMap {
        _ = logger.Rotate()
    }
}
