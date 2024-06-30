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
