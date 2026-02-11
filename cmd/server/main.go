package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goferwplynie/pollingTask/internal/api"
)

func main() {
	ts := api.NewTaskStore()
	h := api.NewHandlers(ts)

	r := gin.Default()

	r.POST("/api/short/task", h.ShortPollHandler)
	r.GET("/api/short/task/:id", h.ShortTaskStatusHandler)
	r.GET("/api/short/task_result/:id", h.ShortTaskFinishedHandler)

	r.Run()
}
