package task

import (
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// CreateTaskRequest: タスク追加時のリクエスト
type FindByUserIDTaskRequest struct {
	Limit  int64 `json:"limit" validate:"required"`
	Offset int64 `json:"offset" validate:"required"`
}

// CreateTaskRequest: タスク追加時のリクエスト
type CreateTaskRequest struct {
	Title   string    `json:"title" validate:"required"`
	Content string    `json:"content" validate:"required"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

// IsCreateRequestValid:
func (r CreateTaskRequest) IsCreateRequestValid() (bool, error) {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateTaskRequest: タスク更新時のリクエスト
type UpdateTaskRequest struct {
	Title   string    `json:"title" validate:"required"`
	Content string    `json:"content" validate:"required"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

func (r UpdateTaskRequest) IsUpdateRequestValid() (bool, error) {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}
