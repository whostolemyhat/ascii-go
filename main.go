package main

import (
	"ascii-go/ascii"
	"log"
	"os"
)

// web version
// validation
// jpg/png only

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: ascii-go path/to/image")
		os.Exit(1)
	}

	img_path := os.Args[1]

	// get input
	img, err := ascii.ReadImage(img_path)
	if err != nil {
		panic(err)
	}

	output := ascii.Convert(img)

	// write converted
	err = ascii.WriteOutput(output)
	if err != nil {
		panic(err)
	}
}
