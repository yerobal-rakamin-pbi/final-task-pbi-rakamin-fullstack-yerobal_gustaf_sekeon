package photos

import (
	"context"
	"fmt"
	"time"

	photoRepo "rakamin-final-task/controllers/repository/photos"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/files"
	"rakamin-final-task/helpers/response"
	"rakamin-final-task/helpers/storage"
	"rakamin-final-task/models"
)

type Interface interface {
	Create(ctx context.Context, param models.CreatePhotoParams, photoFile *files.File) (models.Photos, error)
	Get(ctx context.Context, param models.PhotoParams) (models.Photos, error)
	GetList(ctx context.Context, param models.PhotoParams) ([]models.Photos, *response.PaginationParam, error)
	Update(ctx context.Context, param models.PhotoParams, body models.UpdatePhotoParams) (models.Photos, error)
	Delete(ctx context.Context, param models.PhotoParams) error
}

const (
	photoPath = "photos"
)

type photos struct {
	photo   photoRepo.Interface
	storage storage.Interface
}

type InitParam struct {
	PhotoRepo photoRepo.Interface
	Storage   storage.Interface
}

func Init(param InitParam) Interface {
	return &photos{
		photo:   param.PhotoRepo,
		storage: param.Storage,
	}
}

func (p *photos) Create(ctx context.Context, param models.CreatePhotoParams, photoFile *files.File) (models.Photos, error) {
	var photo models.Photos

	userID := appcontext.GetUserID(ctx)

	// format: {userID}_{timestamp}
	photoFile.SetFileName(fmt.Sprintf("%d_%d", userID, time.Now().Unix()))
	photoURL, err := p.storage.Upload(ctx, photoFile, photoPath)
	if err != nil {
		return photo, err
	}

	photo = models.Photos{
		Title:    param.Title,
		Caption:  param.Caption,
		UserID:   userID,
		PhotoURL: photoURL,
	}

	photo, err = p.photo.Create(ctx, photo)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (p *photos) Get(ctx context.Context, param models.PhotoParams) (models.Photos, error) {
	userID := appcontext.GetUserID(ctx)

	photoParam := models.PhotoParams{
		ID:     param.ID,
		UserID: userID,
	}

	photo, err := p.photo.Get(ctx, photoParam)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (p *photos) GetList(ctx context.Context, param models.PhotoParams) ([]models.Photos, *response.PaginationParam, error) {
	userID := appcontext.GetUserID(ctx)

	photoParam := models.PhotoParams{
		UserID: userID,
	}

	photos, pg, err := p.photo.GetList(ctx, photoParam)
	if err != nil {
		return photos, pg, err
	}

	return photos, pg, nil
}

func (p *photos) Update(ctx context.Context, param models.PhotoParams, body models.UpdatePhotoParams) (models.Photos, error) {
	var photo models.Photos

	userID := appcontext.GetUserID(ctx)

	photoParam := models.PhotoParams{
		ID:     param.ID,
		UserID: userID,
	}

	photo.Title = body.Title
	photo.Caption = body.Caption

	photo, err := p.photo.Update(ctx, photo, photoParam)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (p *photos) Delete(ctx context.Context, param models.PhotoParams) error {
	userID := appcontext.GetUserID(ctx)

	photoParam := models.PhotoParams{
		ID:     param.ID,
		UserID: userID,
	}

	photo, err := p.photo.Get(ctx, photoParam)
	if err != nil {
		return err
	}

	p.storage.Delete(ctx, photo.PhotoURL, photoPath)

	if err := p.photo.Delete(ctx, photoParam); err != nil {
		return err
	}

	return nil
}
