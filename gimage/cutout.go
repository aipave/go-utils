package gimage

// a threshold value threshold (between 0 and 255)
//func CutoutImage(img image.Image, threshold uint8) image.Image {
//	// Select the background using the magic wand tool.
//	selectionMask := imaging.SelectiveBlur(imaging.Grayscale(img), 2.0, threshold, imaging.Gaussian)
//
//	// Cut out the background and fill with transparency.
//	mask := imaging.AdjustFunc(selectionMask, func(c color.NRGBA) color.NRGBA {
//		if c.A > 0 {
//			return color.NRGBA{0, 0, 0, 255}
//		}
//		return color.NRGBA{0, 0, 0, 0}
//	})
//	return imaging.Paste(img, mask, image.Point{})
//}
