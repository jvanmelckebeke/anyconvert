package domain

import (
	"fmt"
	"github.com/google/uuid"
	"jvanmelckebeke/anyconverter-api/pkg/constants"
	"path/filepath"
	"time"
)

type Task struct {
	TaskID     string
	FileName   string
	FileSource string
	CreatedAt  string
	Status     string
	OutputPath string
	ResultURL  string
	Error      string
	TaskType   string
}

type TaskDTO struct {
	TaskID     string `json:"id"`
	FileName   string `json:"fileName"`
	CreatedAt  string `json:"created_at"`
	Status     string `json:"status"`
	FileSource string `json:"fileSource"`
	ResultURL  string `json:"resultURL"`
	Error      string `json:"error"`
	TaskType   string `json:"taskType"`
}

func NewTask(fileName, taskType string) *Task {
	t := &Task{}
	t.TaskID = uuid.New().String()[0:8]
	t.FileName = fileName
	t.FileSource = fmt.Sprintf("%s_%s", t.TaskID, fileName)
	t.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	t.Status = "pending"
	t.TaskType = taskType

	return t
}

func (t *Task) GetFullSourcePath() string {
	return filepath.Join(constants.UploadsDir, t.FileSource)
}

func (t *Task) GetFullOutputPath() string {
	return filepath.Join(constants.UploadsDir, t.OutputPath)
}

func (t *Task) ToResponse() *TaskDTO {
	return &TaskDTO{
		TaskID:     t.TaskID,
		FileName:   t.FileName,
		CreatedAt:  t.CreatedAt,
		Status:     t.Status,
		FileSource: t.FileSource,
		ResultURL:  t.ResultURL,
		Error:      t.Error,
		TaskType:   t.TaskType,
	}
}

func (t *Task) ToMap() map[string]string {
	return map[string]string{
		"taskID":     t.TaskID,
		"fileName":   t.FileName,
		"createdAt":  t.CreatedAt,
		"status":     t.Status,
		"file":       t.FileSource,
		"resultURL":  t.ResultURL,
		"error":      t.Error,
		"taskType":   t.TaskType,
		"outputPath": t.OutputPath,
	}
}

func TaskFromMap(data map[string]string) *Task {
	t := &Task{}
	t.TaskID = data["taskID"]
	t.FileName = data["fileName"]
	t.CreatedAt = data["createdAt"]
	t.Status = data["status"]
	t.FileSource = data["file"]
	t.ResultURL = data["resultURL"]
	t.Error = data["error"]
	t.TaskType = data["taskType"]
	t.OutputPath = data["outputPath"]

	return t
}
