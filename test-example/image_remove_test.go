package test_example

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/aipave/go-utils/gimage"
)

func TestRemove(t *testing.T) {
	t.Log(gimage.RemoveBackground("e2EmV9bShqJ41us1a2PX6aAK", "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/"))
}

func TestExt(t *testing.T) {
	filename := "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/01.png"
	filename2 := "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/gifs"
	splitArr := strings.Split(filename, "/")
	filenameNew := splitArr[len(splitArr)-1]
	t.Log(filepath.Ext(filenameNew))
	t.Log(filepath.Join(filename2, filenameNew+".gifs"))
	t.Log(strings.TrimSuffix(filename, filenameNew) + "togif/" +
		strings.TrimSuffix(filenameNew, filepath.Ext(filenameNew)) +
		".gif")
}
