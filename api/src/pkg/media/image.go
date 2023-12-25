package media

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"jvanmelckebeke/anyconverter-api/pkg/logger"
	"jvanmelckebeke/anyconverter-api/pkg/tools"
	"log"
	"os"
	"path/filepath"

	"github.com/tidbyt/go-libwebp/webp"
)

func imageToJpg(inputPath string, img image.Image) (string, error) {
	outputPath := tools.ConvertToWorkPath(inputPath, "jpg")

	output, err := os.Create(outputPath)
	if err != nil {
		logger.Error("error creating output file", "path", outputPath, err)
		return "", fmt.Errorf("error creating output file")
	}
	defer output.Close()

	if err := jpeg.Encode(output, img, nil); err != nil {
		logger.Error("error encoding jpeg file", err)
		return "", fmt.Errorf("error encoding jpeg file")
	}

	return outputPath, nil
}

func animatedWebpToJpg(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		logger.Error("error reading file", "path", inputPath, err)
		return "", fmt.Errorf("error reading file")
	}

	dec, err := webp.NewAnimationDecoder(data)
	if err != nil {
		logger.Error("error creating animation decoder", err)
		return "", fmt.Errorf("error creating animation decoder")
	}
	defer dec.Close()

	anim, err := dec.Decode()
	if err != nil {
		logger.Error("error creating animation", err)
		return "", fmt.Errorf("error decoding animation")
	}

	return imageToJpg(inputPath, anim.Image[0])
}

func webpToJpg(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		logger.Error("error reading file", "file", inputPath, err)
		return "", fmt.Errorf("error reading file")
	}

	options := &webp.DecoderOptions{}

	img, err := webp.DecodeRGBA(data, options)
	if err != nil {
		logger.Warn("error decoding webp file, trying as animationg")

		return animatedWebpToJpg(inputPath)

	}
	log.Printf("Decoded webp file")

	return imageToJpg(inputPath, img)
}

func gifToJpg(inputPath string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		logger.Error("error reading file", "file", inputPath, err)
		return "", fmt.Errorf("error reading file")
	}
	buffer := bytes.NewBuffer(data)

	img, err := gif.Decode(buffer)
	if err != nil {
		logger.Error("error decoding gif file", err)
		return "", fmt.Errorf("error decoding gif file")
	}

	return imageToJpg(inputPath, img)
}

func ToJpg(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("inputPath is empty")
	}

	extension := filepath.Ext(inputPath)

	switch extension {
	case ".webp":
		return webpToJpg(inputPath)
	case ".jpg":
	case ".jpeg":
		return inputPath, nil
	case ".gif":
		return gifToJpg(inputPath)
	}

	return "", fmt.Errorf("unsupported file type")
}

func imageToPng(outputPath string, img image.Image) error {
	output, err := os.Create(outputPath)
	if err != nil {
		logger.Error("error creating output file", "file", outputPath, err)
		return fmt.Errorf("error creating output file")
	}
	defer output.Close()

	if err := png.Encode(output, img); err != nil {
		logger.Error("error encoding png file", err)
		return fmt.Errorf("error encoding png file")
	}

	return nil
}
