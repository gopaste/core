package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	"github.com/Caixetadev/snippet/internal/entity"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAuthRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator) {
	ur := repository.NewUserRepository(db)

	userService := services.NewUserService(ur, validation, &entity.BcryptPasswordHasher{})

	lc := &controllers.SigninController{
		UserService: userService,
		Env:         cfg,
	}

	sc := &controllers.SignupController{
		UserService: userService,
		Env:         cfg,
	}

	group.POST("/auth/signup", sc.Signup)
	group.POST("/auth/signin", lc.Signin)
}
