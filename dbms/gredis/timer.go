package gredis

import (
    "context"
    _ "embed"
    "sync"
    "time"

    "github.com/alyu01/go-utils/gcast"
    "github.com/alyu01/go-utils/ginfos"
    "github.com/alyu01/go-utils/gracexit"
    redis "github.com/go-redis/redis/v8"
    //_ "github.com/robfig/cron" //The api of the cron package has changed since version v3, requiring a new import path
    "github.com/robfig/cron/v3"

    "github.com/sirupsen/logrus"
)

//go:embed res/get_lock.lua
var lockScript string

// NewCron timer
// cli redis client
// name of timer(redis key)， global unique
// ttl time, validity period of a key is 3*ttl
func NewCron(cli *redis.Client, name string, ttl int64) *Cron {
    exec := "gl:" + ginfos.Runtime.Exec() // global lock:progress name:name
    if len(name) > 0 {
        exec += ":" + name
    }

    c := &Cron{cli: cli, name: exec, c: cron.New(cron.WithSeconds()), ttl: ttl}
    go c.run()
    return c
}

type Cron struct {
    cli      *redis.Client
    name     string
    c        *cron.Cron
    ttl      int64
    isMaster bool
    mux      sync.Mutex
}

func (c *Cron) AddFunc(spec string, cmd func()) error {
    _, err := c.c.AddFunc(spec, c.wrapperCmd(cmd))
    return err
}

func (c *Cron) Start() {
    c.c.Start()
}

func (c *Cron) Stop() {
    c.c.Stop()
}

func (c *Cron) wrapperCmd(cmd func()) func() {
    return func() {
        if !c.isCurrentMaster() {
            logrus.Warningf("machine:%v not master, abort", ginfos.Runtime.IP())
            return
        }

        // master machine
        cmd()
    }
}

func (c *Cron) run() {
    c.grab()

    ticker := time.Tick(time.Duration(c.ttl) * time.Second)
    ctx, cancel := context.WithCancel(context.Background())
    gracexit.Close(func() {
        cancel()
    })

    for {
        select {
        case <-ticker:
            c.grab()

        case <-ctx.Done():
            logrus.Infof("ctx done, cron.run abort")
            return
        }
    }
}

func (c *Cron) grab() {
    result, err := c.cli.Eval(context.Background(), lockScript, []string{c.name}, ginfos.Runtime.IP(), 3*c.ttl).Result()
    if err != nil {
        c.setSlaver()
        logrus.Errorf("grab lock|name=%v|ttl=%v|isMaster=%v|result=%v|err=%v", c.name, c.ttl, c.isMaster, result, err)
        return
    }

    const (
        resultFail    = 0 // fail to grab lock
        resultSuccess = 1 // success to grab lock
        resultRenewal = 2 // renewal success
    )

    switch gcast.ToInt(result) {
    case resultFail:
        c.setSlaver()
        logrus.Warningf("grab lock fail|name=%v|ttl=%v|isMaster=%v|-", c.name, c.ttl, c.isMaster)

    case resultSuccess:
        c.setMaster()
        logrus.Infof("grab lock success|name=%v|ttl=%v|isMaster=%v|-", c.name, c.ttl, c.isMaster)

    case resultRenewal:
        c.setMaster()
        logrus.Infof("renewal success|name=%v|ttl=%v|isMaster=%v|-", c.name, c.ttl, c.isMaster)
    }
}

func (c *Cron) isCurrentMaster() bool {
    c.mux.Lock()
    defer c.mux.Unlock()
    return c.isMaster
}

func (c *Cron) setMaster() {
    c.mux.Lock()
    c.isMaster = true
    c.mux.Unlock()
}

func (c *Cron) setSlaver() {
    c.mux.Lock()
    c.isMaster = false
    c.mux.Unlock()
}
