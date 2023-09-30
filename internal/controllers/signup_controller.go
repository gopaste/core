package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/core/error"
	"github.com/gin-gonic/gin"
)

type SignupController struct {
	UserService domain.SignupService
	Env         *config.Config
}

func (lc *SignupController) Signup(c *gin.Context) {
	var payload domain.User

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.Error(error.BadRequest)
		return
	}

	exist, err := lc.UserService.UserExistsByEmail(c, payload.Email)
	if err != nil {
		c.Error(error.ServerError)
		return
	}

	if exist {
		c.Error(error.UserConflictError)
		return
	}

	err = lc.UserService.Create(c, &payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, "criado com sucesso")
}
