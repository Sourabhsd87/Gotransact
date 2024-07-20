package main

//go:generate swagger generate spec -o docs/swagger.json
import (
	accountModels "GoTransact/apps/accounts/models"
	transactionModels "GoTransact/apps/transaction/models"

	accountsUtils "GoTransact/apps/accounts"
	transactionUtils "GoTransact/apps/transaction"

	"GoTransact/config"
	db "GoTransact/pkg/db"
	"GoTransact/router"
	"log"

	// log "GoTransact/pkg/log"

	_ "GoTransact/docs"

	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api
// @title GoTransact
// @version 1.0
// @description This is a sample server for a project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081

func main() {
	config.LoadEnv()
	db.InitDB("prod")
	accountsUtils.InitValidation()
	transactionUtils.InitValidation()
	accountsUtils.AccountLogInit()
	transactionUtils.TransactionLogInit()

	
	db.DB.AutoMigrate(&accountModels.User{}, &accountModels.Company{}, &transactionModels.Payment_Gateway{}, &transactionModels.TransactionRequest{}, &transactionModels.TransactionHistory{})

	c := cron.New()
	c.AddFunc("@every 24h", func() {
		transactions := transactionUtils.FetchTransactionsLast24Hours()
		filePath, err := transactionUtils.GenerateExcel(transactions)
		if err != nil {
			log.Fatalf("failed to generate excel: %v", err)
		}
		transactionUtils.SendMailWithAttachment("sourabhsd87@gmail.com", filePath)
	})
	c.Start()

	r := router.Router2()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")

}
