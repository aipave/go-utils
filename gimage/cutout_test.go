package gimage

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"gocv.io/x/gocv"
)

// ERROR: pkg-config: exec: "pkg-config": executable file not found in $PATH
// mac: brew install pkg-config  && brew install opencv@4 && echo 'export PATH="/usr/local/opt/opencv@4/bin:$PATH"' >> ~/.bash_profile && source ~/.bash_profile
// linux: sudo apt-get install pkg-config
// test: pkg-config --modversion opencv4
func TestGocv(t *testing.T) {
	//img := gocv.IMRead("/Users/ricco/Downloads/emojo/02.png", gocv.IMReadColor)
	//defer img.Close()

	//if img.Empty() {
	//	println("Error loading image")
	//	return
	//}
	url := "https://media.discordapp.net/attachments/1095555533781618740/1096121121197265006/matthew37_A_3D-rendered_emoticon_for_acting_very_want_to_plot_e_1d298c8c-e8ba-42aa-83ba-0382067bd9ca.png?width=940&height=940"

	// Get the image data from the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the response body into a buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// Decode the image data
	imgData, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadColor)
	if err != nil {
		fmt.Println(err)
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
