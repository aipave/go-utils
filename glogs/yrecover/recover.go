package yrecover

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// Recover catch panic error and trigger alert
func Recover(afterRecover ...func()) {
	if err := recover(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "panic: %v\n\n", err)
		buf := make([]byte, 1024)
		for {
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				buf = buf[:n]
				_, _ = os.Stderr.Write(buf)
				break
			}
			buf = make([]byte, 2*len(buf))
		}

		time.Sleep(500 * time.Millisecond)
		_, _ = fmt.Fprintf(os.Stderr, "progress started at: ---------%v-----------\n", time.Now().Format("2006-01-02 15:04:06"))

		// call after recovered
		for _, fn := range afterRecover {
			fn()
		}
	}
}
