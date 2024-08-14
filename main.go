package main

import (
	"go-rest/controller"
	"go-rest/db"
	"go-rest/repository"
	"go-rest/router"
	"go-rest/usecase"
)

// プログラムのエントリーポイント
func main() {
	db := db.NewDB()
	// dbインスタンスをrepositoryのコンストラクタに渡す
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
}
