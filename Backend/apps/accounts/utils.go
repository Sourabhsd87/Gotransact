package utils

import (
	// utils "GoTransact/apps/accounts"
	accountModels "GoTransact/apps/accounts/models"
	"GoTransact/pkg/db"
	
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SendMail(to string) {

	InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to send mail on registrtion to ", to)

	fmt.Println("start of mail")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "sourabhtrellis@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", "Registration successfull")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "YOU HAVE REGISTERED SUCCESSFULLY ON GOTRANSACT")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error while sending mail")
		panic(err)
	}
	InfoLogger.WithFields(logrus.Fields{}).Info("Registration mail sent to ", to)
}

type RegisterInput struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Companyname string `json:"companyName" binding:"required"`
	Password    string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}

var (
	secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	publicKey = secretKey.Public()                // DO share this one
	ctx       = context.Background()
	rdb       = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func CreateToken(user accountModels.User) (string, error) {

	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.Set("user", user)

	signedToken := token.V4Sign(secretKey, nil)

	return signedToken, nil
}

func VerifyToken(signedToken string) (any, error) {

	val, err := rdb.Get(ctx, signedToken).Result()
	if err == nil && val == "Blacklisted" {
		return nil, fmt.Errorf("token has been revoked")
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	verifiedtoken, err := parser.ParseV4Public(publicKey, signedToken, nil)
	if err != nil {
		return "", err
	}
	var User accountModels.User
	if err := verifiedtoken.Get("user", &User); err != nil {
		return "", err
	}
	return User, nil
}

// ValidatePassword checks if the password meets the complexity requirements
func ValidatePassword(fl validator.FieldLevel) bool {
	Password := fl.Field().String()

	var (
		hasMinLen    = len(Password) >= 8
		hasUpperCase = regexp.MustCompile(`[A-Z]`).MatchString(Password)
		hasLowerCase = regexp.MustCompile(`[a-z]`).MatchString(Password)
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString(Password)
		hasSpecial   = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(Password)
	)

	return hasMinLen && hasUpperCase && hasLowerCase && hasNumber && hasSpecial
}

var validate *validator.Validate

// Init initializes the custom validator
func InitValidation() {
	validate = validator.New()
	validate.RegisterValidation("password_complexity", ValidatePassword)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

func Login(loginuser LoginInput) (int, string, map[string]interface{}) {

	InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to login with ", loginuser.Email)

	// custom validator for additional password validation
	if err := GetValidator().Struct(loginuser); err != nil {

		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}

	var user accountModels.User
	if err := db.DB.Where("email = ?", loginuser.Email).First(&user).Error; err != nil {
		ErrorLogger.WithFields(logrus.Fields{}).Error("Failed to login with ", loginuser.Email)
		return http.StatusUnauthorized, "invalid username or password", map[string]interface{}{}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginuser.Password)); err != nil {
		ErrorLogger.WithFields(logrus.Fields{}).Error("Failed to login with ", loginuser.Email)
		return http.StatusUnauthorized, "invalid username or password", map[string]interface{}{}
	}

	token, err := CreateToken(user)
	if err != nil {

		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error creating token")

		return http.StatusInternalServerError, "error creating token", map[string]interface{}{"data": err.Error()}
	}

	InfoLogger.WithFields(logrus.Fields{}).Info("User logged in with ", loginuser.Email, " and id ", user.Internal_id)

	return http.StatusOK, "Logged in successfull", map[string]interface{}{"token": token}
}

func Signup(user RegisterInput) (int, string, map[string]interface{}) {

	InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to register with ", user.Email, " and company ", user.Companyname)

	if err := GetValidator().Struct(user); err != nil {
		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}

	//chaecking if user with email already exist
	var count int64

	// Check if user with the email already exists
	if err := db.DB.Model(&accountModels.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "email already exists", map[string]interface{}{}
	}

	// Check if company with the name already exists
	if err := db.DB.Model(&accountModels.Company{}).Where("name = ?", user.Companyname).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "company already exists", map[string]interface{}{}
	}

	//hashing the password to store in database
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return http.StatusInternalServerError, "Error while hashing password", map[string]interface{}{}
	}

	//creating user and company model
	newuser := accountModels.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
		Company: accountModels.Company{
			Name: user.Companyname,
		},
	}

	//save the user
	if err := db.DB.Create(&newuser).Error; err != nil {
		//log
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating user in database")

		return http.StatusInternalServerError, "error creating user", map[string]interface{}{}
	}
	//log
	InfoLogger.WithFields(logrus.Fields{}).Info("User created in database ", user.Email, " and company ", user.Companyname)

	go SendMail(user.Email)

	return http.StatusOK, "User created successfully", map[string]interface{}{}
}
