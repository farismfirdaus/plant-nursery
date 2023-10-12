package cart

import (
	"context"
	"errors"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/cart/repository"
	"github.com/farismfirdaus/plant-nursery/services/plant"
	"github.com/farismfirdaus/plant-nursery/utils/db"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cart interface {
	// Get retreive cart and its items
	Get(ctx context.Context, customerID int) (*entity.Cart, error)

	// Get retreive cart and its items by id and customer id
	GetByIDandCustomerID(ctx context.Context, id int, customerID int) (*entity.Cart, error)

	// AddItems
	// creating new cart item if plant id is not available in cart
	// updating quantity cart item if plant id is available in cart
	// TODO: validate if quantity < 0.
	//  	 this validation is ignored now, so user can subtract their cart item by providing negative number
	//		 negative quantity in cart items is possible
	AddItems(ctx context.Context, customerID int, items []*AddItemsRequest) []error

	// CloseCartByID closing cart by id
	CloseCartByID(ctx context.Context, id int) error

	// TODO: delete cart item
}

type Handler struct {
	repo     repository.Cart
	plantSvc plant.Plant
}

func NewHandler(
	repo repository.Cart,
	plantSvc plant.Plant,
) *Handler {
	return &Handler{
		repo:     repo,
		plantSvc: plantSvc,
	}
}

func (h *Handler) Get(ctx context.Context, customerID int) (cart *entity.Cart, err error) {
	if customerID <= 0 {
		return nil, apperr.InvalidCustomerID
	}

	// this function wrap db transaction
	err = db.DBTransaction(h.repo, func(to db.TrxObj) error {
		cart, err = h.repo.GetActiveByCustomerID(ctx, to, customerID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create new cart if customer doesnt have active cart
			cart = &entity.Cart{
				CustomerID:  customerID,
				Status:      entity.CartStatusOpen,
				TotalAmount: decimal.NewFromInt(0),
			}

			if err := h.repo.Create(ctx, to, cart); err != nil {
				return err
			}
		}

		cartItems, err := h.repo.GetListCartItemsByCartID(ctx, to, cart.ID)
		if err != nil {
			return err
		}

		cart.CartItems = cartItems

		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func (h *Handler) GetByIDandCustomerID(ctx context.Context, id int, customerID int) (*entity.Cart, error) {
	if customerID <= 0 {
		return nil, apperr.InvalidCustomerID
	}

	cart, err := h.repo.GetByIDAndCustomerID(ctx, nil, id, customerID)
	if err != nil {
		return nil, err
	}

	cartItems, err := h.repo.GetListCartItemsByCartID(ctx, nil, cart.ID)
	if err != nil {
		return nil, err
	}

	cart.CartItems = cartItems

	return cart, nil
}

func (h *Handler) AddItems(ctx context.Context, customerID int, items []*AddItemsRequest) (errs []error) {
	if customerID <= 0 {
		errs = append(errs, apperr.InvalidCustomerID)
		return
	}

	if len(items) <= 0 {
		errs = append(errs, apperr.BadRequest)
		return
	}

	// steps
	// 1. get cart
	// 2. get cart items
	// 3. iterate request cart item
	// 4. validate stock availibility
	// 		-> create new if plant id not exists in cart
	//		-> update quantity if plant id is exists in cart

	var (
		reqPlantIds = []int{}

		// creating hashmap for o(1) search by plant id
		plantsMap = map[int]*entity.Plant{}
	)

	for _, item := range items {
		reqPlantIds = append(reqPlantIds, item.PlantID)
	}

	plants, err := h.plantSvc.GetListByIDs(ctx, reqPlantIds)
	if err != nil {
		errs = append(errs, err)
		return
	}

	for _, plant := range plants {
		plantsMap[plant.ID] = plant
	}

	// this function wrap db transaction
	err = db.DBTransaction(h.repo, func(to db.TrxObj) error {
		cart, err := h.repo.GetActiveByCustomerID(ctx, to, customerID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create new cart if customer doesnt have active cart
			cart = &entity.Cart{
				CustomerID:  customerID,
				Status:      entity.CartStatusOpen,
				TotalAmount: decimal.NewFromInt(0),
			}

			if err := h.repo.Create(ctx, to, cart); err != nil {
				return err
			}
		}

		// creating hashmap for o(1) search by plant id
		var cartItemsMap = map[int]*entity.CartItem{}

		cartItems, err := h.repo.GetListCartItemsByCartID(ctx, to, cart.ID)
		if err != nil {
			return err
		}
		for _, ci := range cartItems {
			cartItemsMap[ci.PlantID] = ci
		}

		cart.TotalAmount = decimal.Zero

		var (
			toBeUpserted = []*entity.CartItem{}
		)

		for _, item := range items {
			amount := decimal.Decimal{}

			plant, exists := plantsMap[item.PlantID]
			if !exists {
				errs = append(errs, apperr.InvalidPlantID)
				continue
			}

			cartItem, exists := cartItemsMap[item.PlantID]
			if !exists {
				if item.Quantity > plant.Stock {
					errs = append(errs, apperr.InvalidStockNotAvailable)
					continue
				}

				cartItem = &entity.CartItem{
					CartID:   cart.ID,
					PlantID:  plant.ID,
					Quantity: item.Quantity,
				}

				amount = plant.Price.Mul(decimal.NewFromInt(int64(cartItem.Quantity)))

				toBeUpserted = append(toBeUpserted, cartItem)
			} else {
				qty := item.Quantity + cartItem.Quantity

				if qty > plant.Stock {
					errs = append(errs, apperr.InvalidStockNotAvailable)
					continue
				}

				cartItem.Quantity = qty

				amount = plant.Price.Mul(decimal.NewFromInt(int64(qty)))

				toBeUpserted = append(toBeUpserted, cartItem)
			}

			cart.TotalAmount = cart.TotalAmount.Add(amount)
		}

		if len(errs) > 0 {
			return errors.New("skip save changes")
		}

		for _, cartItem := range toBeUpserted {
			if err := h.repo.UpsertItem(ctx, to, cartItem); err != nil {
				return err
			}
		}

		if err := h.repo.Update(ctx, to, cart); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Err(err).Msg("transaction add cart items")
		return
	}

	return nil
}

func (h *Handler) CloseCartByID(ctx context.Context, id int) error {
	return h.repo.UpdateStatus(ctx, nil, id, entity.CartStatusClosed)
}
