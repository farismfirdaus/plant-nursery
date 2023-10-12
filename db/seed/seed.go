package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/farismfirdaus/plant-nursery/config"
	"github.com/farismfirdaus/plant-nursery/entity"
	"github.com/farismfirdaus/plant-nursery/utils/db"
)

func main() {
	ctx := context.Background()

	cfg := config.InitConfig()

	db := db.InitPostgres(cfg.Database)

	repo := setupRepository(db)

	seeds := []func(ctx context.Context, repo *repository) error{
		addPlant,
		addCustomer,
		addCart,
		addOrder,
	}
	for _, seed := range seeds {
		if err := seed(ctx, repo); err != nil {
			panic(err)
		}
	}

	log.Info().Msg("finish seeding data")
}

func addPlant(ctx context.Context, repo *repository) error {
	plantData, err := os.ReadFile("./db/seed/data/plants.json")
	if err != nil {
		return err
	}

	var plants []*entity.Plant
	if err := json.Unmarshal(plantData, &plants); err != nil {
		return err
	}

	err = repo.plantRepo.Creates(ctx, plants)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `plants` table")
	}

	return nil
}

func addCustomer(ctx context.Context, repo *repository) error {
	customerData, err := os.ReadFile("./db/seed/data/customer.json")
	if err != nil {
		return err
	}

	var customer entity.Customer
	if err := json.Unmarshal(customerData, &customer); err != nil {
		return err
	}

	err = repo.customerRepo.Create(ctx, &customer)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `customers` table")
	}

	return nil
}

func addCart(ctx context.Context, repo *repository) error {
	cartData, err := os.ReadFile("./db/seed/data/cart.json")
	if err != nil {
		return err
	}

	var cart entity.Cart
	if err := json.Unmarshal(cartData, &cart); err != nil {
		return err
	}

	err = repo.cartRepo.Create(ctx, nil, &cart)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `carts` table")
	}

	cartItemData, err := os.ReadFile("./db/seed/data/cart_item.json")
	if err != nil {
		return err
	}

	var cartItem entity.CartItem
	if err := json.Unmarshal(cartItemData, &cartItem); err != nil {
		return err
	}

	err = repo.cartRepo.UpsertItem(ctx, nil, &cartItem)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `carts` table")
	}

	return nil
}

func addOrder(ctx context.Context, repo *repository) error {
	orderData, err := os.ReadFile("./db/seed/data/order.json")
	if err != nil {
		return err
	}

	var order entity.Order
	if err := json.Unmarshal(orderData, &order); err != nil {
		return err
	}

	err = repo.orderRepo.Create(ctx, nil, &order)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `orders` table")
	}

	orderItemData, err := os.ReadFile("./db/seed/data/order_item.json")
	if err != nil {
		return err
	}

	var orderItem []*entity.OrderItem
	if err := json.Unmarshal(orderItemData, &orderItem); err != nil {
		return err
	}

	err = repo.orderRepo.CreateItems(ctx, nil, orderItem)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seed `order_items` table")
	}

	return nil
}
