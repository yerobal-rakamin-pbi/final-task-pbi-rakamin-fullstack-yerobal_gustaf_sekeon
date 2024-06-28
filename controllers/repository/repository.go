package repository

import (
	userTokenRepo "rakamin-final-task/controllers/repository/user_token"
	userRepo "rakamin-final-task/controllers/repository/users"
	"rakamin-final-task/database"
)

type Repository struct {
	Users     userRepo.Interface
	UserToken userTokenRepo.Interface
}

func Init(db *database.DB) Repository {
	return Repository{
		Users:     userRepo.Init(db),
		UserToken: userTokenRepo.Init(db),
	}
}
