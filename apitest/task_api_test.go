package apitest_test

import (
	"context"
	"encoding/json"
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
	"github.com/stretchr/testify/assert"
)

const taskURL = "http://localhost:8080/tasks"

var sqlDriver interfaceDB.SqlDriver

func TestMain(m *testing.M) {
	env.NewEnv().LoadEnvFile("../.env")
	sqlDriver = database.NewSqlConnenction()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestFindByUserID(t *testing.T) {
	t.Run("正常系 5件取得し結果が正しいこと", func(t *testing.T) {
		ctx := context.TODO()
		taskRepo := taskRepository.NewTaskRepository(sqlDriver)
		userID, err := createUser(ctx)
		if err != nil {
			t.Fatal(err)
		}
		createdTasks, err := createTasks(5, userID)
		if err != nil {
			t.Fatal(err)
		}
		tasks, err := taskRepo.FindByUserID(ctx, userID, 5, 0)
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
}

// createUser テストユーザーを作成し、ユーザーIDを返却します
func createUser(ctx context.Context) (int64, error) {
	userRepo := userRepository.NewUserRepository(sqlDriver)
	user := domain.User{
		Name:     "test user",
		Email:    "test@example.com",
		Password: "test passsword",
		Salt:     "test salt",
	}
	return userRepo.Create(ctx, user)
}

// createTasks テスト用のタスクを指定したユーザーIDで、指定された数作成します
func createTasks(num int, userID int64) ([]domain.Task, error) {
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
		createdID, err := repo.Create(context.TODO(), task)
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

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		repo := taskRepository.NewTaskRepository(sqlDriver)
		createTask := domain.Task{
			UserID:  1,
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		createdID, err := repo.Create(context.TODO(), createTask)
		if err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequest("GET", taskURL+"/"+strconv.Itoa(int(createdID)), nil)
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
		assert.Equal(t, createdID, resTask.ID)
		assert.Equal(t, createTask.Title, resTask.Title)
		assert.Equal(t, createTask.Content, resTask.Content)
	})

	t.Run("準正常系 存在しないIDで検索した際に404エラーとなること", func(t *testing.T) {
		taskID := int(time.Now().UnixNano())
		req, _ := http.NewRequest("GET", taskURL+"/"+strconv.Itoa(taskID), nil)
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
}

func TestCreate(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}
