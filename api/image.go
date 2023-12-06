package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/tidbyt/go-libwebp/test/util"
	"github.com/tidbyt/go-libwebp/webp"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func webpToJpg(inputPath string) string {
	if inputPath == "" {
		return ""
	}

	outputPath := convertPath(inputPath, "jpg")

	data := util.ReadFile(inputPath)

	options := &webp.DecoderOptions{}
	webpImage, err := webp.DecodeRGBA(data, options)
	if err != nil {
		fmt.Println(err)
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
