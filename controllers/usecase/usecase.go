package usecase

import (
	"rakamin-final-task/config"
	"rakamin-final-task/controllers/repository"
	userUsecase "rakamin-final-task/controllers/usecase/users"
	"rakamin-final-task/helpers/jwt"
)

type Usecase struct {
	Users userUsecase.Interface
}

type InitParam struct {
	Repo       repository.Repository
	ServerConf config.Server
	JwtLib     jwt.Interface
}

func Init(param InitParam) Usecase {
	return Usecase{
		Users: userUsecase.Init(param.Repo.Users, param.Repo.UserToken, param.ServerConf, param.JwtLib),
	}
}
