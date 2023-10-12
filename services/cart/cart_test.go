package cart

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"

	// mocks
	repo_mock "github.com/farismfirdaus/plant-nursery/services/cart/repository/mock"
	plant_mock "github.com/farismfirdaus/plant-nursery/services/plant/mock"
	db_mock "github.com/farismfirdaus/plant-nursery/utils/db/mock"
)

func Test_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			mockFunc       func(*repo_mock.Cart, *plant_mock.Plant)
			in             int
			cartItemsCount int
		}{
			{
				name: "success",
				mockFunc: func(c *repo_mock.Cart, p *plant_mock.Plant) {
					trxObj := db_mock.NewTrxObj(t)
					c.On("Begin").Return(trxObj, nil)
					c.On("GetActiveByCustomerID", mock.Anything, mock.Anything, 1).Return(&entity.Cart{ID: 1}, nil)
					c.On("GetListCartItemsByCartID", mock.Anything, mock.Anything, 1).Return([]*entity.CartItem{{ID: 1}}, nil)
					trxObj.On("Commit").Return(nil)
				},
				in:             1,
				cartItemsCount: 1,
			},
			{
				name: "success create new cart",
				mockFunc: func(c *repo_mock.Cart, p *plant_mock.Plant) {
					trxObj := db_mock.NewTrxObj(t)
					c.On("Begin").Return(trxObj, nil)
					c.On("GetActiveByCustomerID", mock.Anything, mock.Anything, 1).Return(nil, gorm.ErrRecordNotFound)
					c.On("Create", mock.Anything, mock.Anything, mock.IsType(&entity.Cart{})).Return(nil).Run(func(args mock.Arguments) {
						// inject id to cart
						cart := args.Get(2).(*entity.Cart)
						cart.ID = 1
					})
					c.On("GetListCartItemsByCartID", mock.Anything, mock.Anything, 1).Return([]*entity.CartItem{}, nil)
					trxObj.On("Commit").Return(nil)
				},
				in: 1,
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCart(t)
				plantMock := plant_mock.NewPlant(t)

				test.mockFunc(repoMock, plantMock)

				cart, err := NewHandler(repoMock, plantMock).Get(ctx, test.in)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
				if len(cart.CartItems) != test.cartItemsCount {
					t.Fatalf("[%s] cart items count should be %d, but got: %d", test.name, test.cartItemsCount, len(cart.CartItems))
				}
			})
		}
	})

	t.Run("failed", func(t *testing.T) {
		// TODO: add failed unit test
	})
}

func Test_AddItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name         string
			mockFunc     func(*repo_mock.Cart, *plant_mock.Plant)
			inCustomerID int
			in           []*AddItemsRequest
		}{
			{
				name: "success",
				mockFunc: func(c *repo_mock.Cart, p *plant_mock.Plant) {
					p.On("GetListByIDs", mock.Anything, []int{1, 2}).Return([]*entity.Plant{
						{ID: 1, Stock: 10, Price: decimal.NewFromInt(21)},
						{ID: 2, Stock: 5, Price: decimal.NewFromInt(21)},
					}, nil)
					trxObj := db_mock.NewTrxObj(t)
					c.On("Begin").Return(trxObj, nil)
					c.On("GetActiveByCustomerID", mock.Anything, mock.Anything, 1).Return(nil, gorm.ErrRecordNotFound) // err record not found, create new one
					c.On("Create", mock.Anything, mock.Anything, mock.IsType(&entity.Cart{})).Return(nil).Run(func(args mock.Arguments) {
						// inject id to cart
						cart := args.Get(2).(*entity.Cart)
						cart.ID = 1
					})
					c.On("GetListCartItemsByCartID", mock.Anything, mock.Anything, 1).Return([]*entity.CartItem{{ID: 1, CartID: 1, PlantID: 1, Quantity: 1}}, nil)

					// first upsert item
					c.On("UpsertItem", mock.Anything, mock.Anything, &entity.CartItem{
						ID:       1,
						CartID:   1,
						PlantID:  1,
						Quantity: 2,
					}).Return(nil).Once()

					// second upsert item
					c.On("UpsertItem", mock.Anything, mock.Anything, &entity.CartItem{
						ID:       0, // zero because not inserted to db yet
						CartID:   1,
						PlantID:  2,
						Quantity: 1,
					}).Return(nil).Once()

					c.On("Update", mock.Anything, mock.Anything, mock.IsType(&entity.Cart{})).Return(nil)
					trxObj.On("Commit").Return(nil)
				},
				inCustomerID: 1,
				in: []*AddItemsRequest{
					{PlantID: 1, Quantity: 1},
					{PlantID: 2, Quantity: 1},
				},
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCart(t)
				plantMock := plant_mock.NewPlant(t)

				test.mockFunc(repoMock, plantMock)

				err := NewHandler(repoMock, plantMock).AddItems(ctx, test.inCustomerID, test.in)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
			})
		}
	})

	t.Run("failed", func(t *testing.T) {
		// TODO: add failed unit test
	})
}
