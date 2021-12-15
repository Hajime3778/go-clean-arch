package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
	_usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
)

const TaskIndexPath string = "/tasks"

type taskIndexHandler struct {
	taskUsecase _usecase.TaskUsecase
}

// NewTaskHandler タスク機能のHandlerオブジェクトを作成します
func NewTaskIndexHandler(u _usecase.TaskUsecase) *taskIndexHandler {
	return &taskIndexHandler{u}
}

// NewTaskIndexHandler
func (t *taskIndexHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		fmt.Println("fetch all tasks")
	case http.MethodPost:
		err := t.create(ctx, w, r)
		if err != nil {
			log.Println(err.Error())
		}
	default:
		w.WriteHeader(404)
	}
}

func (t *taskIndexHandler) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var requestTask CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&requestTask)
	if err != nil {
		return err
	}

	task := domain.Task{
		Title:   requestTask.Title,
		Content: requestTask.Content,
		DueDate: requestTask.DueDate,
	}

	err = t.taskUsecase.Create(ctx, task)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}
