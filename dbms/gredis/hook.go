package gredis

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aipave/go-utils/ginfos"
	"github.com/go-redis/redis/v8"
)

type startKey struct{}

type customRedisHook struct {
}

func (hook *customRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *customRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	frame := getCaller()
	if frame == nil {
		return nil
	}

	// todo prometheus
	_ = getFunc(filepath.Base(frame.Function))

	if start, ok := ctx.Value(startKey{}).(time.Time); ok {
		// todo prometheus
		_ = time.Since(start)
	}

	if isActualErr(cmd.Err()) {
		// todo prometheus

	} else {
		// todo prometheus
	}

	return nil
}

func (hook *customRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *customRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if err := hook.AfterProcess(ctx, redis.NewCmd(ctx, "pipeline")); err != nil {
		return err
	}

	frame := getCaller()
	if frame == nil {
		return nil
	}

	// todo prometheus
	_ = getFunc(filepath.Base(frame.Function))

	if start, ok := ctx.Value(startKey{}).(time.Time); ok {
		// todo prometheus
		_ = time.Since(start)
		// todo prometheus
	}

	for _, cmd := range cmds {
		if isActualErr(cmd.Err()) {
			// todo prometheus
		} else {
			// todo prometheus
		}
	}

	return nil
}

func isActualErr(err error) bool {
	return err != nil && err != redis.Nil
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	const maximumCallerDepth = 25

	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		// logrus.Info(f.File, f.Function)
		// According to the file name, match the file with the same name under the first a service and return
		if strings.Contains(f.Function, ginfos.Runtime.Exec()) {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getFunc(name string) string {
	if fields := strings.Split(name, "."); len(fields) > 0 {
		return "redis-" + fields[len(fields)-1]
	}
	return "redis-" + name
}
