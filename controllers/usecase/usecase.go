package usecase

import (
	"rakamin-final-task/config"
	"rakamin-final-task/controllers/repository"
	userUsecase "rakamin-final-task/controllers/usecase/users"
	"rakamin-final-task/helpers/jwt"
	"rakamin-final-task/helpers/storage"
	"rakamin-final-task/helpers/validator"
)

type Usecase struct {
	Users userUsecase.Interface
}

type InitParam struct {
	Repo         repository.Repository
	ServerConf   config.Server
	JwtLib       jwt.Interface
	ValidatorLib validator.Interface
	StorageLib   storage.Interface
}

func Init(param InitParam) Usecase {
	userInitParam := userUsecase.InitParam{
		UserRepo:      param.Repo.Users,
		UserTokenRepo: param.Repo.UserToken,
		Config:        param.ServerConf,
		Jwt:           param.JwtLib,
		Validator:     param.ValidatorLib,
	}
	return Usecase{
		Users: userUsecase.Init(userInitParam),
	}
}
