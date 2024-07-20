package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	base.Base
	Name   string `json:"name" gorm:"size:255"`
	UserID int    `json:"userid" gorm:"unique"`
	// User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
