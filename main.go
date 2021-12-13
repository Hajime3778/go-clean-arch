package main

import (
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	_taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/middleware"
	_taskHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/task"
	_taskUsecase "github.com/Hajime3778/go-clean-arch/usecase/task"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env.NewEnv().Load()
	sqlDriver := database.NewSqlConnenction()
	middleware := middleware.Middleware

	// タスクAPI
	taskRepository := _taskRepository.NewTaskRepository(sqlDriver)
	taskUsecase := _taskUsecase.NewTaskUsecase(taskRepository)

	// /tasks
	taskIndexHandler := _taskHandler.NewTaskIndexHandler(taskUsecase).Handler
	taskIndexHandlerFunc := http.HandlerFunc(taskIndexHandler)
	http.Handle(_taskHandler.TaskIndexPath, middleware(taskIndexHandlerFunc))

	// /tasks/:id
	taskHandler := _taskHandler.NewTaskHandler(taskUsecase).Handler
	taskHandlerFunc := http.HandlerFunc(taskHandler)
	http.Handle(_taskHandler.TaskPath, middleware(taskHandlerFunc))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
