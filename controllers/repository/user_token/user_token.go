package user_token

import (
	"context"
	"rakamin-final-task/database"
	"rakamin-final-task/models"
)

type Interface interface {
	CheckTokenExist(ctx context.Context, params models.UserTokenParams) (int64, error)
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

func (u *userToken) CheckTokenExist(ctx context.Context, params models.UserTokenParams) (int64, error) {
	var isExist int64

	if err := u.db.ORM.WithContext(ctx).Model(models.UserToken{}).Where(params).Count(&isExist).Error; err != nil {
		return 0, err
	}

	return isExist, nil
}

func (u *userToken) Create(ctx context.Context, userToken models.UserToken) (models.UserToken, error) {
	if err := u.db.ORM.WithContext(ctx).Create(&userToken).Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}

func (u *userToken) Update(ctx context.Context, userToken models.UserToken, params models.UserTokenParams) (models.UserToken, error) {
	if err := u.db.ORM.WithContext(ctx).Model(models.UserToken{}).Where(params).Updates(&userToken).Error; err != nil {
		return userToken, err
	}

	return userToken, nil
}
