package main

import (
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	userRepository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	authHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/auth"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/middleware"
	taskHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/task"
	authUsecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
	taskUsecase "github.com/Hajime3778/go-clean-arch/usecase/task"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env.NewEnv().Init()
	sqlDriver := database.NewSqlConnenction()
	middleware := middleware.Middleware

	// 認証API
	userRepository := userRepository.NewUserRepository(sqlDriver)
	authUsecase := authUsecase.NewAuthUsecase(userRepository)

	// /auth/sign_up
	signUpHandler := authHandler.NewAuthHandler(authUsecase).SignUpHandler
	signUpHandlerFunc := http.HandlerFunc(signUpHandler)
	http.Handle(authHandler.SignUpPath, middleware(signUpHandlerFunc))

	// /auth/sign_in
	signInHandler := authHandler.NewAuthHandler(authUsecase).SignInHandler
	signInHandlerFunc := http.HandlerFunc(signInHandler)
	http.Handle(authHandler.SignInPath, middleware(signInHandlerFunc))

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
