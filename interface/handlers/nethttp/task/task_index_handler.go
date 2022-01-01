package task

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
	httpUtil "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp"
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
		t.findByUserID(ctx, w, r)
	case http.MethodPost:
		t.create(ctx, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (t *taskIndexHandler) findByUserID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func (t *taskIndexHandler) create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var requestTask CreateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestTask)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var ok bool
	if ok, err = requestTask.IsCreateRequestValid(); !ok {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
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
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
