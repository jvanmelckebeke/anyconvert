package api

import (
	"fmt"
	"jvanmelckebeke/anyconverter-api/pkg/constants"
	"jvanmelckebeke/anyconverter-api/pkg/domain"
	"jvanmelckebeke/anyconverter-api/pkg/env"
	"jvanmelckebeke/anyconverter-api/pkg/logger"
	"jvanmelckebeke/anyconverter-api/pkg/media"
	"jvanmelckebeke/anyconverter-api/pkg/tools"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Task = domain.Task

type TaskService struct {
	repo *TaskRepository
}

func NewTaskService() *TaskService {
	return &TaskService{
		repo: NewTaskRepository(),
	}
}

func (t *TaskService) GetTask(taskID string) (*Task, bool) {

	task := t.repo.Get(taskID)
	if task == nil {
		logger.Warn("task is nil")
		return nil, false
	}

	return task, true

}

func (t *TaskService) GetTaskResponse(taskID string) (*domain.TaskDTO, bool) {
	task, exists := t.GetTask(taskID)
	if !exists {
		return nil, false
	}
	return task.ToResponse(), true
}

func (t *TaskService) GetAllTasksResponse() []domain.TaskDTO {

	tasks := make([]domain.TaskDTO, 0)
	rawTasks, err := t.repo.GetAll()

	if err != nil {
		logger.Error("get all tasks error", err)
		return nil
	}

	for _, task := range rawTasks {
		tasks = append(tasks, *task.ToResponse())
	}
	return tasks
}

func (t *TaskService) NewTask(c *gin.Context, file *multipart.FileHeader, taskType string) string {
	task := domain.NewTask(file.Filename, taskType)
	uid := task.TaskID

	_, err := t.repo.Create(*task)
	if err != nil {
		logger.Error("create task error", err)
		panic(err)
	}

	fpath := task.GetFullSourcePath()
	if err := c.SaveUploadedFile(file, fpath); err != nil {
		logger.Error("save uploaded file error", err)
		panic(err)
	}

	logger.Info("uploaded file saved", "filepath", fpath)
	go t.ProcessUpload(uid)
	return uid
}

func (t *TaskService) DeleteTask(taskID string) {
	if err := t.repo.Delete(taskID); err != nil {
		logger.Error("delete task error", err)
	}
}

func (t *TaskService) DelayedDeleteSource(uploadID string, delay time.Duration) {
	time.Sleep(delay * time.Second)

	task := t.repo.Get(uploadID)
	if task == nil {
		logger.Warn("task is nil")
		return
	}

	fpath := task.GetFullSourcePath()
	if err := os.Remove(fpath); err != nil {
		logger.Error("remove path error", "path", fpath, err)
		return
	}

	logger.Info("deleted source file", "filepath", fpath)

}

func (t *TaskService) DelayedDeleteOutput(uploadID string, delay time.Duration) {
	time.Sleep(delay * time.Second)

	task := t.repo.Get(uploadID)
	if task == nil {
		logger.Warn("task is nil")
		return
	}

	fpath := task.GetFullOutputPath()
	if err := os.Remove(fpath); err != nil {
		logger.Error("remove path error", "path", fpath, err)
		return
	}
	t.DeleteTask(uploadID)
	logger.Info("deleted output file", "filepath", fpath)
	logger.Info("deleted task", "taskID", uploadID)
}

func (t *TaskService) DeleteFailedTasks() {
	tasks, err := t.repo.GetAll()
	if err != nil {
		logger.Error("get all tasks error", err)
		return
	}

	for _, task := range tasks {
		if task.Status == "error" {
			t.DeleteTask(task.TaskID)
		}
	}
}

func (t *TaskService) ProcessUpload(uploadID string) {
	defer func() {
		task := t.repo.Get(uploadID)
		if task == nil {
			logger.Warn("task is nil")
			return
		}

		sourceDelay := env.GetenvInt("SOURCE_DELETE_DELAY", 300)
		resultDelay := env.GetenvInt("RESULT_DELETE_DELAY", 300)

		logger.Info("setting delayed delete", "sourceDelay", sourceDelay, "resultDelay", resultDelay, "taskID", uploadID)

		go t.DelayedDeleteSource(uploadID, time.Duration(sourceDelay))
		go t.DelayedDeleteOutput(uploadID, time.Duration(resultDelay))
	}()

	task, exists := t.GetTask(uploadID)

	if !exists {
		logger.Error("task not found", "taskID", uploadID)
		return
	}

	t.process(task)
}

func (t *TaskService) processAsImage(task *Task) (string, error) {
	path, err := media.ToJpg(task.GetFullSourcePath())
	if err != nil {
		return "", err
	}

	logger.Info("task successfully converted to jpg", "taskID", task.TaskID, "path", path)

	return path, nil
}

func (t *TaskService) processAsVideo(task *Task) (string, error) {
	path, err := media.ToMp4(task.GetFullSourcePath())
	if err != nil {
		log.Println(err)
		return "", err
	}

	logger.Info("task successfully converted to mp4", "taskID", task.TaskID, "path", path)

	t.repo.SetOutputPath(task.TaskID, path)
	task.OutputPath = path

	return path, nil

}

func (t *TaskService) processAsAudio(task *Task) (string, error) {
	path, err := media.Mp4ToMp3(task.GetFullSourcePath())
	if err != nil {
		log.Println(err)
		return "", err
	}

	logger.Info("task successfully converted to mp3", "taskID", task.TaskID, "path", path)

	t.repo.SetOutputPath(task.TaskID, path)
	task.OutputPath = path

	return path, nil
}

func (t *TaskService) process(task *Task) {
	task, err := t.repo.UpdateTaskStatus(task.TaskID, "processing")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info("processing task", "taskID", task.TaskID)

	var outPath string

	if task.TaskType == "image" {
		outPath, err = t.processAsImage(task)
	} else if task.TaskType == "video" {
		outPath, err = t.processAsVideo(task)
	} else if task.TaskType == "audio" {
		outPath, err = t.processAsAudio(task)
	} else {
		err = fmt.Errorf("unknown task type")
	}

	if err != nil {
		t.repo.SetTaskError(task.TaskID, err)
		return
	}

	outPath = tools.ConvertToResultPath(outPath)

	task.OutputPath = outPath
	task.Status = "done"
	task.ResultURL = constants.CreateResultEndpoint(task.TaskID)

	_, err = t.repo.UpdateTask(*task)

	if err != nil {
		fmt.Println(err)

	}
}
