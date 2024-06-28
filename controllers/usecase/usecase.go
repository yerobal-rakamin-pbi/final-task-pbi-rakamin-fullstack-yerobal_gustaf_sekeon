package usecase

import (
	"rakamin-final-task/controllers/repository"
)

type Usecase struct {
	Repository repository.Repository
}

func Init(repo repository.Repository) Usecase {
	return Usecase{
		Repository: repo,
	}
}