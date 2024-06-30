package models

import (
	"gorm.io/gorm"
	"rakamin-final-task/helpers/response"
)

type Photos struct {
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

type PhotoParams struct {
	ID     int64 `json:"id" uri:"photo_id"`
	UserID int64 `json:"userID" uri:"user_id"`
	response.PaginationParam
}

type CreatePhotoParams struct {
	Title   string `json:"title" validate:"required"`
	Caption string `json:"caption" validate:"required"`
}
