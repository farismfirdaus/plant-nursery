package plant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/farismfirdaus/plant-nursery/entity"

	// mocks
	repo_mock "github.com/farismfirdaus/plant-nursery/services/plant/repository/mock"
)

func Test_GetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Plant)
		}{
			{
				name: "success",
				mockFunc: func(p *repo_mock.Plant) {
					p.On("GetList", mock.Anything).Return([]*entity.Plant{}, nil)
				},
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewPlant(t)

				test.mockFunc(repoMock)

				plants, err := NewHandler(repoMock).GetList(ctx)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
				if plants == nil {
					t.Fatalf("[%s] response should be not nil, but got nil", test.name)
				}
			})
		}
	})
}

func Test_GetListByIDs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Plant)
			in       []int
		}{
			{
				name: "success",
				mockFunc: func(p *repo_mock.Plant) {
					p.On("GetListByIDs", mock.Anything, []int{1, 2, 3}).Return([]*entity.Plant{}, nil)
				},
				in: []int{1, 2, 3},
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewPlant(t)

				test.mockFunc(repoMock)

				plants, err := NewHandler(repoMock).GetListByIDs(ctx, test.in)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
				if plants == nil {
					t.Fatalf("[%s] response should be not nil, but got nil", test.name)
				}
			})
		}
	})
}

func Test_UpdateStockByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Plant)
			inId     int
			inStock  int
		}{
			{
				name: "success",
				mockFunc: func(p *repo_mock.Plant) {
					p.On("UpdateStockByID", mock.Anything, 1, 1).Return(nil)
				},
				inId:    1,
				inStock: 1,
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewPlant(t)

				test.mockFunc(repoMock)

				err := NewHandler(repoMock).UpdateStockByID(ctx, test.inId, test.inStock)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
			})
		}
	})
}
