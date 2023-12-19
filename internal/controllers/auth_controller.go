package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserService  entity.UserService
	EmailService entity.EmailService
	Env          *config.Config
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

// @Summary		Submit a request to reset the user's password
// @Description	Submit a request to reset the user's password by providing their email address.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		entity.ForgotPasswordRequest	true	"User's email"
// @Success		200		{object}	entity.Response	"Email sent successfully"
// @Failure		400		{object}	entity.Response	"Bad Request"
// @Failure		401		{object}	entity.Response	"Unauthorized"
// @Failure		404		{object}	entity.Response	"User not found"
// @Failure		500		{object}	entity.Response	"Internal Server Error"
// @Router			/auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload entity.ForgotPasswordRequest

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	user, err := ac.UserService.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	code, err := ac.EmailService.SendResetPasswordEmail(user)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ac.UserService.StoreVerificationData(ctx, user.ID, user.Email, code)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Message: "Email sent successfully",
		Status:  http.StatusOK,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Reset the user's password using a reset token
// @Description	Reset the user's password by providing a valid reset token and the new password.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		entity.ResetPasswordRequest	true	"User's email"
// @Success		200		{object}	entity.Response	"Password updated successfully"
// @Failure		400		{object}	entity.Response	"Bad Request"
// @Failure		401		{object}	entity.Response	"Unauthorized"
// @Failure		500		{object}	entity.Response	"Internal Server Error"
// @Router			/auth/reset-password/{resetToken} [put]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload entity.ResetPasswordRequest
	resetToken := ctx.Params.ByName("resetToken")

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	userID, _, err := ac.UserService.VerifyCodeToResetPassword(ctx, resetToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ac.UserService.UpdatePassword(ctx, payload.Password, payload.PasswordConfirmation, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Message: "Password updated successfully",
		Status:  http.StatusOK,
	}

	ctx.JSON(http.StatusOK, response)
}
