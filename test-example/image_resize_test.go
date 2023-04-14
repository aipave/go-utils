package test_example

import (
	"testing"

	"github.com/aipave/go-utils/gimage"
)

func TestResize(t *testing.T) {
	dirPath := "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/"

	resizeWidth := 120
	resizeHeight := 120

	gimage.Resize(dirPath, resizeWidth, resizeHeight, 0, 0)

}
