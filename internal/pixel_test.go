package ascii

import (
	"fmt"
	"testing"
)

func TestPixelToChar(t *testing.T) {
	// pixelToChar should chunk into same num as charmap len
	// so len = 9, 256/len = ~28
	// so every 29 units the char should change
	// black = 0, white = 255
	// I've defined chars from light to dark so should be flipped (ie first in map is light)
	charMap := [9]string{".", ",", ":", ";", "o", "x", "%", "#", "@"}

	var tests = []struct {
		colourVal int
		expected  string
	}{
		{240, "."},
		{210, ","},
		{180, ":"},
		{160, ";"},
		{125, "o"},
		{98, "x"},
		{68, "%"},
		{32, "#"},
		{12, "@"},
	}

	for _, testData := range tests {
		testName := fmt.Sprintf("%d, %s", testData.colourVal, testData.expected)
		t.Run(testName, func(t *testing.T) {
			char := pixelToChar(uint8(testData.colourVal), charMap[:])
			if char != testData.expected {
				t.Errorf("pixelToChar(%d) = '%s', want '.'", testData.colourVal, char)
			}
		})
	}
}
