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

	if err := addPlant(ctx, repo); err != nil {
		panic(err)
	}
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
