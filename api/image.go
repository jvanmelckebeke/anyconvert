package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"golang.org/x/image/webp"
)

func webpToJpg(inputPath string) string {
	if inputPath == "" {
		return ""
	}

	outputPath := convertPath(inputPath, "jpg")

	inputFile, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer inputFile.Close()

	webpImage, err := webp.Decode(inputFile)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Is this a webp image?")
		return ""
	} else {
		fmt.Println("webp successfully decoded")
	}

	// Convert webpImage to RGB
	rgbImage := image.NewRGBA(webpImage.Bounds())
	if rgba, ok := webpImage.(*image.RGBA); ok {
		rgbImage = rgba
	} else {
		draw.Draw(rgbImage, webpImage.Bounds(), webpImage, image.Point{}, draw.Over)
	}

	// Resize if needed
	resizedImage := resize.Resize(0, 0, rgbImage, resize.Lanczos3)

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer outputFile.Close()

	// Save the resized image as JPEG
	err = jpeg.Encode(outputFile, resizedImage, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Printf("Saved image to %s\n", outputPath)

	return outputPath
}
