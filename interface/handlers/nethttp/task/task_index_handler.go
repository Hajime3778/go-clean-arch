package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
)

const TaskIndexPath string = "/tasks"

type taskIndexHandler struct {
	taskUsecase usecase.TaskUsecase
}

// NewTaskHandler タスク機能のHandlerオブジェクトを作成します
func NewTaskIndexHandler(u usecase.TaskUsecase) *taskIndexHandler {
	return &taskIndexHandler{u}
}

// NewTaskIndexHandler
func (t *taskIndexHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		fmt.Println("fetch all tasks")
	case http.MethodPost:
		t.create(ctx, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (t *taskIndexHandler) create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var requestTask CreateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestTask)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var ok bool
	if ok, err = requestTask.IsCreateRequestValid(); !ok {
		writeJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	task := domain.Task{
		UserID:  1, // TODO: トークンから取得するようにする
		Title:   requestTask.Title,
		Content: requestTask.Content,
		DueDate: requestTask.DueDate,
	}

	err = t.taskUsecase.Create(ctx, task)
	if err != nil {
		writeJSONResponse(w, getStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
