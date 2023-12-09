package domain

import (
	"fmt"
	"github.com/google/uuid"
	"jvanmelckebeke/anyconverter-api/constants"
	"jvanmelckebeke/anyconverter-api/media"
	"jvanmelckebeke/anyconverter-api/tools"
	"log"
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
	OutputPath string `json:"outputPath"`
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
		OutputPath: t.OutputPath,
		ResultURL:  t.ResultURL,
		Error:      t.Error,
		TaskType:   t.TaskType,
	}
}

func (t *Task) ProcessAsImage() (string, error) {
	path, err := media.ToJpg(t.GetFullSourcePath())
	if err != nil {
		return "", err
	}

	fmt.Printf("task %s sucessfully converted to jpg at %s", t.TaskID, path)
	return path, nil
}

func (t *Task) ProcessAsVideo() (string, error) {
	path, err := media.ToMp4(t.GetFullSourcePath())
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("task %s sucessfully converted to mp4 at %s", t.TaskID, path)

	return path, nil

}

func (t *Task) Process() {
	t.Status = "processing"
	fmt.Println("Processing task", t.TaskID)

	var err error
	var outPath string

	if t.TaskType == "image" {
		outPath, err = t.ProcessAsImage()
	} else if t.TaskType == "video" {
		outPath, err = t.ProcessAsVideo()
	} else {
		err = fmt.Errorf("unknown task type")
	}

	if err != nil {
		t.Status = "error"
		t.Error = err.Error()
		return
	}

	outPath = tools.ConvertToResultPath(outPath)

	t.OutputPath = outPath
	t.Status = "done"
	t.ResultURL = constants.CreateResultEndpoint(t.TaskID)

}
