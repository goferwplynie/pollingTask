package api

import (
	"fmt"
	"sync"

	"github.com/goferwplynie/pollingTask/internal/models"
	"github.com/google/uuid"
)

type TaskStore struct {
	tasks         map[string]models.TaskStatusResponse
	longTasks     map[string]chan models.TaskStatusResponse
	finishedTasks map[string]models.TaskResult
	mutex         sync.RWMutex
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:         make(map[string]models.TaskStatusResponse),
		finishedTasks: make(map[string]models.TaskResult),
		longTasks:     make(map[string]chan models.TaskStatusResponse),
		mutex:         sync.RWMutex{},
	}
}

func (ts *TaskStore) AddLongTask(channel chan models.TaskStatusResponse) string {
	taskId := uuid.New().String()

	ts.mutex.Lock()
	ts.longTasks[taskId] = channel
	ts.mutex.Unlock()

	return taskId
}

func (ts *TaskStore) GetLongTask(taskId string) (chan models.TaskStatusResponse, error) {
	ts.mutex.RLock()
	if taskChan, ok := ts.longTasks[taskId]; ok {
		ts.mutex.RUnlock()

		return taskChan, nil
	}

	return nil, fmt.Errorf("task not found")
}

func (ts *TaskStore) AddTask(request models.TaskRequest) string {
	taskId := uuid.New().String()

	ts.mutex.Lock()
	ts.tasks[taskId] = models.TaskStatusResponse{
		Email:  request.Email,
		Status: models.PENDING,
	}
	ts.mutex.Unlock()

	return taskId
}

func (ts *TaskStore) ChangeStatus(taskId string, status models.TaskStatus) {
	task := ts.tasks[taskId]

	task.Status = models.DONE

	ts.mutex.Lock()
	ts.tasks[taskId] = task
	ts.mutex.Unlock()
	ts.FinishTask(taskId, task.Email)
}

func (ts *TaskStore) FinishTask(taskId string, email string) {
	ts.mutex.Lock()
	ts.finishedTasks[taskId] = models.TaskResult{
		Email:  email,
		Emails: "meow",
	}
	ts.mutex.Unlock()
}

func (ts *TaskStore) GetStatus(taskId string) (models.TaskStatusResponse, error) {
	ts.mutex.RLock()
	if task, ok := ts.tasks[taskId]; ok {
		ts.mutex.RUnlock()

		return task, nil
	}

	return models.TaskStatusResponse{}, fmt.Errorf("task not found")
}

func (ts *TaskStore) RemoveTaskStatus(taskId string) error {
	ts.mutex.Lock()
	delete(ts.tasks, taskId)
	ts.mutex.Unlock()
	return fmt.Errorf("task not found")
}

func (ts *TaskStore) GetFinished(taskId string) (models.TaskResult, error) {
	ts.mutex.RLock()
	if task, ok := ts.finishedTasks[taskId]; ok {
		ts.mutex.RUnlock()

		return task, nil
	}

	return models.TaskResult{}, fmt.Errorf("task not found")
}

func (ts *TaskStore) RemoveTaskFinished(taskId string) error {
	ts.mutex.Lock()
	delete(ts.finishedTasks, taskId)
	ts.mutex.Unlock()
	return fmt.Errorf("task not found")
}
