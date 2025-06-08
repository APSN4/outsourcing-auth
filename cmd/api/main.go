package main

import (
	"core/internal"
	"core/internal/api"
	"core/internal/controller"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/security"
	"core/internal/service"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type Entity interface {
	database.ClientDB | database.CompanyDB
}

func main() {
	r := gin.Default()
	
	// Загружаем HTML шаблоны
	r.LoadHTMLGlob("templates/*")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
	err = db.AutoMigrate(&database.Card{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.Order{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.EscrowTransaction{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.Review{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.Notification{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.BalanceTransaction{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.WorkerLink{})
	if err != nil {
		panic(err)
	}

	// Existing repositories and services
	clientRepository := repository.NewClientRepository(db)
	clientService := service.NewClientService(clientRepository)
	clientController := controller.NewClientController(clientService)
	companyRepository := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepository)
	companyController := controller.NewCompanyController(companyService)

	// New repositories
	cardRepository := repository.NewCardRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	balanceRepository := repository.NewBalanceRepository(db)
	escrowRepository := repository.NewEscrowRepository(db)
	workerLinkRepository := repository.NewWorkerLinkRepository(db)
	reviewRepository := repository.NewReviewRepository(db)

	// New services
	cardService := service.NewCardService(cardRepository)
	orderService := service.NewOrderService(orderRepository, cardRepository, balanceRepository, escrowRepository, workerLinkRepository)
	balanceService := service.NewBalanceService(balanceRepository)
	reviewService := service.NewReviewService(reviewRepository, orderRepository)

	// New controllers
	cardController := controller.NewCardController(cardService)
	orderController := controller.NewOrderController(orderService)
	balanceController := controller.NewBalanceController(balanceService)
	reviewController := controller.NewReviewController(reviewService)

	// Публичные маршруты (без авторизации)
	r.GET("/cards", cardController.GetAllCards)
	r.GET("/cards/category/:category", cardController.GetCardsByCategory)
	r.GET("/cards/:id", cardController.GetCardByID)
	r.GET("/cards/search", cardController.SearchCards)
	r.GET("/cards/price-range", cardController.GetCardsByPriceRange)
	r.GET("/orders", orderController.GetAllOrders)
	r.GET("/reviews/company/:company_id", reviewController.GetCompanyReviews)
	r.GET("/reviews/order/:order_id", reviewController.GetOrderReview)
	r.GET("/companies/:company_id/rating", reviewController.GetCompanyRating)

	// Специальная страница для работников (без авторизации)
	r.GET("/worker/complete/:token", orderController.CompleteOrderByWorker)

	v1 := r.Group("v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			request := &api.GeneralAuth{}
			if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
				api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
				return
			}
			resultClient, _, _ := clientRepository.ExistsByEmail(request.GeneralLogin.GeneralLoginAttributes.Email)
			if resultClient {
				clientController.Login(c, request)
			} else {
				companyController.Login(c, request)
			}
		})
		accountGroup := v1.Group("account")
		{
			accountGroup.POST("/", func(c *gin.Context) {
				request := &api.TokenAccess{}
				if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
					api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
					return
				}
				ok, mapClaims := security.CheckToken(request.User.Login.Token)
				if mapClaims == nil {
					api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
					return
				}
				if ok {
					isCompany := mapClaims["isCompany"].(bool)
					if isCompany {
						companyController.GetAccount(c, request)
					} else {
						clientController.GetAccount(c, request)
					}
				} else {
					api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
					return
				}
			})
			cardGroup := accountGroup.Group("card")
			{
				cardGroup.POST("/create", func(c *gin.Context) {
					request := &api.TokenCreateCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.CreateCard(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
				cardGroup.POST("/list", func(c *gin.Context) {
					request := &api.TokenListCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					limit := c.Query("limit")
					page := c.Query("page")
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.ListCard(c, request, limit, page)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
				cardGroup.POST("/delete", func(c *gin.Context) {
					request := &api.TokenDeleteCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.DeleteCard(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
				
				// Новые маршруты для карточек
				cardGroup.POST("/update", func(c *gin.Context) {
					request := &api.TokenUpdateCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							cardController.UpdateCard(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
			}

			// Группа для заказов
			orderGroup := accountGroup.Group("order")
			{
				orderGroup.POST("/create", func(c *gin.Context) {
					request := &api.TokenCreateOrder{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if !isCompany {
							orderController.CreateOrder(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only clients can create orders")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.POST("/pay", func(c *gin.Context) {
					request := &api.TokenOrderAction{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if !isCompany {
							orderController.PayOrder(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only clients can pay orders")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.POST("/start", func(c *gin.Context) {
					request := &api.TokenOrderAction{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							orderController.StartOrder(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only companies can start orders")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.POST("/finish", func(c *gin.Context) {
					request := &api.TokenOrderAction{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if !isCompany {
							orderController.FinishOrder(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only clients can finish orders")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.POST("/cancel", func(c *gin.Context) {
					request := &api.TokenOrderAction{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						orderController.CancelOrder(c, request)
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.POST("/list", func(c *gin.Context) {
					request := &api.TokenAccess{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							orderController.GetCompanyOrders(c, request)
						} else {
							orderController.GetClientOrders(c, request)
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				orderGroup.GET("/:id", orderController.GetOrderByID)
			}

			// Группа для баланса
			balanceGroup := accountGroup.Group("balance")
			{
				balanceGroup.POST("/", func(c *gin.Context) {
					request := &api.TokenAccess{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							balanceController.GetCompanyBalance(c, request)
						} else {
							balanceController.GetClientBalance(c, request)
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				balanceGroup.POST("/deposit", func(c *gin.Context) {
					request := &api.TokenDepositBalance{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if !isCompany {
							balanceController.DepositClientBalance(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only clients can deposit")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				balanceGroup.POST("/withdraw", func(c *gin.Context) {
					request := &api.TokenWithdrawBalance{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							balanceController.WithdrawCompanyBalance(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only companies can withdraw")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})

				balanceGroup.POST("/transactions", func(c *gin.Context) {
					request := &api.TokenAccess{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							balanceController.GetCompanyTransactions(c, request)
						} else {
							balanceController.GetClientTransactions(c, request)
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
			}

			// Группа для отзывов
			reviewGroup := accountGroup.Group("review")
			{
				reviewGroup.POST("/create", func(c *gin.Context) {
					request := &api.TokenCreateReview{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if !isCompany {
							reviewController.CreateReview(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "Only clients can create reviews")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
			}
		}
		registerGroup := v1.Group("register")
		{
			registerGroup.POST("/client", func(c *gin.Context) {
				clientController.Signup(c)
			})
			registerGroup.POST("/company", func(c *gin.Context) {
				companyController.Signup(c)
			})
		}
	}

	r.Run()
}
