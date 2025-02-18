package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rakamin-final-task/config"
	uc "rakamin-final-task/controllers/usecase"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/errors"
	jwtLib "rakamin-final-task/helpers/jwt"
	"rakamin-final-task/helpers/response"
)

const (
	HeaderRequestId = "x-request-id"
)

type Interface interface {
	SetTimeout(c *gin.Context)
	AddFieldsToCtx(c *gin.Context)
	SetCors() gin.HandlerFunc
	CheckJWT() gin.HandlerFunc
}

type middleware struct {
	config   config.Application
	http     *gin.Engine
	jwt      jwtLib.Interface
	response response.Interface
	usecase  uc.Usecase
}

type InitParam struct {
	Config   config.Application
	Http     *gin.Engine
	Response response.Interface
	Usecase  uc.Usecase
}

func Init(param InitParam) Interface {
	jwt := jwtLib.Init(param.Config.Server.JWT.ExpSec, param.Config.Server.JWT.Secret)

	return &middleware{
		http:     param.Http,
		config:   param.Config,
		jwt:      jwt,
		response: param.Response,
		usecase:  param.Usecase,
	}
}

// Timeout middleware wraps the request context with a timeout.
func (m *middleware) SetTimeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(m.config.Server.RequestTimeoutSec) * time.Second)

	// Cancel to clean up resources
	defer cancel()

	// Set the new context and replace the request context
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func (m *middleware) AddFieldsToCtx(c *gin.Context) {
	requestID := uuid.New().String()

	ctx := c.Request.Context()
	ctx = appcontext.SetRequestId(ctx, requestID)
	ctx = appcontext.SetUserAgent(ctx, c.Request.Header.Get(appcontext.HeaderUserAgent))
	ctx = appcontext.SetDeviceType(ctx, c.Request.Header.Get(appcontext.HeaderDeviceType))
	c.Request = c.Request.WithContext(ctx)
	c.Header(HeaderRequestId, requestID)

	c.Next()
}

func (m *middleware) SetCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (m *middleware) CheckJWT() gin.HandlerFunc {
	return m.checkJWT
}

func (m *middleware) checkJWT(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		m.response.Error(c, errors.Unauthorized("Token tidak valid"))
		c.Abort()
		return
	}

	header = header[len("Bearer "):]
	tokenClaims, err := m.jwt.DecodeToken(header)
	if err != nil {
		m.response.Error(c, err)
		c.Abort()
		return
	}

	userData := tokenClaims["data"].(map[string]interface{})
	ctx := c.Request.Context()
	ctx = appcontext.SetUserID(ctx, int64(userData["id"].(float64)))

	msg, isValid := m.usecase.Users.CheckUserToken(ctx, header)
	if !isValid {
		m.response.Error(c, errors.Unauthorized(msg))
		c.Abort()
		return
	}
	
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
