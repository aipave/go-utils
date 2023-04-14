package gimage

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aipave/go-utils/gerr"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

func loadImages(path string) ([]*image.NRGBA, error) {
	images := []*image.NRGBA{}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") {
			img, err := imaging.Open(p)
			if err != nil {
				return err
			}
			nrgbaImg := imaging.Clone(img)
			images = append(images, nrgbaImg)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func createGif(images []*image.NRGBA, filename string, delay int) error {
	outGif := &gif.GIF{}
	for _, img := range images {
		palettedImg := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Src)
		outGif.Image = append(outGif.Image, palettedImg)
		outGif.Delay = append(outGif.Delay, delay)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, outGif)
}

func PngsToGif(path string, delay int, dest_filename string) error {
	images, err := loadImages(path)
	if err != nil {
		return err
	}
	err = createGif(images, fmt.Sprintf("%v.gif", dest_filename), delay)
	if err != nil {
		return err
	}
	return nil

}

// load a single PNG or JPG image and convert it to a GIF
func PngToGif(filename string, delay int) error {
	img, err := imaging.Open(filename)
	if err != nil {
		return err
	}
	palettedImg := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Src)
	outGif := &gif.GIF{}
	for i := 0; i < delay; i++ {
		outGif.Image = append(outGif.Image, palettedImg)
	}
	outGif.Delay = append(outGif.Delay, delay) // the delay parameter represents the delay between frames in hundredths of a second.
	f, err := os.Create(strings.TrimSuffix(filename, filepath.Ext(filename)) + ".gif")
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, outGif)
}

// load all the PNG or JPG images from a specified directory and convert each one to a GIF
func PngsToGifs(path string, delay int) error {
	// recursively walks the directory tree and calls the provided function for each file or subdirectory it encounters.
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") {
			err := PngToGif(p, delay)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func getImageFromUrl(url string) (bytes.Buffer, error) {
	response, err := http.Get(url)
	if err != nil {
		return bytes.Buffer{}, gerr.New(0, "http get error")
	}
	defer response.Body.Close()

	// Decode GIF image
	gifImage, err := gif.DecodeAll(response.Body)
	if err != nil {
		return bytes.Buffer{}, gerr.New(0, "gif decode error")
	}

	contentTp := response.Header.Get("Content-Type") // image/jpg image/png
	// Convert GIF image to JPEG format
	var buffer bytes.Buffer
	randFrame := rand.Int() % len(gifImage.Image)
	switch {
	case strings.Contains(contentTp, "jpg"), strings.Contains(contentTp, "jpeg"):
		err = jpeg.Encode(&buffer, gifImage.Image[randFrame], nil)
	case strings.Contains(contentTp, "png"):
		err = png.Encode(&buffer, gifImage.Image[randFrame])

	}
	if err != nil {
		return bytes.Buffer{}, gerr.New(0, "gif decode error")
	}

	return buffer, nil

}

func DynamicToStatic(url string, destFormat string) error {
	buffer, err := getImageFromUrl(url)
	if err != nil {
		return err
	}
	// Write buffer to file
	ext := filepath.Ext(url)
	filename := strings.TrimSuffix(filepath.Base(url), ext)
	file, err := os.Create(fmt.Sprintf("%v.%v", filename, destFormat))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = buffer.WriteTo(file)
	if err != nil {
		return err
	}
	return nil

}

func GifsToPng(path string) error {
	// Get a list of all the GIF images in the directory
	gifFileNames, err := filepath.Glob(filepath.Join(path, "*.gif"))
	if err != nil {
		logrus.Errorf("Error:", err)
		return err
	}

	for _, gifFileName := range gifFileNames {
		// Open the GIF file
		f, err := os.Open(gifFileName)
		if err != nil {
			logrus.Errorf("Error:", err)
			continue
		}
		defer f.Close()

		// Decode the GIF file
		gifImage, err := gif.DecodeAll(f)
		if err != nil {
			logrus.Errorf("Error:", err)
			continue
		}

		// Get a random frame from the GIF image
		rand.Seed(time.Now().UnixNano())
		frameIndex := rand.Intn(len(gifImage.Image))
		frame := gifImage.Image[frameIndex]

		// Create a new PNG image
		pngImage := image.NewRGBA(image.Rect(0, 0, frame.Bounds().Dx(), frame.Bounds().Dy()))
		// Draw the frame onto the PNG image
		draw.Draw(pngImage, pngImage.Bounds(), frame, image.Point{}, draw.Src)

		err = savePNG(gifFileName, pngImage)
		if err != nil {
			logrus.Errorf("Error:", err)
			continue
		}
	}

	return nil
}

// Save the given image as a PNG file with the same name as the original file but with the ".png" extension
func savePNG(filename string, img *image.RGBA) error {
	pngFilename := fmt.Sprintf("%v.png", strings.TrimSuffix(filename, filepath.Ext(filename)))
	f, err := os.Create(pngFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
