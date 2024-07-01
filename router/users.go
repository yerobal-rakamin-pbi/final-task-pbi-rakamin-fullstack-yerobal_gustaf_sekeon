package router

import (
	"github.com/gin-gonic/gin"
	"rakamin-final-task/models"
)

// @Summary Login
// @Description Login for user
// @Tags Users
// @Produce json
// @Param loginBody body models.UserLoginParams true "Login Body"
// @Success 200 {object} response.HTTPResponse{data=models.AuthResponse}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /users/login [POST]
func (r *router) Login(c *gin.Context) {
	var body models.UserLoginParams

	if err := r.BindBody(c, &body); err != nil {
		r.response.Error(c, err)
		return
	}

	userResponse, err := r.usecase.Users.Login(c.Request.Context(), body)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Login successfull", userResponse, nil)
}

// @Summary Register
// @Description Register for user
// @Tags Users
// @Produce json
// @Param registerBody body models.UserRegisterParams true "Register Body"
// @Success 201 {object} response.HTTPResponse{data=models.AuthResponse}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 409 {object} response.HTTPResponse{}
// @Failure 422 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /users/register [POST]
func (r *router) Register(c *gin.Context) {
	var body models.UserRegisterParams

	if err := r.BindBody(c, &body); err != nil {
		r.response.Error(c, err)
		return
	}

	userResponse, err := r.usecase.Users.Register(c.Request.Context(), body)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Created(c, "Register successfull", userResponse)
}

// @Summary Get User Profile
// @Description Get User Profile
// @Tags Users
// @Produce json
// @Success 200 {object} response.HTTPResponse{data=models.Users}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /users/profile [GET]
func (r *router) GetUserProfile(c *gin.Context) {
	userResponse, err := r.usecase.Users.GetUserProfile(c.Request.Context())
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Get user profile successfull", userResponse, nil)
}

// @Summary Update User
// @Description Update User
// @Tags Users
// @Produce json
// @Param user_id path int true "User ID"
// @Param updateBody body models.UpdateUserParams true "Update Body"
// @Success 200 {object} response.HTTPResponse{data=models.Users}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /users/{user_id} [PUT]
func (r *router) UpdateUser(c *gin.Context) {
	var body models.UpdateUserParams
	if err := r.BindBody(c, &body); err != nil {
		r.response.Error(c, err)
		return
	}

	var params models.UserParams
	if err := r.BindParam(c, &params); err != nil {
		r.response.Error(c, err)
		return
	}

	userResponse, err := r.usecase.Users.UpdateUser(c.Request.Context(), body, params)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Update user successfull", userResponse, nil)
}

// @Summary Deactivate User
// @Description Deactivate User
// @Tags Users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.HTTPResponse{data=models.Users}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /users/{user_id} [DELETE]
func (r *router) DeactivateUser(c *gin.Context) {
	var params models.UserParams
	if err := r.BindParam(c, &params); err != nil {
		r.response.Error(c, err)
		return
	}

	userResponse, err := r.usecase.Users.DeactivateUser(c.Request.Context(), params)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Deactivate user successfull", userResponse, nil)
}
