package main

import (
	"core/internal"
	"core/internal/api"
	"core/internal/controller"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	err := internal.InitEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitialiseDB(&database.DbConfig{
		User:     internal.PostgresUser,
		Password: internal.PostgresPassword,
		DbName:   internal.PostgresDB,
		Host:     internal.PostgresHost,
		Port:     internal.PostgresPort,
		Schema:   "account",
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&database.ClientDB{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.CompanyDB{})
	if err != nil {
		panic(err)
	}

	clientRepository := repository.NewClientRepository(db)
	clientService := service.NewClientService(clientRepository)
	clientController := controller.NewClientController(clientService)

	v1 := r.Group("v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			clientController.LoginByToken(c)
		})
		v1.POST("/account", api.GetMockStandard)
		registerGroup := v1.Group("register")
		{
			registerGroup.POST("/client", func(c *gin.Context) {
				clientController.Signup(c)
			})
			registerGroup.POST("/company", api.GetMockStandard)
		}
	}

	r.Run()
}
