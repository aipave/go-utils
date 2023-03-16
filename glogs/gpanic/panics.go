//go:build !windows

package gpanic

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aipave/go-utils/gexit"
	"github.com/aipave/go-utils/gtime"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var offset int64

var PAlertMgr *AlertMgr = &AlertMgr{}

type AlertMgr struct {
	AlertUrl string
}

// Redirect
func (a *AlertMgr) Redirect(filename string) {
	once.Do(func() {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}

		_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd())) // Redirect to panic file

		// read last block
		var data []byte
		data, err = ioutil.ReadAll(f)
		if err == nil {
			splits := regexp.MustCompile("progress started at: .*-------").Split(string(data), -1)
			if len(splits) > 0 {
				a.filter([]byte(splits[len(splits)-1]))
			}
		}

		offset = int64(len(data))

		// set current block begin
		_, _ = fmt.Fprintf(f, "progress started at: ---------%v-----------\n", time.Now().Format(gtime.FormatDefault))

		go a.watch(filename, f)
	})
}

func (a *AlertMgr) watch(filename string, f *os.File) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("new wacher err:%v", err)
		return
	}

	err = watcher.Add(filename)
	if err != nil {
		logrus.Errorf("watch file:%v err:%v", filename, err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	gexit.Close(func() {
		cancel()
	})

	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Write {
				_, _ = f.Seek(offset, io.SeekStart)
				data, _ := ioutil.ReadAll(f)
				offset += int64(len(data))

				a.filter(data)
			}

		case err = <-watcher.Errors:
			logrus.Errorf("watch file:%v err:%v", filename, err)
			return

		case <-ctx.Done():
			_ = watcher.Remove(filename)
			_ = watcher.Close()
			return
		}
	}
}

var silentMap = make(map[string]bool) // silence policy

func (a *AlertMgr) filter(buf []byte) {
	var matched bool
	var count int
	var stack []string
	for _, line := range strings.Split(string(buf), "\n") {
		if !strings.HasPrefix(line, "panic:") && !matched {
			continue
		}

		if strings.HasPrefix(line, "panic:") {
			matched = true
		}

		if count >= 10 {
			break
		}

		stack = append(stack, line)
	}

	if len(stack) > 0 {
		alertMsg := strings.Join(stack, "\n")
		silentKey := fmt.Sprintf("%v:%v", time.Now().Minute(), localMd5(alertMsg)) // 重复的内容静默一分钟
		if silentMap[silentKey] {
			return
		}

		silentMap[silentKey] = true
		a.triggerAlert(buildAlert(alertMsg))
	}
}

func localMd5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}
