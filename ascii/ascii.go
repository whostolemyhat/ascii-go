package ascii

import (
	"image"
	"image/color"
	"math"
)

// take a pixel value in range 0-255 and list of characters
// convert pixel value to matched string
// accept a map of chars so this can be dynamic
// very basic algo - more adv could look at neighbours for better smoothing/corners
func pixelToChar(px uint8, charMap []string) string {
	// reverse number - default charMap has "." first which should be used for white
	// white is 255 value, black is 0
	offset := int(255 - px)

	// check length of map and chunk 0-255 so each char has a range
	// ie len 9 = 255 / 9 = 28.3
	// so if colour is in range 0-28 pick the first char, 29-58 pick second etc
	multiplier := float64(len(charMap)) / 256.0
	idx := math.Floor(float64(offset) * multiplier)

	return charMap[int(idx)]
}

func Convert(img image.Image) []string {
	// chars for output. Use small chars ('.', ',') for lighter values
	// and larger chars for darker values
	charMap := [9]string{".", ",", ":", ";", "o", "x", "%", "#", "@"}

	bounds := img.Bounds()
	output := make([]string, 0)

	// loop pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// from https://pkg.go.dev/image/png
			px := color.GrayModel.Convert(img.At(x, y)).(color.Gray)

			char := pixelToChar(px.Y, charMap[:])
			// build string
			output = append(output, char)
		}

		// format for text output
		// imgs are usually 1 array of pixels
		output = append(output, "\n")
	}

	return output
}
