package main

import (
	"medods-internship/config"
	"medods-internship/controllers"
	"medods-internship/models"
	"medods-internship/router"
	"medods-internship/services"
	"medods-internship/storage"
	"medods-internship/utils"
)

func main() {
	// читаем конфиг
	appConfig, err := config.LoadEnv()
	utils.LogExit(err)

	// коннектимся к датабазе
	db, err := config.ConnectToDB(appConfig)
	utils.LogExit(err)
	db.AutoMigrate(&models.User{}, &models.RefreshToken{})

	// создаем хранилище
	storage := storage.NewStorage(db)

	// создаем сервис
	tokenService := services.NewTokenService(appConfig)

	// создаем контроллер
	tokenController := controllers.NewTokenController(storage, tokenService, appConfig)

	// создаем роутер
	authRouter := router.Init(tokenController)

	// запускаем сервер
	utils.LogExit(authRouter.Run(appConfig.ServerHost + ":" + appConfig.ServerPort))
}