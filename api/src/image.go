package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/tidbyt/go-libwebp/webp"
)

func webpToJpg(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("inputPath is empty")
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Printf("Error reading file: %s\n", err)
		return "", fmt.Errorf("error reading file")
	}

	outputPath := convertPath(inputPath, "jpg")

	output, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating output file: %s\n", err)
		return "", fmt.Errorf("error creating output file")
	}
	defer output.Close()

	options := &webp.DecoderOptions{}

	img, err := webp.DecodeRGBA(data, options)
	if err != nil {
		log.Printf("Error decoding webp file: %s\n", err)
		log.Printf("trying as animation")

		dec, err := webp.NewAnimationDecoder(data)
		if err != nil {
			log.Printf("Error creating animation decoder: %s\n", err)
			return "", fmt.Errorf("error creating animation decoder")
		}
		defer dec.Close()

		anim, err := dec.Decode()
		if err != nil {
			log.Printf("Error decoding animation: %s\n", err)
			return "", fmt.Errorf("error decoding animation")
		}

		img = anim.Image[0]
	}
	log.Printf("Decoded webp file")

	// write the image to jpeg
	if err := jpeg.Encode(output, img, nil); err != nil {
		log.Printf("Error encoding jpeg file: %s\n", err)
		return "", fmt.Errorf("error encoding jpeg file")
	}

	return outputPath, nil
}
