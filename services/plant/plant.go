package plant

import (
	"context"

	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/services/plant/repository"
)

type Plant interface {
	// GetList retreive list of plants.
	// TODO: implement pagination
	GetList(ctx context.Context) ([]*entity.Plant, error)
}

type Handler struct {
	repo repository.Plant
}

func NewHandler(
	repo repository.Plant,
) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetList(ctx context.Context) ([]*entity.Plant, error) {
	return h.repo.GetList(ctx)
}
