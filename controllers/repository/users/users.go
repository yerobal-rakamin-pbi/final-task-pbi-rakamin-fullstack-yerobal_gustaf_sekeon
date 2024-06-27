package users

import (
	"context"
	"rakamin-final-task/database"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/models"
)

type Interface interface {
	Get(ctx context.Context, params models.UserParams) (models.Users, error)
	Create(ctx context.Context, user models.Users) (models.Users, error)
	Update(ctx context.Context, user models.Users, params models.UserParams) (models.Users, error)
}

type user struct {
	db         *database.DB
}

func Init(db *database.DB) Interface {
	return &user{
		db:         db,
	}
}

func (u *user) Get(ctx context.Context, params models.UserParams) (models.Users, error) {
	var user models.Users

	res := u.db.ORM.Where(params).First(&user)
	if res.Error != nil {
		return user, res.Error
	} else if res.RowsAffected == 0 {
		return user, errors.NotFound("Email or username not found")
	}

	return user, nil
}

func (u *user) Create(ctx context.Context, user models.Users) (models.Users, error) {
	if err := u.db.ORM.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *user) Update(ctx context.Context, user models.Users, params models.UserParams) (models.Users, error) {
	if err := u.db.ORM.Model(models.Users{}).Where(params).Updates(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
