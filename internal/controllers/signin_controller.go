package controllers

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type SigninController struct {
	UserService domain.UserRepository
	Env         *config.Config
}

func (lc *SigninController) Signin(ctx *gin.Context) {
	lc.UserService.GetUserByEmail(ctx, "")
}
