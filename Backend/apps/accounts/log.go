package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
)

func AccountLogInit() {
	// Create the info logger
	InfoLogger = logrus.New()
	infoFile, err := os.OpenFile("./apps/accounts/logs/infolog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		InfoLogger.Out = infoFile
	} else {
		InfoLogger.Info("Failed to log to file, using default stderr" + err.Error())
	}

	// Create the error logger
	ErrorLogger = logrus.New()
	errorFile, err := os.OpenFile("./apps/accounts/logs/errorlog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		ErrorLogger.Out = errorFile
	} else {
		ErrorLogger.Error("Failed to log to file, using default stderr" + err.Error())
	}

	InfoLogger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	InfoLogger.SetReportCaller(true)
	ErrorLogger.SetReportCaller(true)

	// Set log level
	InfoLogger.SetLevel(logrus.InfoLevel)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
