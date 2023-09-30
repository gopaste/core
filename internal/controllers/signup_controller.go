package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/core/error"
	"github.com/gin-gonic/gin"
)

var (
	BadRequest        = error.NewHttpError("Bad request occurred", "Missing required parameters", http.StatusBadRequest)
	NotFound          = error.NewHttpError("Resource not found", "The requested resource does not exist", http.StatusNotFound)
	ServerError       = error.NewHttpError("Internal server error", "An unexpected error occurred on the server", http.StatusInternalServerError)
	UserConflictError = error.NewHttpError("User conflict", "A user with the same email already exists", http.StatusConflict)
)

type SignupController struct {
	UserService domain.SignupService
	Env         *config.Config
}

func (lc *SignupController) Signup(c *gin.Context) {
	var payload domain.User

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.Error(BadRequest)
		return
	}

	exist, err := lc.UserService.UserExistsByEmail(c, payload.Email)
	if err != nil {
		c.Error(ServerError)
		return
	}

	if exist {
		c.Error(UserConflictError)
		return
	}

	err = lc.UserService.Create(c, &payload)
	if err != nil {
		c.Error(ServerError)
		return
	}

	c.JSON(http.StatusCreated, "criado com sucesso")
}
