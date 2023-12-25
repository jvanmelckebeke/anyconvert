package api

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jvanmelckebeke/anyconverter-api/pkg/domain"
	"jvanmelckebeke/anyconverter-api/pkg/env"
	"jvanmelckebeke/anyconverter-api/pkg/logger"
)

type TaskRepository struct {
	client *redis.Client
}

func NewTaskRepository() *TaskRepository {
	redisHost := env.Getenv("REDIS_HOST", "localhost")

	logger.Debug("creating Redis client", "REDIS_HOST", redisHost)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
	})

	logger.Debug("created Redis client", "client", redisClient)

	//check if redis is available
	_, err := redisClient.Ping(context.Background()).Result()

	if err != nil {
		logger.Error("redis is not available", err)
		panic(err)
	} else {
		logger.Info("redis is available")
	}

	return &TaskRepository{
		client: redisClient,
	}
}

func (repo *TaskRepository) Create(task domain.Task) (string, error) {
	ctx := context.Background()
	// convert task to map
	data := task.ToMap()

	logger.Debug("data", "data", data)

	taskID := task.TaskID

	err := repo.client.HSet(ctx, taskID, data).Err()
	if err != nil {
		logger.Error("create error", err)
		return "", err
	}

	logger.Debug("successfully created task", "taskID", taskID)

	return taskID, nil
}

func (repo *TaskRepository) Get(taskID string) *domain.Task {
	ctx := context.Background()
	data, err := repo.client.HGetAll(ctx, taskID).Result()
	if err != nil {
		logger.Error("Get item error", err.Error())
		return nil
	}

	return domain.TaskFromMap(data)
}

func (repo *TaskRepository) GetAll() ([]domain.Task, error) {
	ctx := context.Background()
	keys, err := repo.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	var tasks []domain.Task
	for _, key := range keys {
		task := repo.Get(key)
		if task != nil {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

func (repo *TaskRepository) Delete(taskID string) error {
	ctx := context.Background()
	err := repo.client.Del(ctx, taskID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *TaskRepository) UpdateTaskStatus(taskID string, status string) (*Task, error) {
	ctx := context.Background()
	err := repo.client.HSet(ctx, taskID, "status", status).Err()
	if err != nil {
		logger.Error("update task status error", err)
		return nil, err
	}

	task := repo.Get(taskID)

	return task, nil
}

func (repo *TaskRepository) SetTaskError(taskID string, taskError error) {
	ctx := context.Background()

	err := repo.client.HSet(ctx, taskID,
		map[string]interface{}{
			"status": "error",
			"error":  taskError.Error(),
		})

	if err != nil {
		logger.Error("set task error", err.Err())
		return
	}
}

func (repo *TaskRepository) UpdateTask(task Task) (*Task, error) {
	ctx := context.Background()

	data := task.ToMap()
	taskID := task.TaskID

	err := repo.client.HSet(ctx, taskID, data).Err()

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TaskRepository) SetOutputPath(id string, path string) {
	ctx := context.Background()

	err := repo.client.HSet(ctx, id, "outputPath", path).Err()
	if err != nil {
		logger.Error("set output path error", err)
		return
	}
}
