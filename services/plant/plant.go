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

	// GetListByIDs retreive list of plants by ids.
	GetListByIDs(ctx context.Context, ids []int) ([]*entity.Plant, error)

	// UpdateStockByID update plant stock by id
	UpdateStockByID(ctx context.Context, id int, stock int) error
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

func (h *Handler) GetListByIDs(ctx context.Context, ids []int) ([]*entity.Plant, error) {
	return h.repo.GetListByIDs(ctx, ids)
}

func (h *Handler) UpdateStockByID(ctx context.Context, id int, stock int) error {
	return h.repo.UpdateStockByID(ctx, id, stock)
}
