package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/middleware"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/passwordhash"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAuthRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator, tokenMaker token.Maker) {
	ur := repository.NewUserRepository(db)

	userService := services.NewUserService(ur, validation, &passwordhash.BcryptPasswordHasher{}, tokenMaker)
	emailService, err := services.NewSimpleEmailService(cfg)
	if err != nil {
		panic(err)
	}

	ac := &controllers.AuthController{
		UserService:  userService,
		EmailService: emailService,
		Env:          cfg,
	}

	group.POST("/auth/signup", ac.Signup)
	group.POST("/auth/signin", ac.Signin)
	group.POST("/auth/forgot-password", ac.ForgotPassword)
	group.POST("/auth/refresh-token", ac.RefreshToken)
	group.PUT("/auth/reset-password/:resetToken", ac.ResetPassword)
	group.GET("/auth/me", middleware.AuthPostMiddleware(tokenMaker), ac.Me)
}
