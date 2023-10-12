package db

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/farismfirdaus/plant-nursery/config"
)

func InitPostgres(cfg config.DatabaseConfig) *gorm.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Database)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		CreateBatchSize: 500,
		Logger: logger.New(NewGormLogger(), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
		}),
	})
	log.Info().Msgf("connecting to %s", connStr)
	if err != nil {
		panic("Can't connect to database!")
	}

	return db
}
