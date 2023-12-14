package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserService entity.UserService
	Env         *config.Config
}

// @Summary	Create account
// @Schemes
// @Description	Create a new user account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		entity.User	true	"User"
// @Success		200		{object}	entity.SignupResponse
// @Router			/auth/signup [post]
func (lc *AuthController) Signup(c *gin.Context) {
	var payload *entity.User

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.Error(entity.BadRequest)
		return
	}

	exist, err := lc.UserService.UserExistsByEmail(c, payload.Email)
	if err != nil {
		c.Error(entity.ServerError)
		return
	}

	if exist {
		c.Error(entity.UserConflictError)
		return
	}

	user, err := lc.UserService.Create(c, payload)
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := lc.UserService.CreateAccessToken(user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.Error(err)
		return
	}

	signupResponse := entity.SignupResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, signupResponse)
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
func (sc *AuthController) Signin(ctx *gin.Context) {
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
