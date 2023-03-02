package ginfos

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
)

// NilStruct if v==nil, return a struct
func NilStruct(v any) any {
	if v != nil {
		return v
	}

	return struct{}{}
}

// NilSlice if v==nil, return a struct
func NilSlice(v any) any {
	if v != nil {
		return v
	}

	return []interface{}{}
}

func JsonStr(v any) string {
	jsonStr, _ := json.Marshal(v)
	return string(jsonStr)
}

// Md5 md5 sign
func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}
