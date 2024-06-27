package repository

import (
	userTokenRepo "rakamin-final-task/controllers/repository/user_token"
	userRepo "rakamin-final-task/controllers/repository/users"
	"rakamin-final-task/database"
)

type Repository struct {
	User      userRepo.Interface
	UserToken userTokenRepo.Interface
}

func Init(db *database.DB) Repository {
	return Repository{
		User:      userRepo.Init(db),
		UserToken: userTokenRepo.Init(db),
	}
}
