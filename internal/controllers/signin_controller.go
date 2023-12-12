package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/gin-gonic/gin"
)

type SigninController struct {
	UserService entity.SignupService
	Env         *config.Config
}

// @Summary	Authenticate user
// @Schemes
// @Description	authenticates a user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		entity.SigninRequest	true	"User"
// @Success		200		{object}	entity.SigninResponse
// @Router			/auth/signin [post]
func (sc *SigninController) Signin(ctx *gin.Context) {
	var payload entity.SigninRequest

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	user, err := sc.UserService.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = sc.UserService.CompareHashAndPassword(user.Password, payload.Password)
	if err != nil {
		ctx.Error(entity.Unauthorized)
		return
	}

	token, err := sc.UserService.CreateAccessToken(user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	signinResponse := entity.SigninResponse{
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, signinResponse)
}
