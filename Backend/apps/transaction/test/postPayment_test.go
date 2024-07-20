package test

import (
	accountmodels "GoTransact/apps/accounts/models"
	// "GoTransact/apps/accounts/validators"
	utils "GoTransact/apps/transaction"

	transactionmodels "GoTransact/apps/transaction/models"

	"GoTransact/pkg/db"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostPayment_InvalidAmount(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.Payment_Gateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput with validation errors
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327654",
		ExpiryDate:  "06/27",
		Cvv:         "147",
		Description: "Test payment",
		Amount:      "invalid_amount", // Invalid amount
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Contains(t, message, "error while validating")
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}

func TestPostPayment_InvalidCVV(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.Payment_Gateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput with validation errors
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327654",
		ExpiryDate:  "06/27",
		Cvv:         "1473",
		Description: "Test payment",
		Amount:      "74537.6",
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Contains(t, message, "error while validating")
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}

func TestPostPayment_Invalid_ExpDate(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.Payment_Gateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput with validation errors
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327654",
		ExpiryDate:  "06/21",
		Cvv:         "147",
		Description: "Test payment",
		Amount:      "6465.6",
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Contains(t, message, "error while validating")
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}

func TestPostPayment_Invalid_CardNumber(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.Payment_Gateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput with validation errors
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327",
		ExpiryDate:  "06/27",
		Cvv:         "147",
		Description: "Test payment",
		Amount:      "6343.6",
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Contains(t, message, "error while validating")
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}

// TestPostPayment_InvalidPaymentType tests invalid payment type scenario for PostPayment function.
func TestPostPayment_InvalidPaymentType(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)
	gateway := transactionmodels.Payment_Gateway{
		Slug:  "ach",
		Label: "ACH",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput with validation errors
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327385",
		ExpiryDate:  "06/27",
		Cvv:         "147",
		Description: "Test payment",
		Amount:      "1653.55",
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "invalid payment type", message)
	assert.Empty(t, data)

	CloseTestDb()
}

// TestPostPayment_Success tests successful payment processing for PostPayment function.
func TestPostPayment_Success(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.InitValidation()

	// Create a mock user
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	db.DB.Create(&user)

	// Create a mock payment gateway
	gateway := transactionmodels.Payment_Gateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	// Create a mock PostPaymentInput
	postPaymentInput := utils.PostPaymentInput{
		CardNumber:  "1234567854327654",
		ExpiryDate:  "06/26",
		Cvv:         "147",
		Description: "Test payment",
		Amount:      "6465.63",
	}

	status, message, data := utils.PostPayment(postPaymentInput, user)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "success", message)
	assert.NotEmpty(t, data["transaction ID"])

	CloseTestDb()
}
