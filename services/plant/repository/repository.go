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

	// GetListByIDs retreive list of plants by ids
	GetListByIDs(ctx context.Context, ids []int) ([]*entity.Plant, error)

	// UpdateStockByID update plant stock by id
	UpdateStockByID(ctx context.Context, id int, stock int) error
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

func (r *Repository) GetListByIDs(ctx context.Context, ids []int) (res []*entity.Plant, err error) {
	err = r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	return
}

func (r *Repository) UpdateStockByID(ctx context.Context, id int, stock int) error {
	return r.DB.WithContext(ctx).Model(&entity.Plant{ID: id}).Update("stock", stock).Error
}
