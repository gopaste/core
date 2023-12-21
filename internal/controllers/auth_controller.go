package controllers

import (
	"fmt"
	"net/http"
	"time"

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

	_, err = lc.UserService.Create(c, payload)
	if err != nil {
		c.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "User created successfully",
	}

	c.JSON(http.StatusOK, response)
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

	accessToken, _, err := sc.UserService.CreateAccessToken(user, sc.Env.AccessTokenDuration)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	refreshToken, refreshPayload, err := sc.UserService.CreateRefreshToken(ctx, user, sc.Env.RefreshTokenDuration)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	err = sc.UserService.CreateSession(ctx, refreshPayload, refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	signinResponse := entity.SigninResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

// @Summary	Refresh the user's access token
// @Description	Refresh the user's access token by providing a valid refresh token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param        refresh   header      string  true  "refresh token"
// @Success		200		{object}	entity.Response	"Refreshed successfully"
// @Failure		400		{object}	entity.Response	"Bad Request"
// @Failure		401		{object}	entity.Response	"Unauthorized"
// @Failure		500		{object}	entity.Response	"Internal Server Error"
// @Router			/auth/refresh-token [post]
func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.Request.Header.Get("refresh")
	if refreshToken == "" {
		ctx.Error(entity.BadRequest)
		return
	}

	refreshPayload, err := ac.UserService.VerifyToken(ctx, refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	session, err := ac.UserService.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.Name != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.RefreshToken != refreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		// err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, "expired session")
		return
	}

	user := &entity.User{
		ID:   session.ID,
		Name: session.Name,
	}

	accessToken, _, err := ac.UserService.CreateAccessToken(user, ac.Env.AccessTokenDuration)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Refreshed successfully",
		Data:    accessToken,
	}

	ctx.JSON(http.StatusOK, response)
}
