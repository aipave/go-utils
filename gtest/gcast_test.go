package gtest

import (
    "testing"

    _ "github.com/alyu01/go-utils/gcast"
    "github.com/spf13/cast"
)

func TestGcast(t *testing.T) {
    t.Log(cast.CastToString(""))
}
