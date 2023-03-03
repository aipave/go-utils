package test_example

import (
    "testing"

    "github.com/alyu01/go-utils/gerr"
)

func TestGerr(t *testing.T) {
    var ErrCodeNotImpl int32 = 501
    ErrCodeNotImplMsg := "The server does not support the requested feature and cannot complete the request"
    err := gerr.New(ErrCodeNotImpl, ErrCodeNotImplMsg)
    t.Log(err.Error())
    t.Log(gerr.Append(err, gerr.New(ErrCodeNotImpl, ErrCodeNotImplMsg)))

}
