package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/errors"
	//"rakamin-final-task/src/models"
)

func (r *router) BindParam(ctx *gin.Context, param interface{}) error {
	if err := ctx.ShouldBindUri(param); err != nil {
		return err
	}

	return ctx.ShouldBindWith(param, binding.Query)
}

func (r *router) BindBody(ctx *gin.Context, body interface{}) error {
	return ctx.ShouldBindWith(body, binding.Default(ctx.Request.Method, ctx.ContentType()))
}

func (r *router) SuccessResponse(ctx *gin.Context, message string, data interface{}, pg *model.PaginationParam) {
	ctx.JSON(200, model.HTTPResponse{
		Meta:       getRequestMetadata(ctx),
		Message:    model.ResponseMessage{Title: "Sukses", Description: message},
		IsSuccess:  true,
		Data:       data,
		Pagination: pg,
	})
	r.log.Info(ctx.Request.Context(), message, nil)
}

func (r *router) CreatedResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(201, model.HTTPResponse{
		Meta: getRequestMetadata(ctx),
		Message: model.ResponseMessage{
			Title:       "Sukses",
			Description: message,
		},
		IsSuccess: true,
		Data:      data,
	})
	r.log.Info(ctx.Request.Context(), message, data)
}

func (r *router) ErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(int(errors.GetCode(err)), model.HTTPResponse{
		Meta: getRequestMetadata(ctx),
		Message: model.ResponseMessage{
			Title:       errors.GetType(err),
			Description: errors.GetMessage(err),
		},
		IsSuccess: false,
		Data:      nil,
	})
	r.log.Error(ctx.Request.Context(), err.Error())
}

func getRequestMetadata(ctx *gin.Context) model.Meta {
	meta := model.Meta{
		RequestID: appcontext.GetRequestId(ctx.Request.Context()),
		Time:      time.Now().Format(time.RFC3339),
	}

	requestStartTime := appcontext.GetRequestStartTime(ctx.Request.Context())
	if !requestStartTime.IsZero() {
		elapsedTimeMs := time.Since(requestStartTime).Milliseconds()
		meta.TimeElapsed = fmt.Sprintf("%dms", elapsedTimeMs)
	}

	return meta
}
