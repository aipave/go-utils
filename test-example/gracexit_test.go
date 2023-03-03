package test_example

import (
    "context"
    "testing"
    "time"

    "github.com/alyu01/go-utils/gracexit"
)

type Test struct {
    log *testing.T
}

func (t *Test) sayHello() {
    t.log.Logf("Hello world!")
}

func (t *Test) handleFoo(ctx context.Context) {
    curTimer := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-ctx.Done():
            t.log.Log("timer close")
        case <-curTimer.C:
            t.sayHello()
        }
    }
}

// received term signal, process will exit after 3 seconds
func (t *Test) testGraceExitNoBreakLoop() {
    ctx, cancel := context.WithCancel(context.Background())

    go t.handleFoo(ctx)

    gracexit.Close(cancel)

}

func sayHello(t *testing.T) {
    t.Log("hello world!")
}

func testGraceExitAddBreakLoop(t *testing.T) {
    timeTick := time.Tick(3 * time.Second)
    ctx, cancel := context.WithCancel(context.Background())
    gracexit.Close(cancel)

    ///> add loop, break loop
Loop:
    for {
        select {
        case <-timeTick:
            sayHello(t)
        case <-ctx.Done():
            t.Log("progress received kill, abort.")
            break Loop

        }
    }

}

func TestGracefulExit(tt *testing.T) {
    t := &Test{
        log: tt,
    }
    t.testGraceExitNoBreakLoop()

    gracexit.Wait() //block
}

func TestGracefulTimeCrontab(t *testing.T) {
    testGraceExitAddBreakLoop(t)

    gracexit.Wait() //block
}

//type Dao struct {
//    closeChan chan struct{}
//}
//var dao *Dao
//
//func test-example() {
//    dao = &Dao{
//        closeChan: make(chan struct{}),
//    }
//
//    gracexit.Close(func() {
//        dao.closeChan <- struct {
//        }{}
//    })
//
//}
