package users

import (
	"context"
	"strings"

	"rakamin-final-task/config"
	userTokenRepo "rakamin-final-task/controllers/repository/user_token"
	userRepo "rakamin-final-task/controllers/repository/users"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/helpers/jwt"
	"rakamin-final-task/helpers/password"
	"rakamin-final-task/helpers/validator"
	"rakamin-final-task/models"
)

type Interface interface {
	Login(ctx context.Context, params models.UserLoginParams) (models.AuthResponse, error)
	Register(ctx context.Context, params models.UserRegisterParams) (models.AuthResponse, error)
	// GetUserProfile(ctx context.Context) (models.Users, error)
	// UpdateUserProfile(ctx context.Context, user models.Users, params models.UserParams) (models.Users, error)
	// DeactivateUser(ctx context.Context, params models.UserParams) (models.Users, error)
}

type users struct {
	user      userRepo.Interface
	userToken userTokenRepo.Interface
	config    config.Server
	jwt       jwt.Interface
	validator validator.Interface
}

func Init(userRepo userRepo.Interface, userTokenRepo userTokenRepo.Interface, config config.Server, jwt jwt.Interface, validator validator.Interface) Interface {
	return &users{
		user:      userRepo,
		userToken: userTokenRepo,
		config:    config,
		jwt:       jwt,
		validator: validator,
	}
}

func (u *users) Login(ctx context.Context, params models.UserLoginParams) (models.AuthResponse, error) {
	var res models.AuthResponse

	emailParam := models.UserParams{
		Email: params.Email,
	}
	userRes, err := u.user.Get(ctx, emailParam)
	if err != nil {
		return res, err
	}

	if userRes.ID == 0 {
		return res, errors.NotFound("Email not found")
	}

	if !password.Compare(userRes.Password, params.Password) {
		return res, errors.Unauthorized("Wrong password")
	}

	accessToken, err := u.jwt.GenerateToken(userRes)
	if err != nil {
		return res, err
	}

	userToken := models.UserToken{
		UserID:      userRes.ID,
		AccessToken: accessToken,
	}

	_, err = u.userToken.Create(ctx, userToken)
	if err != nil {
		return res, err
	}

	res.User = userRes
	res.AcessToken = accessToken

	return res, nil
}

func (u *users) Register(ctx context.Context, param models.UserRegisterParams) (models.AuthResponse, error) {
	var res models.AuthResponse

	if err := u.validator.ValidateStruct(param); err != nil {
		return res, errors.ValidationError(u.validator.GetValidationErrors(err))
	}

	hashedPassword, err := password.Hash(param.Password, u.config.Password.SaltRound)
	if err != nil {
		return res, err
	}

	user := models.Users{
		Username: param.Username,
		Email:    param.Email,
		Password: hashedPassword,
	}

	userRes, err := u.user.Create(ctx, user)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return res, errors.Conflict("Email or Username already exists")
	} else if err != nil {
		return res, err
	}

	accessToken, err := u.jwt.GenerateToken(userRes)
	if err != nil {
		return res, err
	}

	userToken := models.UserToken{
		UserID:      userRes.ID,
		AccessToken: accessToken,
	}

	_, err = u.userToken.Create(ctx, userToken)
	if err != nil {
		return res, err
	}

	res.User = userRes
	res.AcessToken = accessToken

	return res, nil
}
