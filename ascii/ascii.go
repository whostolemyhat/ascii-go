package ascii

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

func ReadImage(img_path string) (image.Image, error) {
	// TODO validate

	f, err := os.Open(img_path)

	if err != nil {
		log.Println("File not found")
		return nil, err
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	img, err := png.Decode(reader)
	if err != nil {
		log.Println("Failed to read image")
		return nil, err
	}

	return img, nil
}

// take a pixel value in range 0-255 and list of characters
// convert pixel value to matched string
// accept a map of chars so this can be dynamic
func pixelToChar(px uint8, charMap []string) string {
	offset := int(255 - px)
	multiplier := float64(len(charMap)) / 256.0
	idx := math.Floor(float64(offset) * multiplier)

	return charMap[int(idx)]
}

func WriteOutput(output []string) error {
	filePerms := 0644
	outputString := strings.Join(output, "")
	outputFile := "result.txt"
	return os.WriteFile(outputFile, []byte(outputString), os.FileMode(filePerms))
}

func Convert(img image.Image) []string {
	charMap := [9]string{".", ",", ":", ";", "o", "x", "%", "#", "@"}
	bounds := img.Bounds()
	output := make([]string, 0)

	// loop pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// from https://pkg.go.dev/image/png
			px := color.GrayModel.Convert(img.At(x, y)).(color.Gray)

			char := pixelToChar(px.Y, charMap[:])
			output = append(output, char)
		}

		output = append(output, "\n")
	}

	return output
}
