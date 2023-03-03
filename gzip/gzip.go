package gzip

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "io"
)

type _GZip int

// GZip gzip
var GZip _GZip

// Compress
func (_GZip) Compress(content string) (string, error) {
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)

    _, err := zw.Write([]byte(content))
    if err != nil {
        return content, err
    }

    if err = zw.Close(); err != nil {
        return content, err
    }

    return buf.String(), nil
}

// UnCompress
func (_GZip) UnCompress(content string) (string, error) {
    b := []byte(content)

    var buf = bytes.NewBuffer(b)
    zr, err := gzip.NewReader(buf)
    if err != nil {
        return content, err
    }

    var out bytes.Buffer
    if _, err = io.Copy(&out, zr); err != nil {
        return content, err
    }

    if err = zr.Close(); err != nil {
        return content, err
    }

    return out.String(), nil
}

// Base64Compress
func (_GZip) Base64Compress(content string) (string, error) {
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)

    _, err := zw.Write([]byte(content))
    if err != nil {
        return content, err
    }

    if err = zw.Close(); err != nil {
        return content, err
    }

    return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Base64UnCompress
func (_GZip) Base64UnCompress(content string) (string, error) {
    b, err := base64.StdEncoding.DecodeString(content)
    if err != nil {
        return content, err
    }

    var buf = bytes.NewBuffer(b)
    zr, err := gzip.NewReader(buf)
    if err != nil {
        return content, err
    }

    var out bytes.Buffer
    if _, err = io.Copy(&out, zr); err != nil {
        return content, err
    }

    if err = zr.Close(); err != nil {
        return content, err
    }

    return out.String(), nil
}
