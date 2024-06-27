package models

import (
	"gorm.io/gorm"
)

type Photo struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *int64         `json:"createdBy"`
	UpdatedBy *int64         `json:"updatedBy"`
	DeletedBy *int64         `json:"deletedBy"`

	Title    string `gorm:"not null;type:varchar(255)" json:"title"`
	Caption  string `gorm:"not null;type:text" json:"caption"`
	PhotoURL string `gorm:"not null;type:text" json:"photoURL"`
	UserID   int64  `gorm:"not null" json:"userID"`
}
