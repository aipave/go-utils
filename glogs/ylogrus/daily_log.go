package ylogrus

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "time"

    "github.com/alyu01/go-utils/gstl/options"
    "github.com/alyu01/go-utils/gstl/smap"
    "github.com/alyu01/go-utils/gtime"
    nested "github.com/antonfisher/nested-logrus-formatter"
    "github.com/sirupsen/logrus"
)

const dailyFormat = gtime.FormatDate

var fileMap = smap.NewMap[string, int64](options.WithLocker())

// NewDailyLog 新增单独的日志文件
// yNewDailyLog("svip.log", ylogrus.WithMaxAge(30)) 保留30天日志
func NewDailyLog(filename string, opts ...LogOption) *logrus.Logger {
    var c = config{maxAge: 60}
    for _, fn := range opts {
        fn(&c)
    }

    l := logrus.New()

    l.SetReportCaller(true)
    l.SetFormatter(&nested.Formatter{
        TimestampFormat: "2006-01-02 15:04:05",
        ShowFullLevel:   true,
        CallerFirst:     true,
        CustomCallerFormatter: func(frame *runtime.Frame) string {
            return fmt.Sprintf(" [%v:%v %v]", filepath.Base(frame.File), frame.Line, filepath.Base(frame.Function))
        },
    })

    fileMap.Set(filename, c.maxAge)

    go clean()

    l.SetOutput(&dailyWriter{filename: filename, date: time.Now().Format(dailyFormat)})

    return l
}

func clean() {
    ticker := time.NewTicker(time.Hour)

    for range ticker.C {
        fileMap.Range(func(filename string, maxAge int64) bool {
            for {
                date := time.Now().Add(-time.Duration(maxAge+1) * 86400 * time.Second).Format("2006-01-02") // 2022-11-11
                targetFile := fmt.Sprintf("%v-%v", filename, date)                                          // svip.log-2022-11-11

                _, err := os.Stat(targetFile)
                if os.IsNotExist(err) {
                    break
                }

                if err != nil {
                    fmt.Printf("[ylogs] state logfile:%v err:%v\n", targetFile, maxAge)
                }

                err = os.Remove(targetFile)
                fmt.Printf("[ylogs] remove logfile:%v, maxage:%v, err:%v\n", targetFile, maxAge, err)

                maxAge++
            }

            return true
        })
    }
}

type dailyWriter struct {
    filename string
    date     string
    writer   *os.File
}

func (dw *dailyWriter) openFile() (err error) {
    dw.writer, err = os.OpenFile(fmt.Sprintf("%v-%v", dw.filename, dw.date), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
    return
}

func (dw *dailyWriter) closeFile() (err error) {
    if dw.writer == nil {
        return
    }

    return dw.writer.Close()
}

func (dw *dailyWriter) Write(b []byte) (n int, err error) {
    if dw.writer == nil {
        err = dw.openFile()
        if err != nil {
            return
        }
    }

    if date := time.Now().Format(dailyFormat); date != dw.date {
        err = dw.closeFile()
        if err != nil {
            return
        }

        dw.date = date
        err = dw.openFile()
        if err != nil {
            return
        }
    }

    return dw.writer.Write(b)
}
