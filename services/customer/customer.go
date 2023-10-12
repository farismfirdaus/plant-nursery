package customer

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/farismfirdaus/plant-nursery/auth"
	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/customer/repository"
	"github.com/farismfirdaus/plant-nursery/utils/helper"
)

type Customer interface {
	// Register registers customer profile into database.
	Register(ctx context.Context, customer *entity.Customer) error

	// NewSession validates customer credential and create new session token
	NewSession(ctx context.Context, email string, password string) (string, error)
}

type Handler struct {
	repo repository.Customer
	auth auth.Auth
}

func NewHandler(
	repo repository.Customer,
	auth auth.Auth,
) *Handler {
	return &Handler{
		repo: repo,
		auth: auth,
	}
}

func (h *Handler) Register(ctx context.Context, customer *entity.Customer) error {
	if err := customer.Validate(); err != nil {
		return err
	}

	c, err := h.repo.GetByEmail(ctx, customer.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if c.ID > 0 {
		return apperr.InvalidEmailAlreadyTaken
	}

	// hash password before storing into database.
	hashedPassword, err := helper.HashPassword(customer.Password)
	if err != nil {
		return err
	}
	customer.Password = hashedPassword

	return h.repo.Create(ctx, customer)
}

func (h *Handler) NewSession(ctx context.Context, email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", apperr.BadRequest
	}

	customer, err := h.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperr.InvalidEmailNotFound
		}
		return "", err
	}

	if err := helper.VerifyPassword(customer.Password, password); err != nil {
		return "", apperr.InvalidPasswordNotMatch
	}

	token, err := h.auth.Sign(ctx, customer)
	if err != nil {
		return "", err
	}

	return token, err
}
