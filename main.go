package main

import (
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/middleware"
	taskHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/task"
	taskUsecase "github.com/Hajime3778/go-clean-arch/usecase/task"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env.NewEnv().Load()
	sqlDriver := database.NewSqlConnenction()
	middleware := middleware.Middleware

	// タスクAPI
	taskRepository := taskRepository.NewTaskRepository(sqlDriver)
	taskUsecase := taskUsecase.NewTaskUsecase(taskRepository)

	// /tasks
	taskIndexHandler := taskHandler.NewTaskIndexHandler(taskUsecase).Handler
	taskIndexHandlerFunc := http.HandlerFunc(taskIndexHandler)
	http.Handle(taskHandler.TaskIndexPath, middleware(taskIndexHandlerFunc))

	// /tasks/:id
	taskPathHandler := taskHandler.NewTaskHandler(taskUsecase).Handler
	taskPathHandlerFunc := http.HandlerFunc(taskPathHandler)
	http.Handle(taskHandler.TaskPath, middleware(taskPathHandlerFunc))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
