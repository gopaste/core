package main

import (
	"fmt"
	"log"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/app"
	"github.com/Caixetadev/snippet/internal/infra/db/postgres"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
)

func main() {
	cfg, err := config.NewConfig(".")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	validation := validation.NewValidator(validatorv10.New())

	db, err := postgres.New(cfg.DBURL)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	defer db.Close()

	router := gin.Default()

	app.Run(cfg, db, router, validation)

	router.Run()
}
