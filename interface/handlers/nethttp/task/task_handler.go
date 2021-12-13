package task

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hajime3778/go-clean-arch/domain"
	_usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
)

const TaskPath string = "/tasks/"

type taskHandler struct {
	taskUsecase _usecase.TaskUsecase
}

// NewTaskHandler タスク機能のHandlerオブジェクトを作成します
func NewTaskHandler(u _usecase.TaskUsecase) *taskHandler {
	return &taskHandler{u}
}

// Handler はタスク機能のHandler関数です
func (t *taskHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := strings.TrimPrefix(r.URL.Path, TaskPath)
	taskID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := t.fetchByID(ctx, w, taskID)
		if err != nil {
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	case http.MethodPut:
		err := t.update(ctx, w, r, taskID)
		if err != nil {
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	case http.MethodDelete:
		err := t.delete(ctx, w, taskID)
		if err != nil {
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// fetchByID IDでタスクを1件取得します
func (t *taskHandler) fetchByID(ctx context.Context, w http.ResponseWriter, id int64) error {
	task, err := t.taskUsecase.FetchByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	output, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	log.Println(string(output))
	return nil
}

// update IDでタスクを1件更新します
func (t *taskHandler) update(ctx context.Context, w http.ResponseWriter, r *http.Request, id int64) error {
	var requestTask UpdateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&requestTask)
	if err != nil {
		return err
	}

	task := domain.Task{
		ID:      id,
		Title:   requestTask.Title,
		Content: requestTask.Content,
		DueDate: requestTask.DueDate,
	}

	err = t.taskUsecase.Update(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

// delete IDでタスクを1件削除します
func (t *taskHandler) delete(ctx context.Context, w http.ResponseWriter, id int64) error {
	err := t.taskUsecase.Delete(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return nil
}
