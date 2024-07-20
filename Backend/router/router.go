package router

import (
	account "GoTransact/apps/accounts/handlers"
	transaction "GoTransact/apps/transaction/handlers"
	"GoTransact/middlewares"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	r := gin.Default()
	// docs.SwaggerInfo.BasePath = "/api"

	r.POST("/api/register", account.Signup_handler)
	r.POST("/api/login", account.Login_handler)
	r.GET("/api/confirm-payment", transaction.ConfirmPayment)

	protected := r.Group("/api/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/post-payment", transaction.PaymentRequest)
		protected.POST("/logout", account.LogoutHandler)
	}

	return r
}
