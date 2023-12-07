package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/status", getUploads)
	r.GET("/uploads/:upload_id", getFile)
	r.POST("/image", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uploadID := saveFileAndCreateUploadID(c, file)
		go processUpload(c, uploadID, mediaToJpg)

		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("/status/%s", uploadID)})
	})

	r.POST("/video", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uploadID := saveFileAndCreateUploadID(c, file)
		go processUpload(c, uploadID, mediaToMp4)

		c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("/status/%s", uploadID)})
	})

	r.GET("/status/:upload_id", getUploadStatus)

	if err := r.Run(":8000"); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func mediaToMp4(filepath string) (string, error) {
	// todo
	return filepath, nil
}
