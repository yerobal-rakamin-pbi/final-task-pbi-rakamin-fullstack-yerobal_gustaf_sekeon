package repository

import (
	userTokenRepo "rakamin-final-task/controllers/repository/user_token"
	userRepo "rakamin-final-task/controllers/repository/users"
	photoRepo "rakamin-final-task/controllers/repository/photos"
	"rakamin-final-task/database"
)

type Repository struct {
	Users     userRepo.Interface
	UserToken userTokenRepo.Interface
	Photos    photoRepo.Interface
}

func Init(db *database.DB) Repository {
	return Repository{
		Users:     userRepo.Init(db),
		UserToken: userTokenRepo.Init(db),
		Photos:    photoRepo.Init(db),
	}
}
