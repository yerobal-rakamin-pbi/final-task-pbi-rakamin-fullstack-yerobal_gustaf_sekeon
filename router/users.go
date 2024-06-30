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
// @Router /v1/users/login [POST]
func (r *router) Login(c *gin.Context) {
	var body models.UserLoginParams

	if err := r.BindBody(c, &body); err != nil {
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
// @Success 200 {object} response.HTTPResponse{data=models.AuthResponse}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 409 {object} response.HTTPResponse{}
// @Failure 422 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /v1/users/register [POST]
func (r *router) Register(c *gin.Context) {
	var body models.UserRegisterParams

	if err := r.BindBody(c, &body); err != nil {
		return
	}

	userResponse, err := r.usecase.Users.Register(c.Request.Context(), body)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Register successfull", userResponse, nil)
}
