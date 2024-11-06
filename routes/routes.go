package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wanliqun/go-wallet-app/controllers"
	"github.com/wanliqun/go-wallet-app/middlewares"
	"github.com/wanliqun/go-wallet-app/services"
	"gorm.io/gorm"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	walletService := services.NewWalletService(db)
	userService := services.NewUserService(db)

	router.Use(middlewares.AuthMiddleware(userService))
	router.Use(middlewares.CorsMiddleware())

	walletController := controllers.NewWalletController(walletService, userService)
	walletRouter := router.Group("/wallet")
	{
		walletRouter.POST("/deposit", walletController.Deposit)
		walletRouter.POST("/withdraw", walletController.Withdraw)
		walletRouter.POST("/transfer", walletController.Transfer)
		walletRouter.GET("/balances", walletController.GetBalances)
		walletRouter.GET("/transactions", walletController.GetTransactionHistory)
	}
}
