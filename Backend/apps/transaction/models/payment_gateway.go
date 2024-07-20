package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type Payment_Gateway struct {
	gorm.Model
	base.Base
	Slug               string             `json:"slug" gorm:"size:255;unique"`
	Label              string             `json:"label" gorm:"size:255"`
	TransactionRequest TransactionRequest `json:"transactionrequest" gorm:"foreignKey:Payment_Gateway_id"`
}
