package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadPath = "/tmp/"

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/uploads/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File(filepath.Join(uploadPath, filename))
	})

	router.POST("/image", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := saveUploadedFile(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newFile := webpToJpg(filepath.Join(uploadPath, file.Filename))
		c.JSON(http.StatusOK, gin.H{"file": "/uploads" + newFile})
	})

	router.POST("/video", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := saveUploadedFile(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newFile := mediaToMp4(filepath.Join(uploadPath, file.Filename))
		c.JSON(http.StatusOK, gin.H{"file": "/uploads" + newFile})
	})

	if err := router.Run(":8000"); err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Next()
	}
}

func saveUploadedFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(uploadPath, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func webpToJpg(filePath string) string {
	// Implement your webp to jpg conversion logic in Go
	// and return the new file path
	return filePath
}

func mediaToMp4(filePath string) string {
	// Implement your media to mp4 conversion logic in Go
	// and return the new file path
	return filePath
}
