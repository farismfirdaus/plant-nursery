package customer

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"

	// mocks
	auth_mock "github.com/farismfirdaus/plant-nursery/auth/mock"
	repo_mock "github.com/farismfirdaus/plant-nursery/services/customer/repository/mock"
)

func Test_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Customer, *auth_mock.Auth)
			in       *entity.Customer
		}{
			{
				name: "success",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{}, gorm.ErrRecordNotFound)
					c.On("Create", mock.Anything, mock.IsType(&entity.Customer{})).Return(nil)
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: "123",
				},
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCustomer(t)
				authMock := auth_mock.NewAuth(t)

				test.mockFunc(repoMock, authMock)

				err := NewHandler(repoMock, authMock).Register(ctx, test.in)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
			})
		}
	})

	t.Run("failed", func(t *testing.T) {
		tests := []struct {
			name     string
			mockFunc func(*repo_mock.Customer, *auth_mock.Auth)
			in       *entity.Customer
			wantErr  error
		}{
			{
				name:     "request not valid",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {},
				in:       &entity.Customer{},
				wantErr:  apperr.InvalidCustomerFullName,
			},
			{
				name: "email already taken",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{ID: 1}, nil)
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: "123",
				},
				wantErr: apperr.InvalidEmailAlreadyTaken,
			},
			{
				name: "error get by email",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(nil, errors.New("failed"))
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: "123",
				},
				wantErr: nil, // error can't be check with `errors.Is` so we specify nil
			},
			{
				name: "error generating password",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{}, gorm.ErrRecordNotFound)
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: strings.Repeat("123", 100),
				},
				wantErr: nil, // error cannot be checked with errors.Is so we specify nil
			},
			{
				name: "error generating password",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{}, gorm.ErrRecordNotFound)
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: strings.Repeat("123", 100),
				},
				wantErr: bcrypt.ErrPasswordTooLong,
			},
			{
				name: "error creating record",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{}, gorm.ErrRecordNotFound)
					c.On("Create", mock.Anything, mock.IsType(&entity.Customer{})).Return(errors.New("failed"))
				},
				in: &entity.Customer{
					FullName: uuid.New().String(),
					Email:    "faris@gmail.com",
					Password: "123",
				},
				wantErr: nil, // error can't be check with `errors.Is` so we specify nil
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCustomer(t)
				authMock := auth_mock.NewAuth(t)

				test.mockFunc(repoMock, authMock)

				err := NewHandler(repoMock, authMock).Register(ctx, test.in)
				if err == nil {
					t.Fatalf("[%s] error should be not nil, but got nil", test.name)
				}
				if test.wantErr != nil && !errors.Is(test.wantErr, err) {
					t.Fatalf("[%s] error should be %s, but got %s", test.name, test.wantErr, err)
				}
			})
		}
	})
}

func Test_NewSession(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			mockFunc   func(*repo_mock.Customer, *auth_mock.Auth)
			inEmail    string
			inPassword string
			resp       string
		}{
			{
				name: "success",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{Password: "$2a$04$Pv7Sx3HWR.pzoIOw6YzFsOq50rJ63qdF1UX6PAroNZEBHr5It9dBO"}, nil)
					a.On("Sign", mock.Anything, mock.IsType(&entity.Customer{})).Return("token", nil)
				},
				inEmail:    "faris@gmail.com",
				inPassword: "123",
				resp:       "token",
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCustomer(t)
				authMock := auth_mock.NewAuth(t)

				test.mockFunc(repoMock, authMock)

				token, err := NewHandler(repoMock, authMock).NewSession(ctx, test.inEmail, test.inPassword)
				if err != nil {
					t.Fatalf("[%s] error should be nil, but got: %s", test.name, err)
				}
				if token != test.resp {
					t.Fatalf("[%s] response should be %s, but got: %s", test.name, test.resp, token)
				}
			})
		}
	})

	t.Run("failed", func(t *testing.T) {
		tests := []struct {
			name       string
			mockFunc   func(*repo_mock.Customer, *auth_mock.Auth)
			inEmail    string
			inPassword string
			wantErr    error
		}{
			{
				name:     "request not valid",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {},
				wantErr:  apperr.BadRequest,
			},
			{
				name: "email not found taken",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(nil, gorm.ErrRecordNotFound)
				},
				inEmail:    "faris@gmail.com",
				inPassword: "123",
				wantErr:    apperr.InvalidEmailNotFound,
			},
			{
				name: "error get email",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(nil, errors.New("failed"))
				},
				inEmail:    "faris@gmail.com",
				inPassword: "123",
				wantErr:    nil, // error can't be check with `errors.Is` so we specify nil
			},
			{
				name: "error verify password",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{}, nil)
				},
				inEmail:    "faris@gmail.com",
				inPassword: "123",
				wantErr:    apperr.InvalidPasswordNotMatch,
			},
			{
				name: "error signing token",
				mockFunc: func(c *repo_mock.Customer, a *auth_mock.Auth) {
					c.On("GetByEmail", mock.Anything, "faris@gmail.com").Return(&entity.Customer{Password: "$2a$04$Pv7Sx3HWR.pzoIOw6YzFsOq50rJ63qdF1UX6PAroNZEBHr5It9dBO"}, nil)
					a.On("Sign", mock.Anything, mock.IsType(&entity.Customer{})).Return("", errors.New("failed"))
				},
				inEmail:    "faris@gmail.com",
				inPassword: "123",
				wantErr:    nil, // error can't be check with `errors.Is` so we specify nil
			},
		}
		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()

				repoMock := repo_mock.NewCustomer(t)
				authMock := auth_mock.NewAuth(t)

				test.mockFunc(repoMock, authMock)

				_, err := NewHandler(repoMock, authMock).NewSession(ctx, test.inEmail, test.inPassword)
				if err == nil {
					t.Fatalf("[%s] error should be not nil, but got nil", test.name)
				}
				if test.wantErr != nil && !errors.Is(test.wantErr, err) {
					t.Fatalf("[%s] error should be %s, but got %s", test.name, test.wantErr, err)
				}
			})
		}
	})
}
