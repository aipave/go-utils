package test_example

import (
    "testing"

    "github.com/alyu01/go-utils/gcast"
    "github.com/alyu01/go-utils/gtoken"
)

func TestGtoken(t *testing.T) {
    var usrid uint64 = 123456
    token, err := gtoken.CreateToken(usrid)
    if err != nil {
        t.Fatal(err)
    }
    t.Logf("token:%v", token)

    err = gtoken.VerifyToken(gcast.ToUint64(123456), token)
    if err != nil {
        t.Fatal(err)
    }
    t.Logf("valid token:%v", token)
}
