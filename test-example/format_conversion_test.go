package gimage

import (
	"os/exec"
	"testing"

	"github.com/aipave/go-utils/gimage"
)

// In this example, the delay parameter is set to 100, which represents the delay between frames in hundredths of a second. You can adjust this value to control the speed of the animation.
func TestFormatConver(t *testing.T) {
	err := gimage.PngToGif("/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/01.png", 100)
	if err != nil {
		t.Fatal(err)
	}

}

// gifski --fps 1 -o filename.gif  ori-filename.png
func TestGifski(t *testing.T) {

	// Encode the GIF using gifski
	cmd := exec.Command("gifski", "--fps", "1", "-o", "output.gif", "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/02.png")
	//cmd.Stdin = palettedImage

	// Create a writer to write the paletted image data to the buffer
	err := cmd.Run()
	if err != nil {
		t.Error(err)
		return
	}

}
