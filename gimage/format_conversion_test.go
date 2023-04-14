package gimage

import "testing"

// In this example, the delay parameter is set to 100, which represents the delay between frames in hundredths of a second. You can adjust this value to control the speed of the animation.
func TestFormatConver(t *testing.T) {
	err := PngsToGifs("/Users/ricco/Downloads/emojo/", 100)
	if err != nil {
		t.Fatal(err)
	}

}
