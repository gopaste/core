package app

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/routes"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config, db *pgxpool.Pool, router *gin.Engine, validation validation.Validator) {
	publicRouter := router.Group("")

	routes.NewAuthRouter(cfg, db, publicRouter, validation)
}
