package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/utils/db"
)

type Order interface {
	db.TrxSupportRepo

	// GetListByCustomerID retreives list of order by customer id
	GetListByCustomerID(ctx context.Context, customerID int) ([]*entity.Order, error)

	// GetItemsListByOrderID retreives list of order items by order id
	GetItemsListByOrderID(ctx context.Context, orderID int) ([]*entity.OrderItem, error)

	// GetListUniqueItems retreive list of unique plant id by customer id
	GetListUniqueItems(ctx context.Context, customerId int) ([]*entity.OrderItem, error)

	// Create insert a new order record.
	Create(ctx context.Context, to db.TrxObj, order *entity.Order) error

	// CreateItems insert a list of new order item records.
	CreateItems(ctx context.Context, to db.TrxObj, orderItems []*entity.OrderItem) error
}

type Repository struct {
	db.GormTrxSupport
}

func NewRepository(gormDB *gorm.DB) *Repository {
	return &Repository{db.GormTrxSupport{DB: gormDB}}
}

func (r *Repository) GetListByCustomerID(ctx context.Context, customerID int) (res []*entity.Order, err error) {
	err = r.DB.Where("customer_id = ?", customerID).Find(&res).Error
	return
}

func (r *Repository) GetItemsListByOrderID(ctx context.Context, orderID int) (res []*entity.OrderItem, err error) {
	err = r.DB.Where("order_id = ?", orderID).Find(&res).Error
	return
}

func (r *Repository) Create(ctx context.Context, to db.TrxObj, order *entity.Order) error {
	return r.Trx(to).WithContext(ctx).Create(order).Error
}
func (r *Repository) CreateItems(ctx context.Context, to db.TrxObj, orderItems []*entity.OrderItem) error {
	return r.Trx(to).WithContext(ctx).Create(orderItems).Error
}

func (r *Repository) GetListUniqueItems(ctx context.Context, customerId int) (res []*entity.OrderItem, err error) {
	err = r.DB.WithContext(ctx).Raw("select distinct oi.plant_id from orders s join order_items oi on s.id = oi.order_id where s.customer_id = ?", customerId).Scan(&res).Error
	return
}
