package gimage

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"gocv.io/x/gocv"
)

func TestResize(t *testing.T) {
	dirPath := "/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/"
	outputDir := filepath.Join(dirPath, "resized")

	resizeWidth := 120
	resizeHeight := 120
	resizeWidthFactor := 0
	resizeHeightFactor := 0

	// Get all files in the directory
	//files, err := filepath.Glob(filepath.Join(dirPath, "*.png"))
	//if err != nil {
	//	panic(err)
	//}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	// Get a list of all PNG files in the input directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()
			// Open the image file
			//img := gocv.IMRead(file, gocv.IMReadColor)
			if !strings.Contains(file.Name(), ".png") {
				return
			}
			img := gocv.IMRead(filepath.Join(dirPath, file.Name()), gocv.IMReadColor)
			if img.Empty() {
				fmt.Printf("Failed to read image: %v\n", file)
				return
			}
			defer img.Close()

			// Resize the image
			resized := gocv.NewMat()
			gocv.Resize(img, &resized, image.Point{
				X: resizeWidth,
				Y: resizeHeight,
			}, float64(resizeWidthFactor), float64(resizeHeightFactor), gocv.InterpolationNearestNeighbor)
			/*
				InterpolationNearestNeighbor: This option uses the nearest neighbor algorithm for image resizing. It is the fastest but produces the lowest quality results.
				InterpolationLinear: This option uses the bilinear algorithm for image resizing. It is faster than cubic and produces better quality results than nearest neighbor.
				InterpolationCubic: This option uses the bicubic algorithm for image resizing. It produces higher quality results than bilinear but takes longer to compute.
			*/

			//newFilePath := fmt.Sprintf("%s_resized.png", file[:len(file)-4])
			//gocv.IMWrite(filepath.Join(outputDir, file.Name()), resized, params...)
			// Save the resized image
			// params := []int{gocv.IMWritePngCompression, 9}
			if ok := gocv.IMWrite(filepath.Join(outputDir, file.Name()), resized); !ok {
				fmt.Printf("Failed to write image: %v\n", file)
			}
			defer resized.Close()
			return

		}(file)
	}
	wg.Wait()

}
