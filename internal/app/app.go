package app

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/middleware"
	"github.com/Caixetadev/snippet/internal/routes"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const BASE_PATH = "/api/v1"

func Run(cfg *config.Config, db *pgxpool.Pool, router *gin.Engine, validation validation.Validator, tokenMaker token.Maker) {
	publicRouter := router.Group(BASE_PATH)

	routes.NewAuthRouter(cfg, db, publicRouter, validation, tokenMaker)

	protectedRouter := router.Group(BASE_PATH)

	protectedRouter.Use(middleware.AuthPostMiddleware(tokenMaker))
	routes.NewPostRouter(cfg, db, protectedRouter, validation)
	routes.NewUserRouter(cfg, db, protectedRouter, validation, tokenMaker)
}
