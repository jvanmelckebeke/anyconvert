package tools

import (
	"fmt"
	"jvanmelckebeke/anyconverter-api/pkg/constants"
	"os"
	"path/filepath"
)

func ConvertToWorkPath(inputFname, outputExt string) string {
	baseName, _ := splitExt(filepath.Base(inputFname))
	return filepath.Join(
		constants.UploadsDir,
		fmt.Sprintf("%s.%s", baseName, outputExt))
}

func ConvertToResultPath(inputPath string) string {

	if inputPath != "" && len(inputPath) > len(constants.UploadsDir)+1 {
		return inputPath[len(constants.UploadsDir)+1:]
	}

	return ""
}

func PrepareOutputFile(inputFname string, ext string) string {
	baseName, _ := splitExt(filepath.Base(inputFname))
	if ext[0] != '.' {
		ext = fmt.Sprintf(".%s", ext)
	}
	baseName = fmt.Sprintf("%s%s", baseName, ext)
	outputPath := filepath.Join("/tmp", baseName)

	return outputPath
}

func PrepareFrameDirectory(inputFname string) string {
	baseName, _ := splitExt(filepath.Base(inputFname))
	outputPath := filepath.Join("/tmp", baseName)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.Mkdir(outputPath, os.ModePerm); err != nil {
			fmt.Println(err)
			return ""
		}
	} else {
		files, err := os.ReadDir(outputPath)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		for _, f := range files {
			if err := os.Remove(filepath.Join(outputPath, f.Name())); err != nil {
				fmt.Println(err)
				return ""
			}
		}
	}

	return outputPath
}

func splitExt(filename string) (string, string) {
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	return base, ext
}
