package api

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goferwplynie/pollingTask/internal/models"
)

type Handlers struct {
	taskStore *TaskStore
}

func NewHandlers(ts *TaskStore) *Handlers {
	return &Handlers{
		taskStore: ts,
	}
}

func (h *Handlers) LongPollHandler(c *gin.Context) {
	var request models.TaskRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		b, _ := io.ReadAll(c.Request.Body)
		log.Println(string(b))
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bad request body ~w~"))
	}

	channel := make(chan models.TaskStatusResponse, 1)
	taskId := h.taskStore.AddLongTask(channel)

	go func(request models.TaskRequest, channel chan models.TaskStatusResponse, taskId string) {
		timeout := time.Minute * 3
		time.Sleep(timeout)

		channel <- models.TaskStatusResponse{
			Email:  request.Email,
			Status: models.DONE,
		}
		h.taskStore.FinishTask(taskId, request.Email)
	}(request, channel, taskId)
	c.Header("Location", "/api/long/task/"+taskId)
}

func (h *Handlers) LongTaskStatusHandler(c *gin.Context) {
	taskId := c.Param("id")

	taskChan, err := h.taskStore.GetLongTask(taskId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("task not found TwT"))
		return
	}

	select {
	case taskStatus := <-taskChan:
		c.Header("Location", "/api/long/task_result/"+taskId)
		c.JSON(http.StatusOK, taskStatus)
	case <-c.Request.Context().Done():
		log.Println("client disconnected, releasing channel listener")
		return
	}
}

func (h *Handlers) ShortPollHandler(c *gin.Context) {
	var request models.TaskRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		b, _ := io.ReadAll(c.Request.Body)
		log.Println(string(b))
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bad request body ~w~"))
	}

	taskId := h.taskStore.AddTask(request)

	go func(taskId string) {
		timeout := time.Second * time.Duration(rand.IntN(10))
		time.Sleep(timeout)

		h.taskStore.ChangeStatus(taskId, models.DONE)
	}(taskId)
	c.Header("Location", "/api/short/task/"+taskId)
}

func (h *Handlers) ShortTaskStatusHandler(c *gin.Context) {
	taskId := c.Param("id")

	task, err := h.taskStore.GetStatus(taskId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if task.Status != models.PENDING {
		c.Header("Location", "/api/short/task_result/"+taskId)

		h.taskStore.RemoveTaskStatus(taskId)
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handlers) TaskFinishedHandler(c *gin.Context) {
	taskId := c.Param("id")

	result, err := h.taskStore.GetFinished(taskId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	h.taskStore.RemoveTaskFinished(taskId)

	c.JSON(http.StatusOK, result)
}
