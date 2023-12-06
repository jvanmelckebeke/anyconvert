package main

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const uploadsDir = "/tmp"

var uploads = make(map[string]*Upload)

type Upload struct {
	uploadID   string
	fileName   string
	fileSource string
	createdAt  string
	status     string
	outputPath string
	resultURL  string
	error      string
}

func getFilePath(filesource string) string {
	return filepath.Join(uploadsDir, filesource)
}

func saveFileAndCreateUploadID(c *gin.Context, file *multipart.FileHeader) string {
	uid := uuid.New().String()[0:8]
	upload := &Upload{
		uploadID:   uid,
		fileName:   file.Filename,
		fileSource: fmt.Sprintf("%s_%s", uid, file.Filename),
		createdAt:  time.Now().Format("2006-01-02 15:04:05"),
		status:     "pending",
	}
	uploads[uid] = upload

	fmt.Printf("uploads: %+v\n", uploads)

	fpath := getFilePath(upload.fileSource)
	if err := c.SaveUploadedFile(file, fpath); err != nil {
		panic(err)
	}

	fmt.Printf("File saved to %s\n", fpath)
	return uid
}

func deleteUpload(uploadID string, delay time.Duration) {
	time.Sleep(delay * time.Second)
	fpath := getFilePath(uploads[uploadID].fileSource)
	if err := os.Remove(fpath); err != nil {
		fmt.Println(err)
	}
	delete(uploads, uploadID)
}

func processUpload(c *gin.Context, uploadID string, converterFunc func(string) string) {
	defer func() {
		// Delete the source file after processing
		fpath := getFilePath(uploads[uploadID].fileSource)
		if err := os.Remove(fpath); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Deleted source file %s\n", fpath)
		// After 5 minutes, delete the file
		go deleteUpload(uploadID, 300)
	}()

	uploads[uploadID].status = "processing"
	fmt.Printf("Processing upload %s\n", uploadID)
	fpath := getFilePath(uploads[uploadID].fileSource)
	outputPath := converterFunc(fpath)

	if outputPath == "" {
		uploads[uploadID].status = "error"
		uploads[uploadID].error = "Failed to process file"
		return
	}

	// Remove the uploadsDir from the output path
	if fpath != "" && strings.HasPrefix(fpath, uploadsDir) {
		fmt.Println("Removing uploadsDir from output path")
		fmt.Println(outputPath)
		outputPath = outputPath[len(uploadsDir)+1:]
	}

	uploads[uploadID].outputPath = outputPath
	uploads[uploadID].resultURL = fmt.Sprintf("/uploads/%s", uploadID)
	uploads[uploadID].status = "done"
}
