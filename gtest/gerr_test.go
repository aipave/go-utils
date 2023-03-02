package gtest

import (
    "testing"

    "github.com/alyu01/go-utils/gerr"
)

func TestGerrNew(t *testing.T) {
    var ErrCodeNotImpl int32 = 501
    ErrCodeNotImplMsg := "The server does not support the requested feature and cannot complete the request"
    err := gerr.NewErrInst(ErrCodeNotImpl, ErrCodeNotImplMsg)
    t.Logf("%v", err.Error())

}
