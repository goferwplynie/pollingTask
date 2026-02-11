package models

type TaskStatus string

const (
	DONE    TaskStatus = "done"
	PENDING TaskStatus = "pending"
	ERROR   TaskStatus = "error"
)

type TaskStatusResponse struct {
	Email  string     `json:"email,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

type TaskResult struct {
	Email  string `json:"email,omitempty"`
	Emails string `json:"emails,omitempty"`
}
