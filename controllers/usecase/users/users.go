package users

import (
	"context"
	"strings"

	"rakamin-final-task/config"
	userTokenRepo "rakamin-final-task/controllers/repository/user_token"
	userRepo "rakamin-final-task/controllers/repository/users"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/helpers/jwt"
	"rakamin-final-task/helpers/password"
	"rakamin-final-task/helpers/validator"
	"rakamin-final-task/models"
)

type Interface interface {
	Login(ctx context.Context, params models.UserLoginParams) (models.AuthResponse, error)
	Register(ctx context.Context, params models.UserRegisterParams) (models.AuthResponse, error)
	CheckUserToken(ctx context.Context, token string) (string, bool)
	UpdateUser(ctx context.Context, body models.UpdateUserParams, params models.UserParams) (models.Users, error)
	GetUserProfile(ctx context.Context) (models.Users, error)
	DeactivateUser(ctx context.Context, params models.UserParams) (models.Users, error)
}

type users struct {
	user      userRepo.Interface
	userToken userTokenRepo.Interface
	config    config.Server
	jwt       jwt.Interface
	validator validator.Interface
}

type InitParam struct {
	UserRepo      userRepo.Interface
	UserTokenRepo userTokenRepo.Interface
	Config        config.Server
	Jwt           jwt.Interface
	Validator     validator.Interface
}

func Init(param InitParam) Interface {
	return &users{
		user:      param.UserRepo,
		userToken: param.UserTokenRepo,
		config:    param.Config,
		jwt:       param.Jwt,
		validator: param.Validator,
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
		validationErr, _ := u.validator.GetValidationErrors(err)
		return res, errors.ValidationError(validationErr)
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

func (u *users) CheckUserToken(ctx context.Context, token string) (string, bool) {
	userID := appcontext.GetUserID(ctx)

	userTokenParam := models.UserTokenParams{
		UserID:      userID,
		AccessToken: token,
		IsRevoked:   &[]bool{false}[0],
	}

	userTokenRes, err := u.userToken.CheckTokenExist(ctx, userTokenParam)
	if err != nil || userTokenRes == 0 {
		return "Token is unauthorized", false
	}

	return "", true
}

func (u *users) GetUserProfile(ctx context.Context) (models.Users, error) {
	var res models.Users

	userId := appcontext.GetUserID(ctx)
	userParam := models.UserParams{
		ID: userId,
	}

	userRes, err := u.user.Get(ctx, userParam)
	if err != nil {
		return res, err
	}

	return userRes, nil
}

func (u *users) UpdateUser(ctx context.Context, body models.UpdateUserParams, params models.UserParams) (models.Users, error) {
	var res models.Users

	userId := appcontext.GetUserID(ctx)
	if userId != params.ID {
		return res, errors.Forbidden("You are not allowed to update this user")
	}

	if err := u.validator.ValidateStruct(body); err != nil {
		validationErr, _ := u.validator.GetValidationErrors(err)
		return res, errors.ValidationError(validationErr)
	}

	userParam := models.UserParams{
		ID: userId,
	}

	userField := models.Users{
		Username: body.Username,
		Email:    body.Email,
	}

	userRes, err := u.user.Update(ctx, userField, userParam)
	if err != nil {
		return res, err
	}

	return userRes, nil
}

func (u *users) DeactivateUser(ctx context.Context, params models.UserParams) (models.Users, error) {
	var res models.Users

	userId := appcontext.GetUserID(ctx)
	if userId != params.ID {
		return res, errors.Forbidden("You are not allowed to deactivate this user")
	}

	userParam := models.UserParams{
		ID: userId,
	}

	userField := models.Users{
		IsActived: &[]bool{false}[0],
	}

	userRes, err := u.user.Update(ctx, userField, userParam)
	if err != nil {
		return res, err
	}

	userTokenParam := models.UserTokenParams{
		UserID: userId,
	}

	userTokenField := models.UserToken{
		IsRevoked: &[]bool{true}[0],
	}

	_, err = u.userToken.Update(ctx, userTokenField, userTokenParam)
	if err != nil {
		return res, err
	}

	return userRes, nil
}
