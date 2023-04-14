package test_example

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aipave/go-utils/gimage"
)

// low quality
// In this example, the delay parameter is set to 100, which represents the delay between frames in hundredths of a second. You can adjust this value to control the speed of the animation.
func TestFormatConver(t *testing.T) {
	err := gimage.PngsToGifs("/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/", 1)
	if err != nil {
		t.Fatal(err)
	}

}

// high quality
//
// gifski -r 50  03.png 03.png -o 03.gif --extra
func TestGifski(t *testing.T) {
	// Directory containing PNG images
	dirname := "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/"

	// Output directory for GIF images
	outdir := "/Users/ricco/Downloads/01-道场小恶魔-3D版/gif/"

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.MkdirAll(outdir, os.ModePerm)
	}

	// Find all PNG images in the input directory
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate over PNG images and convert each one to a GIF
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg" {
			// Construct input and output filenames
			infile := filepath.Join(dirname, file.Name())
			outfile := filepath.Join(outdir, file.Name()[0:len(file.Name())-4]+".gif")

			// Run gifski command to convert PNG to GIF
			cmd := exec.Command("gifski", "-r", "50", infile, infile, "-o", outfile)
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
