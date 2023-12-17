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
	emailService, _ := services.NewSimpleEmailService()

	ac := &controllers.AuthController{
		UserService:  userService,
		EmailService: emailService,
		Env:          cfg,
	}

	group.POST("/auth/signup", ac.Signup)
	group.POST("/auth/signin", ac.Signin)
	group.POST("/auth/forgot-password", ac.ForgotPassword)
	group.PUT("/auth/resetpassword/:resetToken", ac.ResetPassword)
}
