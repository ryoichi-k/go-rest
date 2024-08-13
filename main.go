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
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)
	e := router.NewRouter(userController)
	e.Logger.Fatal(e.Start(":8080"))
}
