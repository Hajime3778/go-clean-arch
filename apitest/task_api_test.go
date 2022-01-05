package apitest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	repository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	"github.com/stretchr/testify/assert"
)

const taskURL = "http://localhost:8080/tasks"

func TestGetByID(t *testing.T) {
	env.NewEnv().LoadEnvFile("../.env")
	sqlDriver := database.NewSqlConnenction()
	repo := repository.NewTaskRepository(sqlDriver)

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		createTask := domain.Task{
			UserID:  1,
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		createdID, err := repo.Create(context.TODO(), createTask)
		if err != nil {
			t.Error(err)
		}

		req, _ := http.NewRequest("GET", taskURL+"/"+strconv.Itoa(int(createdID)), nil)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		var resTask domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resTask)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, createdID, resTask.ID)
		assert.Equal(t, createTask.Title, resTask.Title)
		assert.Equal(t, createTask.Content, resTask.Content)
	})
}
