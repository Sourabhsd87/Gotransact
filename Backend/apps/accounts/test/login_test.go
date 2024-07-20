package test

import (
	utils "GoTransact/apps/accounts"
	
	"GoTransact/apps/accounts/models"

	
	"GoTransact/pkg/db"
	log "GoTransact/settings"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_success(t *testing.T) {
	fmt.Println("---------------------------------------in TestLogin_success")
	SetupTestDb()
	ClearDatabase()
	utils.Init()
	log.Init()
	//create a user
	password, _ := utils.HashPassword("Password@123")
	existinguser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	input := utils.LoginInput{
		Email:    "test@gmail.com",
		Password: "Password@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Logged in successfull", message)
	assert.NotEmpty(t, data["token"])
	ClearDatabase()
	CloseTestDb()
}

func TestLogin_InvalidEmail(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.Init()
	log.Init()
	password, _ := utils.HashPassword("Password@123")
	existinguser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	// Attempt to log in with invalid email
	input := utils.LoginInput{
		Email:    "wrongtestemail@example.com",
		Password: "Password@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, "invalid username or password", message)
	assert.Empty(t, data)
	CloseTestDb()
}

func TestLogin_InvalidPassword(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	utils.Init()
	log.Init()
	// Create a user
	password, _ := utils.HashPassword("Password@123")
	existinguser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	input := utils.LoginInput{
		Email:    "test@gmail.com",
		Password: "WrongPassword@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, "invalid username or password", message)
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}
