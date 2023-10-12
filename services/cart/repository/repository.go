package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/utils/db"
)

type Cart interface {
	db.TrxSupportRepo

	// Create insert a new cart record.
	Create(ctx context.Context, to db.TrxObj, cart *entity.Cart) error

	// Update update cart record.
	Update(ctx context.Context, to db.TrxObj, cart *entity.Cart) error

	// UpdateStatus update cart record.
	UpdateStatus(ctx context.Context, to db.TrxObj, id int, status entity.CartStatus) error

	// GetActiveByCustomerID retreive active cart by customer id
	GetActiveByCustomerID(ctx context.Context, to db.TrxObj, customerID int) (*entity.Cart, error)

	// GetByIDAndCustomerID retreive cart by id and customer id
	GetByIDAndCustomerID(ctx context.Context, to db.TrxObj, id int, customerID int) (*entity.Cart, error)

	// GetListCartItemsByCartID retreive cart items by cart id
	GetListCartItemsByCartID(ctx context.Context, to db.TrxObj, cartID int) ([]*entity.CartItem, error)

	// UpsertItem upsert a cart item record.
	UpsertItem(ctx context.Context, to db.TrxObj, cartItem *entity.CartItem) error
}

type Repository struct {
	db.GormTrxSupport
}

func NewRepository(gormDB *gorm.DB) *Repository {
	return &Repository{db.GormTrxSupport{DB: gormDB}}
}

func (r *Repository) Create(ctx context.Context, to db.TrxObj, cart *entity.Cart) error {
	return r.Trx(to).WithContext(ctx).Create(cart).Error
}

func (r *Repository) Update(ctx context.Context, to db.TrxObj, cart *entity.Cart) error {
	return r.Trx(to).WithContext(ctx).Save(cart).Error
}

func (r *Repository) GetActiveByCustomerID(ctx context.Context, to db.TrxObj, customerID int) (cart *entity.Cart, err error) {
	err = r.Trx(to).WithContext(ctx).
		Where("customer_id = ?", customerID).
		Where("status = ?", entity.CartStatusOpen).
		Take(&cart).Error
	return
}

func (r *Repository) GetByIDAndCustomerID(ctx context.Context, to db.TrxObj, id int, customerID int) (cart *entity.Cart, err error) {
	err = r.Trx(to).WithContext(ctx).
		Where("id = ?", id).
		Where("customer_id = ?", customerID).
		Take(&cart).Error
	return
}

func (r *Repository) GetListCartItemsByCartID(ctx context.Context, to db.TrxObj, cartID int) (cartItems []*entity.CartItem, err error) {
	err = r.Trx(to).WithContext(ctx).
		Where("cart_id = ?", cartID).
		Find(&cartItems).Error
	return
}

func (r *Repository) UpdateStatus(ctx context.Context, to db.TrxObj, id int, status entity.CartStatus) error {
	return r.Trx(to).WithContext(ctx).Model(&entity.Cart{ID: id}).Update("status", status).Error
}

func (r *Repository) UpsertItem(ctx context.Context, to db.TrxObj, cartItem *entity.CartItem) error {
	return r.Trx(to).WithContext(ctx).Save(cartItem).Error
}
