package order

import (
	"context"
	"errors"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/cart"
	"github.com/farismfirdaus/plant-nursery/services/order/repository"
	"github.com/farismfirdaus/plant-nursery/services/plant"
	"github.com/farismfirdaus/plant-nursery/utils/db"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order interface {
	// Checkout all of items in cart
	Checkout(ctx context.Context, customerID int, cartID int) error

	// GetList retreives all order history
	GetList(ctx context.Context, customerID int) ([]*entity.Order, error)
}

type Handler struct {
	repo repository.Order

	plantSvc plant.Plant
	cartSvc  cart.Cart
}

func NewHandler(
	repo repository.Order,
	plantSvc plant.Plant,
	cartSvc cart.Cart,
) *Handler {
	return &Handler{
		repo:     repo,
		plantSvc: plantSvc,
		cartSvc:  cartSvc,
	}
}
func (h *Handler) GetList(ctx context.Context, customerID int) ([]*entity.Order, error) {
	if customerID <= 0 {
		return nil, apperr.InvalidCustomerID
	}

	orders, err := h.repo.GetListByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		order.OrderItems, err = h.repo.GetItemsListByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (h *Handler) Checkout(ctx context.Context, customerID int, cartID int) error {
	cart, err := h.cartSvc.GetByIDandCustomerID(ctx, cartID, customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.InvalidCartNotFound
		}
		return err
	}

	if cart.Status == entity.CartStatusClosed {
		return apperr.InvalidCartNotFound
	}

	var (
		plantIds = []int{}

		plantsMap = map[int]*entity.Plant{}
	)
	for _, ci := range cart.CartItems {
		plantIds = append(plantIds, ci.PlantID)
	}

	plants, err := h.plantSvc.GetListByIDs(ctx, plantIds)
	if err != nil {
		return err
	}
	// creating hashmap for o(1) search by plant id
	for _, plant := range plants {
		plantsMap[plant.ID] = plant
	}

	var (
		order = &entity.Order{
			CustomerID: customerID,
		}

		orderItems = []*entity.OrderItem{}
	)

	return db.DBTransaction(h.repo, func(to db.TrxObj) error {
		for _, ci := range cart.CartItems {
			plant, exists := plantsMap[ci.PlantID]
			if !exists {
				return apperr.InvalidPlantID
			}

			if ci.Quantity > plant.Stock {
				return apperr.InvalidStockNotAvailable
			}

			stock := plant.Stock - ci.Quantity
			if err := h.plantSvc.UpdateStockByID(ctx, plant.ID, stock); err != nil {
				return err
			}

			totalAmount := plant.Price.Mul(decimal.NewFromInt(int64(ci.Quantity)))

			orderItems = append(orderItems, &entity.OrderItem{
				PlantID:     plant.ID,
				Price:       plant.Price,
				Quantity:    ci.Quantity,
				TotalAmount: totalAmount,
			})

			order.TotalAmount = order.TotalAmount.Add(totalAmount)
		}

		if err := h.repo.Create(ctx, to, order); err != nil {
			return err
		}

		for _, oi := range orderItems {
			oi.OrderID = order.ID
		}

		if err := h.repo.CreateItems(ctx, to, orderItems); err != nil {
			return err
		}

		if err := h.cartSvc.CloseCartByID(ctx, cart.ID); err != nil {
			return err
		}

		return nil
	})
}
