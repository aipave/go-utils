package gtest

import "testing"

func TestGerrNew(t testing.T) {
    ErrCodeNotImpl := 501
    ErrCodeNotImplMsg := "The server does not support the requested feature and cannot complete the request"
    err := gerr.NewErrInst(ErrCodeNotImpl, ErrCodeNotImplMsg)

}
