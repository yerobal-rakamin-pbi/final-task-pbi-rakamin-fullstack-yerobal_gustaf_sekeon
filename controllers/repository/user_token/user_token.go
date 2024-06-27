package user_token

import (
	"context"
	"rakamin-final-task/database"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/models"
)

type Interface interface {
	Get(ctx context.Context, params models.UserTokenParams) (models.UserToken, error)
	Create(ctx context.Context, userToken models.UserToken) (models.UserToken, error)
	Update(ctx context.Context, userToken models.UserToken, params models.UserTokenParams) (models.UserToken, error)
}

type userToken struct {
	db *database.DB
}

func Init(db *database.DB) Interface {
	return &userToken{
		db: db,
	}
}

func (u *userToken) Get(ctx context.Context, params models.UserTokenParams) (models.UserToken, error) {
	var userToken models.UserToken

	res := u.db.ORM.Where(params).First(&userToken)
	if res.Error != nil {
		return userToken, res.Error
	} else if res.RowsAffected == 0 {
		return userToken, errors.NotFound("Token not found")
	}

	return userToken, nil
}

func (u *userToken) Create(ctx context.Context, userToken models.UserToken) (models.UserToken, error) {
	if err := u.db.ORM.Create(&userToken).Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}

func (u *userToken) Update(ctx context.Context, userToken models.UserToken, params models.UserTokenParams) (models.UserToken, error) {
	if err := u.db.ORM.Model(models.UserToken{}).Where(params).Updates(&userToken).Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}
