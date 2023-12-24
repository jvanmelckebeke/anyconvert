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
// @Success 200 {array} TaskDTO
// @Router /status [get]
func GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllTasksResponse())
}

// GetResult godoc
// @Summary Get a specific task result
// @Description Get the result of a specific task by its ID
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {file} file "The result of the task"
// @Failure 404 {object} map[string]any "Task not found"
// @Router /result/{id} [get]
// GetResult is a function that handles the GET request at the /result/{id} endpoint.
// It uses the service to get a task by its ID and sends the result as a file.
// If the task does not exist, it sends a JSON response with a status code of 404 and a detail message.
func GetResult(c *gin.Context) {
	taskID := c.Param("id")
	task, exists := service.GetTask(taskID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	filepath := task.GetFullOutputPath()

	fmt.Println("returning full path: ", filepath)

	c.File(filepath)
}

// PostImage godoc
// @Summary Convert an image
// @Description Convert an image to a jpg image
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to convert"
// @Success 200 {string} string "The status endpoint of the task"
// @Router /image [post]
func PostImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uploadID := service.NewTask(c, file, "image")

	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("/status/%s", uploadID)})
}

// PostVideo godoc
// @Summary Convert a video
// @Description Convert a video to a mp4 video
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to convert"
// @Success 200 {object} TaskDTO "The status endpoint of the task"
// @Router /video [post]
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

// GetTaskStatus godoc
// @Summary Get a specific task status
// @Description Get the status of a specific task by its ID
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} TaskDTO "The status of the task"
// @Failure 404 {object} map[string]string "Task not found"
func GetTaskStatus(c *gin.Context) {
	uploadID := c.Param("id")

	response, exists := service.GetTaskResponse(uploadID)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}
