package router

import (
	"mime/multipart"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/helpers/files"

	"rakamin-final-task/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create Photo
// @Description Create photo
// @Tags Photos
// @Produce json
// @Param title formData string true "Title"
// @Param caption formData string true "Caption"
// @Param photo formData file true "Photo"
// @Accept multipart/form-data
// @Success 201 {object} response.HTTPResponse{data=models.Photos}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /photos [POST]
func (r *router) CreatePhoto(c *gin.Context) {
	var body models.CreatePhotoParams

	if err := r.BindBody(c, &body); err != nil {
		r.response.Error(c, err)
		return
	}

	photoFile, meta, err := c.Request.FormFile("photo")
	if err != nil {
		r.response.Error(c, errors.BadRequest("File not found"))
		return
	}

	image, err := r.getPhotos(photoFile, meta)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	photo, err := r.usecase.Photos.Create(c.Request.Context(), body, image)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Created(c, "Photo created", photo)
}

func (r *router) getPhotos(file multipart.File, meta *multipart.FileHeader) (*files.File, error) {
	image := files.Init(file, meta)

	if !image.IsImage() {
		return nil, errors.BadRequest("File is not an image")
	}

	return image, nil
}
