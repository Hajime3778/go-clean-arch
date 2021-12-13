package task

import "time"

// CreateTaskRequest: タスク追加時のリクエスト
type CreateTaskRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	DueDate time.Time `json:"due_date"`
}

// UpdateTaskRequest: タスク更新時のリクエスト
type UpdateTaskRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	DueDate time.Time `json:"due_date"`
}
