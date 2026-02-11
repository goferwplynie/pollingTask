package models

type TaskRequest struct {
	Email string `json:"email,omitempty"`
	Count int    `json:"count,omitempty"`
}
