package api

import (
	"fmt"
	"jvanmelckebeke/anyconverter-api/pkg/domain"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Task = domain.Task

type TaskService struct {
	tasks map[string]*Task
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*Task),
	}
}

func (t *TaskService) GetTask(taskID string) (*Task, bool) {
	task, exists := t.tasks[taskID]
	return task, exists
}

func (t *TaskService) GetTaskResponse(taskID string) (*domain.TaskDTO, bool) {
	task, exists := t.tasks[taskID]
	if !exists {
		return nil, false
	}
	return task.ToResponse(), true
}

func (s *TaskService) GetAllTasksResponse() []domain.TaskDTO {
	var tasks []domain.TaskDTO
	for _, task := range s.tasks {
		tasks = append(tasks, *task.ToResponse())
	}
	return tasks
}

func (s *TaskService) NewTask(c *gin.Context, file *multipart.FileHeader, taskType string) string {
	task := domain.NewTask(file.Filename, taskType)
	uid := task.TaskID

	s.tasks[uid] = task

	fpath := task.GetFullSourcePath()
	if err := c.SaveUploadedFile(file, fpath); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("File saved to %s\n", fpath)
	go s.ProcessUpload(uid)
	return uid
}

func (s *TaskService) DeleteTask(taskID string) {
	delete(s.tasks, taskID)
}

func (s *TaskService) DelayedDeleteUpload(uploadID string, delay time.Duration) {
	time.Sleep(delay * time.Second)
	fpath := s.tasks[uploadID].GetFullOutputPath()
	if err := os.Remove(fpath); err != nil {
		fmt.Println(err)
	}
	delete(s.tasks, uploadID)
}

func (s *TaskService) DeleteFailedTasks() {
	for _, task := range s.tasks {
		if task.Status == "error" {
			s.DeleteTask(task.TaskID)
		}
	}
}

func (t *TaskService) ProcessUpload(uploadID string) {
	defer func() {
		fpath := t.tasks[uploadID].GetFullSourcePath()
		if err := os.Remove(fpath); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Deleted source file %t\n", fpath)
		// After 5 minutes, delete the file
		go t.DelayedDeleteUpload(uploadID, 300)
	}()

	task, exists := t.GetTask(uploadID)

	if !exists {
		fmt.Println("Task not found")
		return
	}

	task.Process()
}
