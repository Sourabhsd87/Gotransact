package handlers

import (
	accountmodels "GoTransact/apps/accounts/models"
	basemodels "GoTransact/apps/base"
	utils "GoTransact/apps/transaction"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PaymentRequest handles the payment request
// @Summary Create a new payment request
// @Description Create a new payment request with the provided details
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header with bearer token"
// @Param paymentInput body utils.PostPaymentInput true "Payment Request Input"
// @Success 200 {object} basemodels.Response "Successfully created payment request"
// @Failure 400 {object} basemodels.Response "Invalid input"
// @Failure 500 {object} basemodels.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /protected/post-payment [post]
func PaymentRequest(c *gin.Context) {

	utils.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Post Payment Request received")

	var Postpaymentinput utils.PostPaymentInput
	if err := c.ShouldBindJSON(&Postpaymentinput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	UserFromRequest, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "User not found in token",
			Data:    map[string]interface{}{},
		})
		return
	}

	user, ok := UserFromRequest.(accountmodels.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	status, message, data := utils.PostPayment(Postpaymentinput, user)

	c.JSON(status, basemodels.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
