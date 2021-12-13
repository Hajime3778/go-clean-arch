package task

import "time"

// Task ...
type CreateTaskRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	DueDate time.Time `json:"due_date"`
}
