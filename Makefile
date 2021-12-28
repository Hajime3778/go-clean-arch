BINARY=engine

task_test:
	make task_repository_test task_usecase_test

task_repository_test: 
	go test -v -cover -covermode=atomic ./interface/database/task

task_usecase_test: 
	go test -v -cover -covermode=atomic ./usecase/task