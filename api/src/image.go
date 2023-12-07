package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
)

func webpToJpg(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("inputPath is empty")
	}

	outputPath := convertPath(inputPath, "jpg")

	file, err := os.Open(inputPath)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return "", fmt.Errorf("error opening file")
	}

	output, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating file: %s\n", err)
		return "", fmt.Errorf("error creating file")
	}
	defer output.Close()

	img, err := webp.Decode(file, &decoder.Options{})
	if err != nil {
		log.Printf("Error decoding webp file: %s\n", err)
		return "", fmt.Errorf("error decoding webp file")
	}

	if err := jpeg.Encode(output, img, nil); err != nil {
		log.Printf("Error encoding jpeg file: %s\n", err)
		return "", fmt.Errorf("error encoding jpeg file")
	}

	// Convert webpImage to RGB
	fmt.Printf("Saved image to %s\n", outputPath)

	return outputPath, nil
}
