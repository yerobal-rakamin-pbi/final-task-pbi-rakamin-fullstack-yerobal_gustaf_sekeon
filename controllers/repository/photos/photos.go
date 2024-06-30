package photos

import (
	"context"

	"rakamin-final-task/database"
	"rakamin-final-task/models"
)

type Interface interface {
}

type photos struct {
	db *database.DB
}

func Init(db *database.DB) Interface {
	return &photos{
		db: db,
	}
}

func (p *photos) Create(ctx context.Context, photo models.Photos) (models.Photos, error) {
	if err := p.db.ORM.WithContext(ctx).Create(&photo).Error; err != nil {
		return photo, err
	}

	return photo, nil
}
