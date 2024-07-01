package photos

import (
	"context"
	"fmt"
	"time"

	photoRepo "rakamin-final-task/controllers/repository/photos"
	userRepo "rakamin-final-task/controllers/repository/users"
	"rakamin-final-task/helpers/appcontext"
	"rakamin-final-task/helpers/files"
	"rakamin-final-task/helpers/storage"
	"rakamin-final-task/models"
)

type Interface interface {
}

type photos struct {
	photo   photoRepo.Interface
	storage storage.Interface
	user    userRepo.Interface
}

type InitParam struct {
	PhotoRepo photoRepo.Interface
	UserRepo  userRepo.Interface
	Storage   storage.Interface
}

func Init(param InitParam) Interface {
	return &photos{
		photo:   param.PhotoRepo,
		storage: param.Storage,
		user:    param.UserRepo,
	}
}

func (p *photos) Create(ctx context.Context, param models.CreatePhotoParams, photoFile *files.File) (models.Photos, error) {
	var photo models.Photos

	userID := appcontext.GetUserID(ctx)

	// format: {userID}_{timestamp}
	fileName := fmt.Sprintf("%d_%d", userID, time.Now().Unix())
	photoURL, err := p.storage.Upload(ctx, photoFile, fileName)
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

func (p *photos) GetList(ctx context.Context, param models.PhotoParams) ([]models.Photos, error) {
	userID := appcontext.GetUserID(ctx)

	photoParam := models.PhotoParams{
		UserID: userID,
	}

	photos, _, err := p.photo.GetList(ctx, photoParam)
	if err != nil {
		return photos, err
	}

	return photos, nil
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

	err := p.photo.Delete(ctx, photoParam)
	if err != nil {
		return err
	}

	return nil
}
