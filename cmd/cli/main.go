package main

import (
	ascii "ascii-go/internal"
	"bufio"
	"image"
	"image/png"
	"log"
	"os"
)

// lib is only concerned with conversion so file handling is outside
// only pass an image.Image to lib = flexible so we
// can have a cli and web version

func readImage(img_path string) (image.Image, error) {
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

func writeOutput(output string) error {
	filePerms := 0644
	outputFile := "result.txt"
	return os.WriteFile(outputFile, []byte(output), os.FileMode(filePerms))
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: ascii-go path/to/image")
		os.Exit(1)
	}

	img_path := os.Args[1]

	// get input
	img, err := readImage(img_path)
	if err != nil {
		panic(err)
	}

	output := ascii.Convert(img)

	// write converted
	err = writeOutput(output)
	if err != nil {
		panic(err)
	}
}
