package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wanliqun/go-wallet-app/controllers"
	"gorm.io/gorm"
)

func SetupRouter(router *gin.Engine, db *gorm.DB) {
	walletController := controllers.NewWalletController(db)

	wallet := router.Group("/wallet")
	{
		wallet.POST("/deposit", walletController.Deposit)
		wallet.POST("/withdraw", walletController.Withdraw)
		wallet.POST("/transfer", walletController.Transfer)
		wallet.GET("/balances", walletController.GetBalances)
		wallet.GET("/transactions", walletController.GetTransactionHistory)
	}
}
