package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (r *router) BindParam(c *gin.Context, param interface{}) error {
	if err := c.ShouldBindUri(param); err != nil {
		return err
	}

	return c.ShouldBindWith(param, binding.Query)
}

func (r *router) BindBody(c *gin.Context, body interface{}) error {
	return c.ShouldBindWith(body, binding.Default(c.Request.Method, c.ContentType()))
}
