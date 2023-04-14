package gimage

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

func Resize(inDirPath string, resizeWidth int, resizeHeight int, resizeWidthFactor int, resizeHeightFactor int) error {
	outDirPath := filepath.Join(inDirPath, fmt.Sprintf("resized%vx%v", resizeWidth, resizeHeight))

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		logrus.Errorf("Error:", err)
		return err
	}

	// Get a list of all PNG files in the input directory
	files, err := ioutil.ReadDir(inDirPath)
	if err != nil {
		logrus.Errorf("Error:", err)
		return err
	}

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()
			// Open the image file
			if !strings.Contains(file.Name(), ".png") {
				return
			}
			img := gocv.IMRead(filepath.Join(inDirPath, file.Name()), gocv.IMReadColor)
			defer img.Close()
			if img.Empty() {
				logrus.Warningf("Failed to read image: %v\n", file)
				return
			}

			// Resize the image
			resized := gocv.NewMat()
			defer resized.Close()
			/*
				InterpolationNearestNeighbor: This option uses the nearest neighbor algorithm for image resizing. It is the fastest but produces the lowest quality results.
				InterpolationLinear: This option uses the bilinear algorithm for image resizing. It is faster than cubic and produces better quality results than nearest neighbor.
				InterpolationCubic: This option uses the bicubic algorithm for image resizing. It produces higher quality results than bilinear but takes longer to compute.
			*/
			gocv.Resize(img, &resized, image.Point{
				X: resizeWidth,
				Y: resizeHeight,
			}, float64(resizeWidthFactor), float64(resizeHeightFactor), gocv.InterpolationNearestNeighbor)

			// Save the resized image
			if ok := gocv.IMWrite(filepath.Join(outDirPath, file.Name()), resized); !ok {
				logrus.Errorf("Failed to write image: %v\n", file)
			}
			return

		}(file)
	}
	wg.Wait()
	return nil

}
