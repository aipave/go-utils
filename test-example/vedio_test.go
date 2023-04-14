package test_example

//  go list golang.org/x/image/...

//func TestVideo(t *testing.T) {
//	// Open the video file
//	f, err := os.Open("video.mp4")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer f.Close()
//
//	// Create a video player
//	player, err := ebiten.NewVideoPlayerFromFile(f)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer player.Close()
//
//	// Get the video configuration
//	config := player.VideoInfo()
//
//	// Print some information about the video
//	fmt.Printf("Format: %s\n", config.Format())
//	fmt.Printf("Width: %d\n", config.Width())
//	fmt.Printf("Height: %d\n", config.Height())
//	fmt.Printf("Frame rate: %f\n", config.FPS())
//
//	// Decode and resize the first frame
//	frame, err := player.Frame()
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	dst := image.NewRGBA(image.Rect(0, 0, config.Width()/2, config.Height()/2))
//	graphics.Scale(dst, frame, dst.Bounds(), graphics.Linear)
//
//	// Save the resized frame to a file
//	out, err := os.Create("frame.png")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer out.Close()
//	png.Encode(out, dst)
//}
