package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/utils/db"
)

type Plant interface {
	db.TrxSupportRepo

	// Creates insert new plant records.
	// returning inserted ids on success.
	Creates(ctx context.Context, plants []*entity.Plant) error

	// GetList retreive list of plants
	GetList(ctx context.Context) ([]*entity.Plant, error)
}

type Repository struct {
	db.GormTrxSupport
}

func NewRepository(gormDB *gorm.DB) *Repository {
	return &Repository{db.GormTrxSupport{DB: gormDB}}
}

func (r *Repository) Creates(ctx context.Context, plants []*entity.Plant) error {
	return r.DB.WithContext(ctx).Create(plants).Error
}

func (r *Repository) GetList(ctx context.Context) (res []*entity.Plant, err error) {
	err = r.DB.WithContext(ctx).Find(&res).Error
	return
}
