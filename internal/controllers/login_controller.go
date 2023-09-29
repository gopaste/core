package controllers

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginService domain.LoginService
	Env          *config.Config
}

func (lc *LoginController) Login(ctx *gin.Context) {

}
