package utils

import (
	accountmodels "GoTransact/apps/accounts/models"
	transactionmodels "GoTransact/apps/transaction/models"
	"GoTransact/pkg/db"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

type PostPaymentInput struct {
	CardNumber  string `json:"cardnumber" binding:"required" validate:"card_number" `
	ExpiryDate  string `json:"expirydate" binding:"required" validate:"expiry_date" `
	Cvv         string `json:"cvv" validate:"cvv" binding:"required"`
	Amount      string `json:"amount" binding:"required" validate:"amount"`
	Description string `json:"description" `
}

type TemplateData struct {
	Username     string
	TrasactionID uuid.UUID
	Amount       float64
	ConfirmURL   string
	CancelURL    string
	DateTime     time.Time
}

func SendMail(user accountmodels.User, request transactionmodels.TransactionRequest) {

	InfoLogger.WithFields(logrus.Fields{
		"email": user.Email,
		"id":    user.Internal_id,
	}).Info("Attempted to send confirm payment mail")

	fmt.Println("start of mail")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "sourabhtrellis@gmail.com")

	// Set E-Mail receivers
	// m.SetHeader("To", user.Email)
	m.SetHeader("To", user.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Payment Confirmation Required")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("/home/trellis/Sourabh/GoTransact/Backend/apps/transaction/templates/email_template.html")
	if err != nil {
		fmt.Printf("Error parsing email template: %s", err)
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	baseURL := "http://localhost:8080/api/confirm-payment" // Replace with your actual domain and endpoint
	params := url.Values{}
	params.Add("transaction_id", request.Internal_id.String())
	params.Add("status", "true")
	ConfirmActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	baseURL = "http://localhost:8080/api/confirm-payment" // Replace with your actual domain and endpoint
	params = url.Values{}
	params.Add("transaction_id", request.Internal_id.String())
	params.Add("status", "false")
	CancelActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Execute the template with the data
	TemplateData := TemplateData{
		Username:     user.FirstName,
		TrasactionID: request.Internal_id,
		Amount:       request.Amount,
		ConfirmURL:   ConfirmActionURL,
		CancelURL:    CancelActionURL,
	}
	fmt.Println(TemplateData)
	if err := tmpl.Execute(&body, TemplateData); err != nil {
		fmt.Printf("Error executing email template: %s", err)
	}

	fmt.Println()
	// Set E-Mail body as HTML
	m.SetBody("text/html", body.String())

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": user.Email,
		}).Error("Error sending confirm payment mail")
		panic(err)
	}
	InfoLogger.WithFields(logrus.Fields{
		"email": user.Email,
		"id":    user.Internal_id,
	}).Info("comfirmation mail sent")
}

func FetchTransactionsLast24Hours() []transactionmodels.TransactionRequest {
	var transactions []transactionmodels.TransactionRequest
	last24Hours := time.Now().Add(-24 * time.Hour)
	db.DB.Where("created_at >= ?", last24Hours).Find(&transactions)
	return transactions
}

func GenerateExcel(transactions []transactionmodels.TransactionRequest) (string, error) {
	f := excelize.NewFile()
	sheetName := "Transactions"
	index := f.NewSheet(sheetName)

	// f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "Transaction ID")
	// f.SetCellValue(sheetName, "C1", "UserID")
	f.SetCellValue(sheetName, "D1", "Status")
	// f.SetCellValue(sheetName, "E1", "PaymentGatewayID")
	f.SetCellValue(sheetName, "F1", "Description")
	f.SetCellValue(sheetName, "G1", "Transaction Amount")
	f.SetCellValue(sheetName, "H1", "Transaction date")
	// f.SetCellValue(sheetName, "I1", "UpdatedAt")

	for i, tr := range transactions {
		row := i + 2
		// f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tr.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tr.Internal_id)
		// f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tr.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tr.Status)
		// f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tr.Payment_Gateway_id)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), tr.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tr.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), tr.CreatedAt)
		// f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), tr.UpdatedAt)
	}

	f.SetActiveSheet(index)
	filePath := "transactions.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func SendMailWithAttachment(email, filePath string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "sourabhtrellis@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Daily Transactions Report")
	m.SetBody("text/plain", "Please find attached the daily transactions report.")
	m.Attach(filePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("could not send email: %v", err)
	}
	fmt.Println("Email sent successfully")
}

var validate *validator.Validate

func cardNumberValidator(fl validator.FieldLevel) bool {
	fmt.Println("-------in card validator--------")
	cardNumber := fl.Field().String()
	// Check if card number is 16 or 18 digits
	match, _ := regexp.MatchString(`^\d{16}|\d{18}$`, cardNumber)
	fmt.Println(match)
	return match
}

func expiryDateValidator(fl validator.FieldLevel) bool {
	expiryDate := fl.Field().String()
	// Check if expiry date is in the format MM/YY and within 10 years span
	t, err := time.Parse("01/06", expiryDate)
	if err != nil {
		return false
	}
	currentYear := time.Now().Year() % 100
	expiryYear := t.Year() % 100
	return expiryYear >= currentYear
}

func cvvValidator(fl validator.FieldLevel) bool {
	cvv := fl.Field().String()
	// Check if CVV is exactly 3 digits
	match, _ := regexp.MatchString(`^\d{3}$`, cvv)
	return match
}

func amountValidation(fl validator.FieldLevel) bool {
	amount := fl.Field().String()
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false
	}
	return value > 0
}

// CustomErrorMessages contains custom error messages for validation
var CustomErrorMessages = map[string]string{
	"card_number": "Card number must be 16 or 18 digits.",
	"expiry_date": "Expiry date must be in MM/YY format and within a 10 year span.",
	"cvv":         "CVV must be exactly 3 digits.",
	"amount":      "Amount must be greater than 0.",
}

// InitValidation initializes the custom validators

func InitValidation() {
	validate = validator.New()
	validate.RegisterValidation("card_number", cardNumberValidator)
	validate.RegisterValidation("expiry_date", expiryDateValidator)
	validate.RegisterValidation("cvv", cvvValidator)
	validate.RegisterValidation("amount", amountValidation)
}
func GetValidator() *validator.Validate {
	return validate
}

func PostPayment(Postpaymentinput PostPaymentInput, user accountmodels.User) (int, string, map[string]interface{}) {

	InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to create transaction request with email ", user.Email, " id ", user.Internal_id)

	if err := GetValidator().Struct(Postpaymentinput); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			tag := fieldErr.Tag()
			errors[fieldName] = CustomErrorMessages[tag]
		}
		return http.StatusBadRequest, "error while validating", map[string]interface{}{}
	}

	floatAmount, _ := strconv.ParseFloat(Postpaymentinput.Amount, 64)

	var gateway transactionmodels.Payment_Gateway
	if err := db.DB.Where("slug = ?", "card").First(&gateway).Error; err != nil {
		return http.StatusBadRequest, "invalid payment type", map[string]interface{}{}
	}

	TransactionRequest := transactionmodels.TransactionRequest{
		UserID:             user.ID,
		Status:             transactionmodels.StatusProcessing,
		Description:        Postpaymentinput.Description,
		Amount:             floatAmount,
		Payment_Gateway_id: gateway.ID,
	}

	if err := db.DB.Create(&TransactionRequest).Error; err != nil {
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction request transaction-id=", TransactionRequest.Internal_id)
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}
	InfoLogger.WithFields(logrus.Fields{}).Info("created record in transaction request with email ", user.Email, " id ", user.Internal_id)

	TransactionHistory := transactionmodels.TransactionHistory{
		TransactionID: TransactionRequest.ID,
		Status:        TransactionRequest.Status,
		Description:   TransactionRequest.Description,
		Amount:        TransactionRequest.Amount,
	}

	if err := db.DB.Create(&TransactionHistory).Error; err != nil {
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction history")
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}

	InfoLogger.WithFields(logrus.Fields{}).Info("Created record in transaction history with email ", user.Email, " id ", user.Internal_id)

	go SendMail(user, TransactionRequest)

	return http.StatusOK, "success", map[string]interface{}{"transaction ID": TransactionRequest.Internal_id}
}

func ConfirmPayment(transactionIdStr, statusStr string) (int, string, map[string]interface{}) {
	// Convert the string ID to a uuid.UUID
	InfoLogger.WithFields(logrus.Fields{
		// "email": user.Email,
		// "id":    user.Internal_id,
	}).Info("Attempted to confirm/cancel payment transaction-id=", transactionIdStr)
	transactionId, err := uuid.Parse(transactionIdStr)
	fmt.Println("parsed", transactionId)
	if err != nil {
		return http.StatusBadRequest, "Invalid transaction ID", map[string]interface{}{}
	}

	var transactionRequest transactionmodels.TransactionRequest
	if err := db.DB.Where("internal_id = ?", transactionId).First(&transactionRequest).Error; err != nil {
		return http.StatusBadRequest, "transaction request not found", map[string]interface{}{}
	}

	var trasactionHistory transactionmodels.TransactionHistory
	trasactionHistory.TransactionID = transactionRequest.ID
	trasactionHistory.Description = transactionRequest.Description
	trasactionHistory.Amount = transactionRequest.Amount

	if strings.EqualFold(statusStr, "true") {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transactionmodels.StatusSuccess).Error; err != nil {
			ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error completing payment transaction-id=", transactionRequest.Internal_id)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		InfoLogger.WithFields(logrus.Fields{}).Info("Payment completed transaction-id=", transactionRequest.Internal_id)
		{
			trasactionHistory.Status = transactionmodels.StatusSuccess
		}
	} else {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transactionmodels.StatusFailed).Error; err != nil {
			ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error canceling the payment transaction-id=", transactionRequest.Internal_id)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		InfoLogger.WithFields(logrus.Fields{}).Info("Payment canceled transaction-id=", transactionRequest.Internal_id)
		{
			trasactionHistory.Status = transactionmodels.StatusFailed
		}

	}

	if err := db.DB.Create(&trasactionHistory).Error; err != nil {
		ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to update transaction history to confirm/cancel transaction-id=", transactionRequest.Internal_id)
		return http.StatusInternalServerError, "Failed to update transaction history", map[string]interface{}{}
	}
	InfoLogger.WithFields(logrus.Fields{}).Info("transaction history updated to confirm/cancel transaction-id=", transactionRequest.Internal_id)
	if strings.EqualFold(statusStr, "true") {
		return http.StatusOK, "Transaction successful", map[string]interface{}{"Amount": transactionRequest.Amount}
	} else {
		return http.StatusOK, "Transaction Canceled", map[string]interface{}{"Amount": transactionRequest.Amount}
	}
}
