package main

import (
	"go-rest/controller"
	"go-rest/db"
	"go-rest/repository"
	"go-rest/router"
	"go-rest/usecase"

	"go-rest/validator"
)

// プログラムのエントリーポイント
func main() {
	db := db.NewDB()
	// dbインスタンスをrepositoryのコンストラクタに渡す
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	taskValidator := validator.NewTaskValidator()
	userValidator := validator.NewUserValidator()
	//usecaseにvalidatorをdiする
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
}
