package models

import (
	"github.com/google/uuid"
)

type Base struct {
	Internal_id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"internal_id"`
	IsActive    bool      `gorm:"default:true" json:"isactive"`
}
