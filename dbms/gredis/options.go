package gredis

import (
    "time"

    "github.com/alyu01/go-utils/gcast"
)

type RedisConfig struct {
    Addr         string        `json:"Addr"         yaml:"Addr"`
    PoolSize     int           `json:"PoolSize"     yaml:"PoolSize"`
    ReadTimeout  int64         `json:"ReadTimeout"  yaml:"ReadTimeout"`
    WriteTimeout int64         `json:"WriteTimeout" yaml:"WriteTimeout"`
    Password     string        `json:"Password"     yaml:"Password"`
    DB           int           `json:"DB"           yaml:"DB"`
}

var defaultConfig = RedisConfig{
    Addr:         "",
    PoolSize:     100,
    ReadTimeout:  gcast.ToInt64(10 * time.Second),
    WriteTimeout: gcast.ToInt64(20 * time.Second),
    DB:           0,
}

type Option func(c *RedisConfig)

func SetDefaultConfig(rc RedisConfig) {
    if len(rc.Addr) > 0 {
        defaultConfig.Addr = rc.Addr
    }

    if rc.PoolSize > 0 {
        defaultConfig.PoolSize = rc.PoolSize
    }

    if rc.ReadTimeout > 0 {
        defaultConfig.ReadTimeout = rc.ReadTimeout
    }

    if rc.WriteTimeout > 0 {
        defaultConfig.WriteTimeout = rc.WriteTimeout
    }

    if len(rc.Password) > 0 {
        defaultConfig.Password = rc.Password
    }

    if rc.DB > 0 {
        defaultConfig.DB = rc.DB
    }
}

func WithAddr(addr string) Option {
    return func(c *RedisConfig) {
        c.Addr = addr
    }
}

func WithPoolSize(n int) Option {
    return func(c *RedisConfig) {
        c.PoolSize = n
    }
}

func WithReadTimeOut(n time.Duration) Option {
    return func(c *RedisConfig) {
        c.ReadTimeout = int64(n)
    }
}

func WithWriteTimeOut(n time.Duration) Option {
    return func(c *RedisConfig) {
        c.WriteTimeout = int64(n)
    }
}

func WithPassword(p string) Option {
    return func(c *RedisConfig) {
        c.Password = p
    }
}

func WithDB(db int) Option {
    return func(c *RedisConfig) {
        c.DB = db
    }
}

func WithDefaultConfig(config RedisConfig) Option {
    return func(c *RedisConfig) {
        *c = config
    }
}
