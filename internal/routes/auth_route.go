package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	"github.com/Caixetadev/snippet/internal/core/services"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAuthRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator) {
	ur := repository.NewUserRepository(db)

	lc := &controllers.SigninController{
		UserService: services.NewUserService(ur, validation),
		Env:         cfg,
	}

	sc := &controllers.SignupController{
		UserService: services.NewUserService(ur, validation),
		Env:         cfg,
	}

	group.POST("/signup", sc.Signup)
	group.POST("/signin", lc.Signin)
}
