package db

import (
	"GoTransact/config"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(typ string) {
	var dbURI string
	if strings.EqualFold(typ, "test") {
		fmt.Println("in test db setup")
		config.LoadEnv()
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, "testgotransact", config.DbPort, config.DbTimezone)
	} else {
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort, config.DbTimezone)
	}

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Microsecond, // Slow SQL threshold
	// 		LogLevel:                  logger.Info,      // Log level
	// 		IgnoreRecordNotFoundError: true,             // Ignore ErrRecord Not Found error for logger
	// 		Colorful:                  true,             // Disable color
	// 	},
	// )

	// db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{Logger: newLogger})
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db
}
