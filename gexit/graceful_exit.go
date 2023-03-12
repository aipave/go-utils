package gexit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var sig = make(chan os.Signal, 1)
var blockChan = make(chan int)
var releaseHandlers []func()
var closeHandlers []func()

func init() {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go wait()
}

// Close
func Close(handler func()) {
	closeHandlers = append(closeHandlers, handler)
}

// Release
func Release(handler func()) {
	releaseHandlers = append(releaseHandlers, handler)
}

// Wait
func Wait() {
	<-blockChan
}

func wait() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	select {
	case <-sig:
		for _, handler := range closeHandlers {
			if handler != nil {
				handler() // close all
			}
		}
	}

	// exit
	logrus.Infof("received term signal, process will exit after 3 seconds\n")

	time.Sleep(3 * time.Second)

	for _, handler := range releaseHandlers {
		if handler != nil {
			handler() // release all
		}
	}

	blockChan <- 0
	os.Exit(0)
}
