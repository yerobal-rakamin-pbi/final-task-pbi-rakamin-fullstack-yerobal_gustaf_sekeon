package models

import (
	"gorm.io/gorm"
	"rakamin-final-task/helpers/response"
)

type Users struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *int64         `json:"createdBy"`
	UpdatedBy *int64         `json:"updatedBy"`
	DeletedBy *int64         `json:"deletedBy"`

	Username  string   `gorm:"not null;unique;type:varchar(255)" json:"username"`
	Email     string   `gorm:"not null;unique;type:varchar(255)" json:"email"`
	Password  string   `gorm:"not null;type:text" json:"-"`
	IsActived *bool    `gorm:"default:true" json:"isActived"`
	Photos    []Photos `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type UserParams struct {
	ID       int64  `json:"id" uri:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	response.PaginationParam
}

type UserLoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"min=6"`
}

type AuthResponse struct {
	User       Users  `json:"user"`
	AcessToken string `json:"accessToken"`
}

type UserToken struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *int64         `json:"createdBy"`
	UpdatedBy *int64         `json:"updatedBy"`
	DeletedBy *int64         `json:"deletedBy"`

	UserID      int64  `gorm:"not null;index:user_id_access_token_idx,unique" json:"userID"`
	AccessToken string `gorm:"not null;index:user_id_access_token_idx,unique;type:text" json:"accessToken"`
	IsRevoked   *bool  `gorm:"default:false" json:"isRevoked"`
}

type UserTokenParams struct {
	UserID      int64  `json:"userID"`
	AccessToken string `json:"accessToken"`
	IsRevoked   *bool  `json:"isRevoked"`
}
