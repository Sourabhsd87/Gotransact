package handlers

import (
	// "GoTransact/apps/accounts/functions"
	// "GoTransact/apps/accounts/utils"
	utils "GoTransact/apps/accounts"
	basemodels "GoTransact/apps/base"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /api
// @Summary 			Login
// @Description 		User login
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Param 				loginInput body   utils.LoginInput true "Login input"
// @in 					header
// @Success 			200 {object} basemodels.Response
// @Failure 			400 {object} basemodels.Response
// @Failure 			401 {object} basemodels.Response
// @Router 				/login [post]
func Login_handler(c *gin.Context) {

	utils.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Login Request received")

	var loginInput utils.LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {

		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	status, message, data := utils.Login(loginInput)

	c.JSON(status, basemodels.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
