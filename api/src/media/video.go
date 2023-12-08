package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ffmpegProcess(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func webmToMp4(inputPath string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func gifToMp4(inputPath string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func animatedWebpToMp4(inputPath string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
func mediaToMp4(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("empty inputPath provided")
	}

	extenstion := filepath.Ext(inputPath)

	switch extenstion {
	case ".webm":
		return webmToMp4(inputPath)
	case ".gif":
		return gifToMp4(inputPath)
	case ".webp":
		return animatedWebpToMp4(inputPath)
	}

	return "", fmt.Errorf("unsupported file type")
}
