package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *int64         `json:"createdBy"`
	UpdatedBy *int64         `json:"updatedBy"`
	DeletedBy *int64         `json:"deletedBy"`

	Username string  `gorm:"not null;unique;type:varchar(255)" json:"username"`
	Email    string  `gorm:"not null;unique;type:varchar(255)" json:"email"`
	Password string  `gorm:"not null;type:text" json:"-"`
	Photos   []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"photos"`
}
