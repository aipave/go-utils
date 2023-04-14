package gimage

import "testing"

// In this example, the delay parameter is set to 100, which represents the delay between frames in hundredths of a second. You can adjust this value to control the speed of the animation.
func TestFormatConver(t *testing.T) {
	err := PngsToGifs("/Users/me/Downloads/emojo/", 100)
	if err != nil {
		t.Fatal(err)
	}

}

// If you need advanced processing functionalities like selective blur or color manipulation,
// then github.com/disintegration/imaging would be a good choice.