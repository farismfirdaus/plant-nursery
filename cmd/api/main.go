package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/farismfirdaus/plant-nursery/auth"
	"github.com/farismfirdaus/plant-nursery/config"
	"github.com/farismfirdaus/plant-nursery/utils/db"
	"github.com/farismfirdaus/plant-nursery/utils/middleware"
)

func main() {
	// initialize env variables config
	cfg := config.InitConfig()

	// connect to db
	pg := db.InitPostgres(cfg.Database)

	client := setupCert()
	repo := setupRepository(pg)
	service := setupService(repo, client)

	// init router
	r := gin.New()
	middleware.New(r)

	v1 := r.Group("api/v1")

	setupController(v1, service)

	log.Info().Msgf("running on port %d", cfg.Port)

	// run router
	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}

func setupCert() *auth.CustomerAuth {
	priv, err := os.ReadFile("cert/id_rsa")
	if err != nil {
		panic(err)
	}

	pub, err := os.ReadFile("cert/id_rsa.pub")
	if err != nil {
		panic(err)
	}

	client, err := auth.NewCustomerAuth(priv, pub)
	if err != nil {
		panic(err)
	}

	return client
}
