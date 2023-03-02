package test

import (
	"testing"

	gcast "github.com/alyu01/go-utils/gcast"
)

func TestGcast(t *testing.T) {
	t.Log(gcast.ToBool("0"))
	t.Log(gcast.ToBool(1))
	t.Log(gcast.ToString("11"))
	t.Log(gcast.ToInt32("12"))
	t.Log(gcast.ToUint64("12.000"))
}
