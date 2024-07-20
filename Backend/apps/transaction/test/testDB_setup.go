package test

import (
	"GoTransact/apps/accounts/models"
	transactionmodel "GoTransact/apps/transaction/models"
	"GoTransact/pkg/db"

	"fmt"
)

func SetupTestDb() {
	fmt.Println("----------------in setup test db----------")
	db.InitDB("test")
	if err := db.DB.AutoMigrate(&models.User{}, &models.Company{}, &transactionmodel.Payment_Gateway{}, &transactionmodel.TransactionRequest{}, &transactionmodel.TransactionHistory{}); err != nil {
		fmt.Printf("Error automigrating models : %s", err.Error())
	}
}

func CloseTestDb() {
	fmt.Println("---------------in close db-------------")
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Printf("Error getting sqlDB from gorm DB: %s", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Printf("Error closing database: %s", err.Error())
	}
}

func ClearDatabase() {
	fmt.Println("----------------in clear db-----------------")

	db.DB.Exec("DELETE FROM companies")

	db.DB.Exec("DELETE FROM users")

	db.DB.Exec("DELETE FROM transaction_histories")

	db.DB.Exec("DELETE FROM transaction_requests")
	db.DB.Exec("DELETE FROM payment_gateways")
}
