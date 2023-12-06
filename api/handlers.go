package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUploads(c *gin.Context) {
	c.JSON(http.StatusOK, uploads)
}

func getFile(c *gin.Context) {
	uploadID := c.Param("upload_id")
	upload, exists := uploads[uploadID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Upload not found"})
		return
	}

	filepath := getFilePath(upload.outputPath)
	c.File(filepath)
}

func getUploadStatus(c *gin.Context) {
	uploadID := c.Param("upload_id")

	fmt.Printf("upload: %+v\n", uploads)

	upload, exists := uploads[uploadID]

	if !exists {
		fmt.Println("Upload not found")
		c.JSON(http.StatusNotFound, gin.H{"detail": "Upload not found"})
		return
	}

	fmt.Println("Upload found")

	// convert upload to json response
	response := gin.H{
		"upload_id":  upload.uploadID,
		"fileName":   upload.fileName,
		"created_at": upload.createdAt,
		"status":     upload.status,
		"fileSource": upload.fileSource,
		"outputPath": upload.outputPath,
		"resultURL":  upload.resultURL,
		"error":      upload.error,
	}

	c.JSON(http.StatusOK, response)
}
