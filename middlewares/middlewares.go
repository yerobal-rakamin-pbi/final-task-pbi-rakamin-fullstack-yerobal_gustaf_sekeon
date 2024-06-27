package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rakamin-final-task/config"
	"rakamin-final-task/helpers/appcontext"
)

const (
	HeaderRequestId    = "x-request-id"
)

type Interface interface {
	SetTimeout(c *gin.Context)
	AddFieldsToCtx(c *gin.Context)
	SetCors(c *gin.Context) gin.HandlerFunc
	CheckJWT(c *gin.Context) gin.HandlerFunc
}

type middleware struct {
	config config.Application
	http   *gin.Engine
}

func Init(http *gin.Engine) Interface {
	return &middleware{
		http: http,
	}
}

// Timeout middleware wraps the request context with a timeout.
func (m *middleware) SetTimeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(m.config.Server.RequestTimeoutSecond)*time.Second)

	// Cancel to clean up resources
	defer cancel()

	// Set the new context and replace the request context
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func (m *middleware) AddFieldsToContext(c *gin.Context) {
	requestID := uuid.New().String()

	ctx := c.Request.Context()
	ctx = appcontext.SetRequestId(ctx, requestID)
	ctx = appcontext.SetUserAgent(ctx, c.Request.Header.Get(appcontext.HeaderUserAgent))
	ctx = appcontext.SetDeviceType(ctx, c.Request.Header.Get(appcontext.HeaderDeviceType))
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func (m *middleware) SetCors(c *gin.Context) gin.HandlerFunc {
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
