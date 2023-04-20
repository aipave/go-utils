package test_example

import (
	"fmt"
	"image/gif"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CuteReimu/neuquant"
	"github.com/aipave/go-utils/gerr"
	"github.com/disintegration/imaging"
)

// low quality
// In this example, the delay parameter is set to 100, which represents the delay between frames in hundredths of a second. You can adjust this value to control the speed of the animation.
func TestFormatConver(t *testing.T) {
	err := PngToGif("/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/", "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/01.png", 95)
	if err != nil {
		t.Fatal(err)
	}

}

func PngToGif(path, filename string, delay int) error {
	outDirPath := filepath.Join(path, "gifs")
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		return gerr.New(0, err.Error())
	}
	img, err := imaging.Open(filename)
	if err != nil {
		return gerr.New(0, filename)
	}
	//palettedImg := image.NewPaletted(img.Bounds(), palette.WebSafe)
	//draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Src)
	//outGif := &gif.GIF{}
	//outGif.Image = append(outGif.Image, palettedImg)
	//outGif.Delay = append(outGif.Delay, delay) // the delay parameter represents the delay between frames in hundredths of a second.
	// Convert the image to a slice of colorful.Color

	/*
		// Quantize the colors using the NeuQuant algorithm
		quantizedColors := colorful.QuantizeNeuQuant(colors, 256)

		// Convert the quantized colors back to a slice of color.Color
		palettedColors := make([]color.Color, len(quantizedColors))
		for i, c := range quantizedColors {
			palettedColors[i] = c
		}

		// Create a new Paletted image using the quantized colors
		palette := color.Palette(palettedColors)
		bounds := img.Bounds()
		palettedImg := image.NewPaletted(bounds, palette)
		index = 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				palettedImg.Set(x, y, palettedColors[index])
				index++
			}
		}
	*/

	splitArr := strings.Split(filename, "/")
	filenameRaw := splitArr[len(splitArr)-1]
	createName := filepath.Join(outDirPath,
		strings.TrimSuffix(filenameRaw, filepath.Ext(filenameRaw))+".gif")
	f, err := os.Create(createName)
	if err != nil {
		return gerr.New(100, createName)
	}
	defer f.Close()
	return gif.Encode(f, img, neuquant.Opt())
	//return gif.Encode(f, img, &gif.Options{
	//	NumColors: 256,
	//	Quantizer: nil,
	//	Drawer:    draw.FloydSteinberg,
	//})
}

// high quality
//
// gifski -r 50  03.png 03.png -o 03.gif --extra
func TestGifski(t *testing.T) {
	// Directory containing PNG images
	dirname := "/Users/ricco/Downloads/01-道场小恶魔-3D版/a0-github/prompt/resized240x240"

	// Output directory for GIF images
	outdir := "/Users/ricco/Downloads/01-道场小恶魔-3D版/a0-github/prompt/resized240x240/gifs"

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
