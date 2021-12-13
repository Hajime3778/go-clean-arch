package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := t.fetchByID(ctx, w, taskID)
		if err != nil {
			log.Println(err)
		}
	case http.MethodPost:
		fmt.Println("create tasks")
	case http.MethodPut:
		fmt.Println("update tasks")
	case http.MethodDelete:
		fmt.Println("delete tasks")
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// fetchByID IDでタスクを1件取得します
func (t *taskHandler) fetchByID(ctx context.Context, w http.ResponseWriter, id int64) error {
	task, err := t.taskUsecase.FetchByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	output, err := json.Marshal(task)
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}

	_, err = w.Write(output)
	if err != nil {
		return err
	}

	log.Println(output)
	return nil
}
