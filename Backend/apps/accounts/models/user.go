package models

import (
	base "GoTransact/apps/base"
	"GoTransact/apps/transaction/models"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	base.Base
	FirstName          string                    `json:"firstname" gorm:"size:255"`
	LastName           string                    `json:"lastname" gorm:"size:255"`
	Email              string                    `json:"email" gorm:"size:255;unique" `
	Password           string                    `json:"password" gorm:"size:255"`
	Company            Company                   `json:"company" gorm:"foreignKey:UserID"`
	TransactionRequest models.TransactionRequest `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
