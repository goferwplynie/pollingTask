package api

import (
	"fmt"
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

func (h *Handlers) ShortPollHandler(c *gin.Context) {
	var request models.TaskRequest

	if err := c.ShouldBindJSON(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("bad request body ~w~"))
	}

	taskId := h.taskStore.AddTask(request)

	go func(taskId string) {
		timeout := rand.IntN(10)
		time.Sleep(time.Second * time.Duration(timeout))

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
	c.JSON(http.StatusOK, task)
	if task.Status != models.PENDING {
		c.Header("Location", "/api/short/task_result/"+taskId)

		h.taskStore.RemoveTaskStatus(taskId)
	}
}

func (h *Handlers) ShortTaskFinishedHandler(c *gin.Context) {
	taskId := c.Param("id")

	result, err := h.taskStore.GetFinished(taskId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	h.taskStore.RemoveTaskFinished(taskId)

	c.JSON(http.StatusOK, result)
}
