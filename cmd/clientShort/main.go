package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/goferwplynie/pollingTask/internal/models"
)

const (
	Host      = "http://127.0.0.1:8080"
	ShortBase = "/api/short"
	LongBase  = "/api/long"
)

type Logger struct{}

func (*Logger) Log(info any) {
	t := time.Now().Format(time.DateTime)
	color.RGB(235, 52, 235).Printf("[INFO]%s %v\n", t, info)
}
func (*Logger) Error(info any) {
	t := time.Now().Format(time.DateTime)
	color.RGB(255, 0, 0).Printf("[ERROR]%s %v\n", t, info)
}

var logger Logger

func main() {
	long()
	short()
}

func long() {
	logger = Logger{}

	taskRequest := models.TaskRequest{
		Email: "nEY9R@example.com",
		Count: 5,
	}
	marshalled, _ := json.Marshal(taskRequest)

	req, err := http.NewRequest(http.MethodPost, Host+LongBase+"/task", bytes.NewBuffer(marshalled))
	if err != nil {
		logger.Error("failed to build request T~T")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("request failed T~T")
		return
	}

	resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Error("Server returned error status: " + resp.Status)
		return
	}

	statusLocation := resp.Header.Get("Location")
	if statusLocation == "" {
		logger.Error("Empty Location header! Backend didn't tell us where to look.")
		return
	}

	logger.Log("Task submitted! Checking status at: " + statusLocation)
}

func short() {
	logger = Logger{}
	taskRequest := models.TaskRequest{
		Email: "nEY9R@example.com",
		Count: 5,
	}
	marshalled, _ := json.Marshal(taskRequest)

	req, err := http.NewRequest(http.MethodPost, Host+ShortBase+"/task", bytes.NewBuffer(marshalled))
	if err != nil {
		logger.Error("failed to build request T~T")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("request failed T~T")
		return
	}

	resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Error("Server returned error status: " + resp.Status)
		return
	}

	statusLocation := resp.Header.Get("Location")
	if statusLocation == "" {
		logger.Error("Empty Location header! Backend didn't tell us where to look.")
		return
	}

	logger.Log("Task submitted! Checking status at: " + statusLocation)

	resultLocation, err := shortPoll(statusLocation)
	if err != nil {
		logger.Error("short polling failed T~T")
		logger.Error(err)
		return
	}

	finalURL := Host + resultLocation

	req, err = http.NewRequest(http.MethodGet, finalURL, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	var result models.TaskResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to unmarshall result")
		return
	}

	logger.Log("Final Result received! >^.^<")
	logger.Log(result)
}

func shortPoll(location string) (string, error) {
	targetURL := Host + location

	for {
		req, err := http.NewRequest(http.MethodGet, targetURL, nil)
		if err != nil {
			return "", err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Error("polling request failed")
			return "", err
		}

		var response models.TaskStatusResponse
		err = json.NewDecoder(resp.Body).Decode(&response)

		resp.Body.Close()

		if err != nil {
			logger.Error("failed to decode status body")
			return "", err
		}

		if response.Status != models.PENDING {
			logger.Log("Meeeow! task finished, status: " + response.Status + " >^w^<")

			return resp.Header.Get("Location"), nil
		}

		logger.Log("status: " + response.Status + " mrrrp TwT")
		time.Sleep(1 * time.Second)
	}
}
