package photos

import (
	"context"

	"rakamin-final-task/database"
	"rakamin-final-task/helpers/response"
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

func (p *photos) Get(ctx context.Context, params models.PhotoParams) (models.Photos, error) {
	var photo models.Photos

	res := p.db.ORM.WithContext(ctx).Where(params).First(&photo)
	if res.RowsAffected == 0 {
		return photo, nil
	} else if res.Error != nil {
		return photo, res.Error
	}

	return photo, nil
}

func (p *photos) GetList(ctx context.Context, params models.PhotoParams) ([]models.Photos, *response.PaginationParam, error) {
	var photos []models.Photos

	pg := response.PaginationParam{
		Limit: params.Limit,
		Page:  params.Page,
	}
	pg.SetDefaultPagination()

	res := p.db.ORM.WithContext(ctx).Where(params).Offset(int(pg.Offset)).Limit(int(pg.Limit)).Find(&photos)
	if res.RowsAffected == 0 {
		return photos, &pg, nil
	} else if res.Error != nil {
		return photos, &pg, res.Error
	}

	pg.ProcessPagination(res.RowsAffected)

	return photos, &pg, nil
}
