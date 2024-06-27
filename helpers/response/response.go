package response

import (
	"context"
	goerr "errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/helpers/log"
)

type response struct {
	log log.LogInterface
}

type Interface interface {
	Success(c *gin.Context, message string, data interface{}, pg *PaginationParam)
	Created(c *gin.Context, message string, data interface{})
	Error(c *gin.Context, err error)
}

func Init(log log.LogInterface) Interface {
	return &response{log: log}
}

func (r *response) Success(c *gin.Context, message string, data interface{}, pg *PaginationParam) {
	c.JSON(200, HTTPResponse{
		Meta:       getRequestMetadata(c),
		Message:    ResponseMessage{Title: "Sukses", Description: message},
		IsSuccess:  true,
		Data:       data,
		Pagination: pg,
	})
	r.log.Info(c.Request.Context(), message, nil)
}

func (r *response) Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, HTTPResponse{
		Meta: getRequestMetadata(c),
		Message: ResponseMessage{
			Title:       "Sukses",
			Description: message,
		},
		IsSuccess: true,
		Data:      data,
	})
	r.log.Info(c.Request.Context(), message, data)
}

func (r *response) Error(c *gin.Context, err error) {
	if goerr.Is(err, context.DeadlineExceeded) {
		err = errors.RequestTimeout("Request timeout")
	}

	c.JSON(int(errors.GetCode(err)), HTTPResponse{
		Meta: getRequestMetadata(c),
		Message: ResponseMessage{
			Title:       errors.GetType(err),
			Description: errors.GetMessage(err),
		},
		IsSuccess: false,
		Data:      nil,
	})
	r.log.Error(c.Request.Context(), err.Error())
}

func getRequestMetadata(c *gin.Context) Meta {
	meta := Meta{
		RequestID: appcontext.GetRequestId(c.Request.Context()),
		Time:      time.Now().Format(time.RFC3339),
	}

	requestStartTime := appcontext.GetRequestStartTime(c.Request.Context())
	if !requestStartTime.IsZero() {
		elapsedTimeMs := time.Since(requestStartTime).Milliseconds()
		meta.TimeElapsed = fmt.Sprintf("%dms", elapsedTimeMs)
	}

	return meta
}
