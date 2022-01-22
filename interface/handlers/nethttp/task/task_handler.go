package task

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/domain/constant"
	httpUtil "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
)

const TaskPath string = "/tasks/"

type taskHandler struct {
	taskUsecase usecase.TaskUsecase
}

// NewTaskHandler タスク機能のHandlerオブジェクトを作成します
func NewTaskHandler(u usecase.TaskUsecase) *taskHandler {
	return &taskHandler{u}
}

// Handler はタスク機能のHandler関数です
func (t *taskHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := strings.TrimPrefix(r.URL.Path, TaskPath)
	taskID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		log.Println(err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		t.getByID(ctx, w, r, taskID)
	case http.MethodPut:
		t.update(ctx, w, r, taskID)
	case http.MethodDelete:
		t.delete(ctx, w, r, taskID)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetByID IDでタスクを1件取得します
func (t *taskHandler) getByID(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) {
	token, userID, err := httpUtil.VerifyAccessToken(r)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx = context.WithValue(ctx, constant.UserIDContextKey, userID)
	ctx = context.WithValue(ctx, constant.AuthTokenContextKey, token)
	task, err := t.taskUsecase.GetByID(ctx, id)
	if err != nil {
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}
	httpUtil.WriteJSONResponse(w, http.StatusOK, task)
}

// update IDでタスクを1件更新します
func (t *taskHandler) update(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) {
	token, userID, err := httpUtil.VerifyAccessToken(r)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx = context.WithValue(ctx, constant.UserIDContextKey, userID)
	ctx = context.WithValue(ctx, constant.AuthTokenContextKey, token)

	var requestTask UpdateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&requestTask)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var ok bool
	if ok, err = requestTask.IsUpdateRequestValid(); !ok {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	task := domain.Task{
		ID:      id,
		Title:   requestTask.Title,
		Content: requestTask.Content,
		DueDate: requestTask.DueDate,
	}

	err = t.taskUsecase.Update(ctx, task)
	if err != nil {
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// delete IDでタスクを1件削除します
func (t *taskHandler) delete(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) {
	token, userID, err := httpUtil.VerifyAccessToken(r)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx = context.WithValue(ctx, constant.UserIDContextKey, userID)
	ctx = context.WithValue(ctx, constant.AuthTokenContextKey, token)

	err = t.taskUsecase.Delete(ctx, id)
	if err != nil {
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
