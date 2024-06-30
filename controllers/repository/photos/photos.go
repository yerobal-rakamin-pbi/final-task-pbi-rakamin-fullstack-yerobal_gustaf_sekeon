package photos

import (
	"context"

	"rakamin-final-task/database"
	"rakamin-final-task/helpers/response"
	"rakamin-final-task/models"
)

type Interface interface {
	Create(ctx context.Context, photo models.Photos) (models.Photos, error)
	Get(ctx context.Context, params models.PhotoParams) (models.Photos, error)
	GetList(ctx context.Context, params models.PhotoParams) ([]models.Photos, *response.PaginationParam, error)
	Update(ctx context.Context, photo models.Photos, params models.PhotoParams) (models.Photos, error)
	Delete(ctx context.Context, params models.PhotoParams) error
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

func (p *photos) Update(ctx context.Context, photo models.Photos, params models.PhotoParams) (models.Photos, error) {
	if err := p.db.ORM.WithContext(ctx).Model(models.Photos{}).Where(params).Updates(&photo).Error; err != nil {
		return photo, err
	}

	return photo, nil
}

func (p *photos) Delete(ctx context.Context, params models.PhotoParams) error {
	res := p.db.ORM.WithContext(ctx).Where(params).Delete(&models.Photos{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
