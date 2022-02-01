package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	interfaceDB "github.com/Hajime3778/go-clean-arch/interface/database"
	taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	userRepository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	taskHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/task"
	"github.com/Hajime3778/go-clean-arch/util/token"
	"github.com/stretchr/testify/assert"
)

const taskURL = "http://localhost:8080/tasks"

var sqlDriver interfaceDB.SqlDriver

func TestMain(m *testing.M) {
	env.NewEnv().LoadEnvFile("../../.env")
	sqlDriver = database.NewSqlConnenction()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestFindByUserID(t *testing.T) {
	t.Run("正常系 取得結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		user, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createdTasks, err := createTasks(ctx, 5, user.ID)
		if err != nil {
			t.Fatal(err)
		}

		query := fmt.Sprintf("?limit=%d&offset=%d", 5, 0)
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		var tasks []domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&tasks)
		if err != nil {
			t.Fatal(err)
		}

		assertOrderByDueDate(t, tasks)

		mapTasks := map[int64]domain.Task{}
		for _, task := range tasks {
			mapTasks[task.ID] = task
		}
		for _, createdTask := range createdTasks {
			task := mapTasks[createdTask.ID]
			assert.Equal(t, createdTask.ID, task.ID)
			assert.Equal(t, createdTask.Title, task.Title)
			assert.Equal(t, createdTask.Content, task.Content)
			assert.True(t, createdTask.DueDate.Equal(task.DueDate))
		}
	})

	t.Run("正常系 limit, offsetを指定し、結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		user, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createdTasks, err := createTasks(ctx, 5, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		query := fmt.Sprintf("?limit=%d&offset=%d", 2, 1)
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)

		var tasks []domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&tasks)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 2, len(tasks))
		assertOrderByDueDate(t, tasks)

		sort.Slice(createdTasks, func(i, j int) bool {
			return createdTasks[i].DueDate.Before(createdTasks[j].DueDate)
		})
		expectedTasks := createdTasks[1:3]

		mapTasks := map[int64]domain.Task{}
		for _, task := range tasks {
			mapTasks[task.ID] = task
		}
		for _, expectedTask := range expectedTasks {
			task := mapTasks[expectedTask.ID]
			assert.Equal(t, expectedTask.ID, task.ID)
			assert.Equal(t, expectedTask.Title, task.Title)
			assert.Equal(t, expectedTask.Content, task.Content)
			assert.True(t, expectedTask.DueDate.Equal(task.DueDate))
		}
	})

	t.Run("正常系 存在しない場合、0件取得しステータス200であること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		query := fmt.Sprintf("?limit=%d&offset=%d", 5, 0)
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var tasks []domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&tasks)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 0, len(tasks))
	})

	t.Run("準正常系 トークンが指定されてない場合、401エラーとなること", func(t *testing.T) {
		query := fmt.Sprintf("?limit=%d&offset=%d", 5, 0)
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 トークンが間違っている場合、401エラーとなること", func(t *testing.T) {
		query := fmt.Sprintf("?limit=%d&offset=%d", 5, 0)
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid token")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 パラメータが指定されてない場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest("GET", taskURL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 パラメータの型が間違っている場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		query := fmt.Sprint("?limit=foo&offset=bar")
		req, _ := http.NewRequest("GET", taskURL+query, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		ctx := context.TODO()
		user, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createdTasks, err := createTasks(ctx, 1, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		createdTask := createdTasks[0]

		req, _ := http.NewRequest("GET", taskURL+"/"+strconv.Itoa(int(createdTask.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resTask domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resTask)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, createdTask.ID, resTask.ID)
		assert.Equal(t, createdTask.Title, resTask.Title)
		assert.Equal(t, createdTask.Content, resTask.Content)
	})

	t.Run("準正常系 存在しないIDで検索した際に404エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		taskID := int(time.Now().UnixNano())
		req, _ := http.NewRequest("GET", taskURL+"/"+strconv.Itoa(taskID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, resError.Message, domain.ErrRecordNotFound.Error())
	})

	t.Run("準正常系 トークンが指定されてない場合、401エラーとなること", func(t *testing.T) {
		req, _ := http.NewRequest("GET", taskURL+"/1", nil)
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 トークンが間違っている場合、401エラーとなること", func(t *testing.T) {
		req, _ := http.NewRequest("GET", taskURL+"/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid token")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)

	})

	t.Run("準正常系 指定されたIDが数字でない場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest("GET", taskURL+"/hoge", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})
}

func TestCreate(t *testing.T) {
	t.Run("正常系 1件作成し登録結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		user, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createRequest := taskHandler.CreateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", taskURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		repo := taskRepository.NewTaskRepository(sqlDriver)
		tasks, err := repo.FindByUserID(ctx, user.ID, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		task := tasks[0]

		s := createRequest.DueDate.String()
		d := task.DueDate.String()

		fmt.Println(s + d)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, createRequest.Title, task.Title)
		assert.Equal(t, createRequest.Content, task.Content)
		assert.True(t, createRequest.DueDate.Equal(task.DueDate))
	})

	t.Run("準正常系 トークンが指定されてない場合、401エラーとなること", func(t *testing.T) {
		createRequest := taskHandler.CreateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", taskURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 トークンが間違っている場合、401エラーとなること", func(t *testing.T) {
		createRequest := taskHandler.CreateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", taskURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid token")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createRequest := taskHandler.CreateTaskRequest{
			Title: "test title",
		}
		byteRequest, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", taskURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createRequest := domain.ErrorResponse{
			Message: "test",
		}
		byteRequest, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", taskURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("正常系 1件更新し結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		user, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		createdTask, err := createTasks(ctx, 1, user.ID)
		if err != nil {
			t.Fatal(err)
		}

		updateRequest := taskHandler.UpdateTaskRequest{
			Title:   "updated title",
			Content: "updated content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(updateRequest)

		req, _ := http.NewRequest("PUT", taskURL+"/"+strconv.Itoa(int(createdTask[0].ID)), bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		repo := taskRepository.NewTaskRepository(sqlDriver)
		updatedTask, err := repo.GetByID(ctx, createdTask[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, createdTask[0].ID, updatedTask.ID)
		assert.Equal(t, updateRequest.Title, updatedTask.Title)
		assert.Equal(t, updateRequest.Content, updatedTask.Content)
		assert.True(t, updateRequest.DueDate.Equal(updatedTask.DueDate))
	})

	t.Run("準正常系 存在しないIDを指定した際に404エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		updateRequest := taskHandler.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		taskID := int(time.Now().UnixNano())
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/"+strconv.Itoa(taskID), bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, resError.Message, domain.ErrRecordNotFound.Error())
	})

	t.Run("準正常系 指定されたIDが数字でない場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		updateRequest := taskHandler.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/hoge", bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			t.Fatal("成功レスポンスのためテスト失敗")
		}

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 トークンが指定されてない場合、401エラーとなること", func(t *testing.T) {
		updateRequest := taskHandler.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/1", bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid token")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 トークンが間違っている場合、401エラーとなること", func(t *testing.T) {
		updateRequest := taskHandler.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now().Round(time.Second),
		}
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/1", bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		updateRequest := taskHandler.UpdateTaskRequest{
			Title: "test title",
		}
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/123", bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {
		ctx := context.TODO()
		_, token, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}

		updateRequest := domain.ErrorResponse{
			Message: "test",
		}
		byteRequest, _ := json.Marshal(updateRequest)
		req, _ := http.NewRequest("PUT", taskURL+"/123", bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

func TestDelete(t *testing.T) {
	t.Run("正常系 1件削除し結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		user, _, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createdTasks, err := createTasks(ctx, 1, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		createdTask := createdTasks[0]

		req, _ := http.NewRequest("DELETE", taskURL+"/"+strconv.Itoa(int(createdTask.ID)), nil)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		repo := taskRepository.NewTaskRepository(sqlDriver)
		_, err = repo.GetByID(ctx, createdTask.ID)
		if err == nil {
			t.Fatal("エラーが起きていない場合失敗")
		}
		assert.Equal(t, http.StatusNoContent, response.StatusCode)
		assert.Equal(t, domain.ErrRecordNotFound, err)
	})

	t.Run("正常系 存在しないIDを指定した際にエラーとならないこと(204が返却されること)", func(t *testing.T) {
		taskID := time.Now().UnixNano()

		req, _ := http.NewRequest("DELETE", taskURL+"/"+strconv.Itoa(int(taskID)), nil)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		assert.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("準正常系 トークンが指定されてない場合、401エラーとなること", func(t *testing.T) {})

	t.Run("準正常系 トークンが間違っている場合、401エラーとなること", func(t *testing.T) {})

	t.Run("準正常系 指定されたIDが数字でない場合、400エラーとなること", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", taskURL+"/hoge", nil)
		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})
}

// createUser テストユーザーを作成し、ユーザーとトークンを返却します
func createUser(ctx context.Context) (domain.User, string, error) {
	email := fmt.Sprintf("%d@example.com", time.Now().UnixNano())
	userRepo := userRepository.NewUserRepository(sqlDriver)
	user := domain.User{
		Name:     "test user",
		Email:    email,
		Password: "test passsword",
		Salt:     "test salt",
	}

	userID, err := userRepo.Create(ctx, user)
	user.ID = userID

	token := token.GenerateAccessToken(user)

	return user, token, err
}

// createTasks テスト用のタスクを指定したユーザーIDで、指定された数作成します
func createTasks(ctx context.Context, num int, userID int64) ([]domain.Task, error) {
	repo := taskRepository.NewTaskRepository(sqlDriver)

	tasks := make([]domain.Task, 0)
	for i := 0; i < num; i++ {
		dueDate := time.Now().Add(time.Duration(i) * time.Hour)
		task := domain.Task{
			UserID:  userID,
			Title:   "test title" + strconv.Itoa(i+1),
			Content: "test content" + strconv.Itoa(i+1),
			DueDate: dueDate.Round(time.Second),
		}
		createdID, err := repo.Create(ctx, task)
		if err != nil {
			return nil, err
		}
		task.ID = createdID
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// assertOrderByDueDate タスクが期限日昇順になっているか確認します
func assertOrderByDueDate(t *testing.T, tasks []domain.Task) {
	isSorted := sort.SliceIsSorted(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})
	if !isSorted {
		t.Fatal("DueDate順になっていません")
	}
}
