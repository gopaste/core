package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/passwordhash"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator, tokenMaker token.Maker) {
	ur := repository.NewUserRepository(db)

	userService := services.NewUserService(ur, validation, &passwordhash.BcryptPasswordHasher{}, tokenMaker)

	uc := &controllers.UserController{
		UserService: userService,
		Env:         cfg,
	}

	group.GET("/user", uc.GetAuthenticatedUser)
}
