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

// @Summary	Create account
// @Schemes
// @Description	Create a new user account
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		domain.User	true	"User"
// @Success		200		{object}	domain.SignupResponse
// @Router			/auth/signup [post]
func (lc *SignupController) Signup(c *gin.Context) {
	var payload *domain.User

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

	signupResponse := domain.SignupResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, signupResponse)
}
