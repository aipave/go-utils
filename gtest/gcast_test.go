package gtest

import (
	"testing"

	"github.com/alyu01/go-utils/gcast"
	_ "github.com/spf13/cast"
)

func TestGcast(t *testing.T) {
	t.Log(gcast.CastToString(""))
}
