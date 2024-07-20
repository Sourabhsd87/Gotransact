package router

import (
	account "GoTransact/apps/accounts/handlers"
	transaction "GoTransact/apps/transaction/handlers"
	docs "GoTransact/docs"
	"GoTransact/middlewares"

	"github.com/gin-gonic/gin"
)

func Router2() *gin.Engine {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api"

	api := r.Group("/api")
	{
		api.POST("/register", account.Signup_handler)
		api.POST("/login", account.Login_handler)
		api.GET("/confirm-payment", transaction.ConfirmPayment)

		protected := api.Group("/protected")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/post-payment", transaction.PaymentRequest)
			protected.POST("/logout", account.LogoutHandler)
		}
	}
	return r
}
