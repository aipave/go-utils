package test_example

import (
	"testing"

	"github.com/aipave/go-utils/gimage"
)

func TestRemove(t *testing.T) {
	t.Log(gimage.RemoveBackground("e2EmV9bShqJ41us1a2PX6aAK", "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/"))
}
