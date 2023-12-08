package api

import (
	"fmt"
	"jvanmelckebeke/anyconverter-api/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

var service = NewTaskService()

// GetTasks godoc
// @Summary Get all tasks
// @Description Get all tasks
// @Produce json
// @Success 200 {array} Task
// @Router /status [get]
func GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllTasksResponse())
}

func GetResult(c *gin.Context) {
	taskID := c.Param("id")
	task, exists := service.GetTask(taskID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	filepath := task.GetFullOutputPath()

	fmt.Println("retruning full path: ", filepath)

	c.File(filepath)
}

func PostImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uploadID := service.NewTask(c, file, "image")

	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("/status/%s", uploadID)})
}

func PostVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadID := service.NewTask(c, file, "video")

	statusEndpoint := constants.CreateStatusEndpoint(uploadID)

	c.JSON(http.StatusOK, gin.H{"status": statusEndpoint})
}

func GetTaskStatus(c *gin.Context) {
	uploadID := c.Param("id")

	response, exists := service.GetTaskResponse(uploadID)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}
