package ginfos

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type _Runtime struct {
}

// Runtime
var Runtime _Runtime

// Exec name
func (_Runtime) Exec() string {
	pwd, _ := os.Executable()
	_, exec := filepath.Split(pwd)
	return exec
}

// Pwd exe path
func (_Runtime) Pwd() string {
	pwd, _ := os.Executable()
	pwd, _ = filepath.Split(pwd)
	return pwd
}

// Func name
func (_Runtime) Func() string {
	const unknown = "Unknown_Func"
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return unknown
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return unknown
	}

	return f.Name()
}

func FuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	funcName := strings.ReplaceAll(f.Name(), "main.", "")
	arr := strings.Split(funcName, ".")
	if len(arr) >= 1 {
		return fmt.Sprintf("%s", arr[len(arr)-1])
	}
	return funcName
}

func CallerFuncName(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	arr := strings.Split(name, ".")
	if len(arr) >= 1 {
		return fmt.Sprintf("[%s-%d]", arr[len(arr)-1], line)
	}
	return fmt.Sprintf("[%s-%d]", name, line)
}

// IP
func (_Runtime) IP() string {
	ips, _ := net.InterfaceAddrs()

	for _, ip := range ips {
		if ipAddr, ok := ip.(*net.IPNet); ok && !ipAddr.IP.IsLoopback() && ipAddr.IP.To4() != nil {
			return ipAddr.IP.To4().String()
		}
	}

	return ""
}

func (_Runtime) IsPrivateIP(ip string) bool {
	addr := net.ParseIP(ip)
	return addr.IsPrivate()
}
