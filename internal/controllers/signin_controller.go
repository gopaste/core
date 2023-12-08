package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/core/error"
	"github.com/gin-gonic/gin"
)

type SigninController struct {
	UserService domain.SignupService
	Env         *config.Config
}

func (sc *SigninController) Signin(ctx *gin.Context) {
	var payload domain.SigninRequest

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(error.BadRequest)
		return
	}

	user, err := sc.UserService.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = sc.UserService.CompareHashAndPassword(user.Password, payload.Password)
	if err != nil {
		ctx.Error(error.Unauthorized)
		return
	}

	token, err := sc.UserService.CreateAccessToken(user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		ctx.Error(error.BadRequest)
		return
	}

	signinResponse := domain.SigninResponse{
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, signinResponse)
}
