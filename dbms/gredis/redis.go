package gredis

import (
    "context"
    "time"

    "github.com/alyu01/go-utils/gracexit"
    redisprom "github.com/globocom/go-redis-prometheus"
    "github.com/go-redis/redis/v8"
    "github.com/sirupsen/logrus"
)

// NewRedisClient initialredis client， goroutine safety, no need to close
// client := NewRedisClient(Addr("host:port"))
// client.SetNX().Result()
func NewRedisClient(opts ...Option) (redisClient *redis.Client) {
    var c = defaultConfig
    for _, fn := range opts {
        fn(&c)
    }

    redisClient = redis.NewClient(&redis.Options{
        Addr:         c.Addr,
        PoolSize:     c.PoolSize,
        ReadTimeout:  time.Duration(c.ReadTimeout) * time.Millisecond, // time.ParseDuration(str+"ms")
        WriteTimeout: time.Duration(c.WriteTimeout) * time.Millisecond,
        Password:     c.Password,
    })

    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        logrus.Fatalf("init redis client err:%v config:%v", err, c)
    }

    gracexit.Release(func() {
        if redisClient == nil {
            return
        }

        if err := redisClient.Close(); err != nil {
            logrus.Errorf("close redis client err:%v", err)
        }

        logrus.Infof("redis addr:%v resource released.", c.Addr)
    })

    redisClient.AddHook(&customRedisHook{})
    redisClient.AddHook(redisprom.NewHook())

    logrus.Infof("init redis client done, config:%+v", c)
    return
}

// ResetNil reset redis.Nil error
func ResetNil(err error) error {
    if err == redis.Nil {
        return nil
    }
    return err
}
