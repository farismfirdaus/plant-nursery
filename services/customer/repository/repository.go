package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/utils/db"
)

type Customer interface {
	db.TrxSupportRepo

	// GetByEmail retreive customer data by email
	GetByEmail(ctx context.Context, email string) (*entity.Customer, error)

	// Creates insert new customer record.
	// returning inserted ids on success.
	Create(ctx context.Context, customer *entity.Customer) error
}

type Repository struct {
	db.GormTrxSupport
}

func NewRepository(gormDB *gorm.DB) *Repository {
	return &Repository{db.GormTrxSupport{DB: gormDB}}
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (res *entity.Customer, err error) {
	err = r.DB.WithContext(ctx).Where("email = ?", email).Take(&res).Error
	return
}

func (r *Repository) Create(ctx context.Context, customer *entity.Customer) error {
	return r.DB.WithContext(ctx).Create(customer).Error
}
