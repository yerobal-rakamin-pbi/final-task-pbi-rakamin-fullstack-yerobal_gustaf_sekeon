package router

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/helpers/files"
	"rakamin-final-task/models"
)

// @Summary Create Photo
// @Description Create photo
// @Tags Photos
// @Produce json
// @Param title formData string true "Title"
// @Param caption formData string true "Caption"
// @Param photo formData file true "Photo"
// @Accept multipart/form-data
// @Security BearerAuth
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

// @Summary Get Photo
// @Description Get photo
// @Tags Photos
// @Produce json
// @Param photo_id path int true "Photo ID"
// @Security BearerAuth
// @Success 200 {object} response.HTTPResponse{data=models.Photos}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /photos/{photo_id} [GET]
func (r *router) GetPhoto(c *gin.Context) {
	var photoParam models.PhotoParams
	if err := r.BindParam(c, &photoParam); err != nil {
		r.response.Error(c, err)
		return
	}

	photo, err := r.usecase.Photos.Get(c.Request.Context(), photoParam)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Get photo successfull", photo, nil)
}

// @Summary Get List Photo
// @Description Get list photo
// @Tags Photos
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Security BearerAuth
// @Success 200 {object} response.HTTPResponse{data=[]models.Photos,meta=response.PaginationParam}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /photos [GET]
func (r *router) GetListPhoto(c *gin.Context) {
	var photoParam models.PhotoParams
	if err := r.BindParam(c, &photoParam); err != nil {
		r.response.Error(c, err)
		return
	}

	photos, pg, err := r.usecase.Photos.GetList(c.Request.Context(), photoParam)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Get list photo successfull", photos, pg)
}

// @Summary Update Photo
// @Description Update photo
// @Tags Photos
// @Produce json
// @Param photo_id path int true "Photo ID"
// @Param photoBody body models.UpdatePhotoParams true "Update Body"
// @Security BearerAuth
// @Success 200 {object} response.HTTPResponse{data=models.Photos}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /photos/{photo_id} [PUT]
func (r *router) UpdatePhoto(c *gin.Context) {
	var body models.UpdatePhotoParams
	if err := r.BindBody(c, &body); err != nil {
		r.response.Error(c, err)
		return
	}

	var photoParam models.PhotoParams
	if err := r.BindParam(c, &photoParam); err != nil {
		r.response.Error(c, err)
		return
	}

	photo, err := r.usecase.Photos.Update(c.Request.Context(), photoParam, body)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Update photo successfull", photo, nil)
}

// @Summary Delete Photo
// @Description Delete photo
// @Tags Photos
// @Produce json
// @Param photo_id path int true "Photo ID"
// @Security BearerAuth
// @Success 200 {object} response.HTTPResponse{data=models.Photos}
// @Failure 400 {object} response.HTTPResponse{}
// @Failure 404 {object} response.HTTPResponse{}
// @Failure 500 {object} response.HTTPResponse{}
// @Router /photos/{photo_id} [DELETE]
func (r *router) DeletePhoto(c *gin.Context) {
	var photoParam models.PhotoParams
	if err := r.BindParam(c, &photoParam); err != nil {
		r.response.Error(c, err)
		return
	}

	err := r.usecase.Photos.Delete(c.Request.Context(), photoParam)
	if err != nil {
		r.response.Error(c, err)
		return
	}

	r.response.Success(c, "Delete photo successfull", nil, nil)
}
