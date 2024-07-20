package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusProcessing TransactionStatus = "processing"
	StatusSuccess    TransactionStatus = "success"
	StatusFailed     TransactionStatus = "failed"
)

type TransactionRequest struct {
	gorm.Model
	base.Base
	UserID             uint               `json:"user_id" gorm:""`
	Status             TransactionStatus  `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Payment_Gateway_id uint               `json:"payment_gateway_id" gorm:""`
	Description        string             `json:"description" gorm:"size:255"`
	Amount             float64            `json:"amount" gorm:"type:float"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
