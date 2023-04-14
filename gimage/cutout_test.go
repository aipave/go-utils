package gimage

import (
	"bytes"
	"image"
	"net/http"
	"testing"

	"gocv.io/x/gocv"
)

// ERROR: pkg-config: exec: "pkg-config": executable file not found in $PATH
// mac: brew install pkg-config  && brew install opencv@4 && echo 'export PATH="/usr/local/opt/opencv@4/bin:$PATH"' >> ~/.bash_profile && source ~/.bash_profile
// linux: sudo apt-get install pkg-config
// test: $ pkg-config --modversion opencv4
//       4.7.0

func ImageFromUrl() ([]byte, error) {
	url := "https://media.discordapp.net/attachments/1095555533781618740/1096121121197265006/matthew37_A_3D-rendered_emoticon_for_acting_very_want_to_plot_e_1d298c8c-e8ba-42aa-83ba-0382067bd9ca.png?width=940&height=940"

	// Get the image data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Read the response body into a buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}

func ImageFromFile() ([]byte, error) {
	img := gocv.IMRead("/Users/path", gocv.IMReadColor)
	defer img.Close()

	if img.Empty() {
		println("Error loading image")
		return []byte{}, nil
	}
	return img.ToBytes(), nil

}
func TestGocv(t *testing.T) {
	bytes, _ := ImageFromUrl()

	// Decode the image data
	imgData, err := gocv.IMDecode(bytes, gocv.IMReadColor)
	if err != nil {
		t.Log(err)
		return
	}

	window := gocv.NewWindow("My Window")
	defer window.Close()

	for {
		window.IMShow(imgData)
		if window.WaitKey(1) >= 0 {
			break
		}
	}

}

func TestCutout(t *testing.T) {
	input := "/Users/me/Downloads/01-道场小恶魔-3D版/主图/02.png"
	output := "/Users/me/Downloads/01-道场小恶魔-3D版/主图/02-new.png"
	// Load the input PNG image
	img := gocv.IMRead(input, gocv.IMReadColor)
	if img.Empty() {
		t.Error("Error loading image")
		return
	}
	defer img.Close()

	// Create a mask with the same size as the image
	mask := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	defer mask.Close()

	// Initialize the mask with zeros
	mask.SetTo(gocv.Scalar{
		//Val1: 0, Val2: 0, Val3: 0, Val4: 255,
		Val1: 0, Val2: 0, Val3: 0, Val4: 0,
	})

	// Define the rectangle enclosing the object of interest
	rectVal := 1
	rect := image.Rect(rectVal, rectVal, img.Cols()-rectVal, img.Rows()-rectVal)

	t.Logf("%vx%v", img.Cols(), mask.Rows())

	// Run the GrabCut algorithm to remove the background
	// cols == 13*componentsCount:
	// The matrix should have a number of columns equal to 13 times the number of components in the Gaussian mixture model used by the algorithm.
	bgdModel := gocv.NewMatWithSize(1, 65, gocv.MatTypeCV64FC1)
	bgdModel.SetTo(gocv.Scalar{Val1: 0, Val2: 255, Val3: 0, Val4: 0})

	fgdModel := gocv.NewMatWithSize(1, 65, gocv.MatTypeCV64FC1)
	fgdModel.SetTo(gocv.Scalar{Val1: 0, Val2: 255, Val3: 255, Val4: 0})

	gocv.GrabCut(img, &mask, rect, &bgdModel, &fgdModel, 5, gocv.GCInitWithRect)
	t.Logf("%vx%v", img.Cols(), mask.Rows())

	// Create a new image with transparent background
	result := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8UC4)
	defer result.Close()

	// Set the pixels of the result image to transparent or opaque based on the mask
	for i := 0; i < img.Rows(); i++ {
		for j := 0; j < img.Cols(); j++ {

			if mask.GetUCharAt(i, j) == 2 || mask.GetUCharAt(i, j) == 1 {
				result.SetUCharAt(i, j*4, 0)
				result.SetUCharAt(i, j*4+1, 0)
				result.SetUCharAt(i, j*4+2, 0)
				result.SetUCharAt(i, j*4+3, 0)
			} else {
				result.SetUCharAt(i, j*4, img.GetUCharAt(i, j*3))
				result.SetUCharAt(i, j*4+1, img.GetUCharAt(i, j*3+1))
				result.SetUCharAt(i, j*4+2, img.GetUCharAt(i, j*3+2))
				result.SetUCharAt(i, j*4+3, 255)
			}
		}
	}

	// Save the result as a new PNG image
	gocv.IMWrite(output, result)

}
