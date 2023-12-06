package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func fileSize(filename string) int64 {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return fileInfo.Size()
}

func bytesToHuman(raw int64) string {
	return humanSize(raw)
}

func convertPath(inputFname, outputExt string) string {
	baseName, _ := splitExt(filepath.Base(inputFname))
	return filepath.Join("/tmp", fmt.Sprintf("%s.%s", baseName, outputExt))
}

func convertDirectories(inputFname string) string {
	baseName, _ := splitExt(filepath.Base(inputFname))
	outputDir := filepath.Join("/tmp", baseName)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
			fmt.Println(err)
			return ""
		}
	} else {
		files, err := os.ReadDir(outputDir)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		for _, f := range files {
			if err := os.Remove(filepath.Join(outputDir, f.Name())); err != nil {
				fmt.Println(err)
				return ""
			}
		}
	}

	return outputDir
}

func splitExt(filename string) (string, string) {
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	return base, ext
}

func humanSize(size int64) string {
	const (
		_        = iota
		kB int64 = 1 << (10 * iota)
		mB
		gB
		tB
		pB
		eB
	)

	sizeUnits := []struct {
		suffix string
		size   int64
	}{
		{"EB", eB},
		{"PB", pB},
		{"TB", tB},
		{"GB", gB},
		{"MB", mB},
		{"kB", kB},
	}

	for _, unit := range sizeUnits {
		if size >= unit.size {
			return fmt.Sprintf("%.2f %s", float64(size)/float64(unit.size), unit.suffix)
		}
	}

	return fmt.Sprintf("%d B", size)
}
