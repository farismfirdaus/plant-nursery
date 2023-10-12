package order

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"

	// mocks
	cart_mock "github.com/farismfirdaus/plant-nursery/services/cart/mock"
	repo_mock "github.com/farismfirdaus/plant-nursery/services/order/repository/mock"
	plant_mock "github.com/farismfirdaus/plant-nursery/services/plant/mock"
	db_mock "github.com/farismfirdaus/plant-nursery/utils/db/mock"
)

func Test_GetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Order)
			in       int
		}{
			{
				name: "success",
				mockFunc: func(o *repo_mock.Order) {
					o.On("GetListByCustomerID", mock.Anything, 1).Return([]*entity.Order{{ID: 1}}, nil)
					o.On("GetItemsListByOrderID", mock.Anything, 1).Return([]*entity.OrderItem{}, nil)
				},
				in: 1,
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewOrder(t)

				test.mockFunc(repoMock)

				orders, err := NewHandler(repoMock, nil, nil).GetList(ctx, test.in)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
				if orders == nil {
					t.Fatalf("[%s] response should be not nil, but got nil", test.name)
				}
			})
		}
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Order)
			in       int
			wantErr  error
		}{
			{
				name:     "invalid customer id",
				mockFunc: func(o *repo_mock.Order) {},
				wantErr:  apperr.InvalidCustomerID,
			},
			{
				name: "error get list orders",
				mockFunc: func(o *repo_mock.Order) {
					o.On("GetListByCustomerID", mock.Anything, 1).Return(nil, errors.New("failed"))
				},
				in: 1,
			},
			{
				name: "success",
				mockFunc: func(o *repo_mock.Order) {
					o.On("GetListByCustomerID", mock.Anything, 1).Return([]*entity.Order{{ID: 1}}, nil)
					o.On("GetItemsListByOrderID", mock.Anything, 1).Return(nil, errors.New("failed"))
				},
				in: 1,
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewOrder(t)

				test.mockFunc(repoMock)

				_, err := NewHandler(repoMock, nil, nil).GetList(ctx, test.in)
				if err == nil {
					t.Fatalf("[%s] error should be not nil, but got nil", test.name)
				}
				if test.wantErr != nil && !errors.Is(err, test.wantErr) {
					t.Fatalf("[%s] error should be %s, but got %s", test.name, test.wantErr, err)
				}
			})
		}
	})
}

func Test_Checkout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name         string
			mockFunc     func(*repo_mock.Order, *plant_mock.Plant, *cart_mock.Cart)
			inCustomerID int
			inCartID     int
		}{
			{
				name: "success",
				mockFunc: func(o *repo_mock.Order, p *plant_mock.Plant, c *cart_mock.Cart) {
					c.On("GetByIDandCustomerID", mock.Anything, 2, 1).Return(&entity.Cart{
						ID:     2,
						Status: entity.CartStatusOpen,
						CartItems: []*entity.CartItem{
							{PlantID: 10, CartID: 2, Quantity: 10},
						},
					}, nil)
					p.On("GetListByIDs", mock.Anything, []int{10}).Return([]*entity.Plant{
						{ID: 10, Stock: 10, Price: decimal.NewFromInt(15)},
					}, nil)

					trxObj := db_mock.NewTrxObj(t)
					o.On("Begin").Return(trxObj, nil)
					p.On("UpdateStockByID", mock.Anything, 10, 0).Return(nil)
					o.On("Create", mock.Anything, mock.Anything, mock.IsType(&entity.Order{})).Return(nil)
					o.On("CreateItems", mock.Anything, mock.Anything, mock.IsType([]*entity.OrderItem{})).Return(nil)
					c.On("CloseCartByID", mock.Anything, 2).Return(nil)
					trxObj.On("Commit").Return(nil)
				},
				inCustomerID: 1,
				inCartID:     2,
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewOrder(t)
				plantMock := plant_mock.NewPlant(t)
				cartMock := cart_mock.NewCart(t)

				test.mockFunc(repoMock, plantMock, cartMock)

				err := NewHandler(repoMock, plantMock, cartMock).Checkout(ctx, test.inCustomerID, test.inCartID)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
			})
		}
	})
}
