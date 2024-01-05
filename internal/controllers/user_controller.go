package controllers

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService entity.UserService
	Env         *config.Config
}

// @Summary		Get authenticated user
// @Schemes		http
// @Description	Get details of the authenticated user
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	entity.Response	"User retrieved successfully"
// @Failure		400	{object}	typesystem.Http	"Bad Request"
// @Failure		401	{object}	typesystem.Http	"Unauthorized"
// @Failure		404	{object}	typesystem.Http	"User not found"
// @Router			/user [get]
func (uc *UserController) GetAuthenticatedUser(ctx *gin.Context) {
	userID := ctx.GetString("x-user-id")
	if userID == "" {
		ctx.Error(typesystem.Unauthorized)
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(typesystem.BadRequest)
		return
	}

	user, err := uc.UserService.GetUserByID(ctx, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Data:    user,
		Status:  200,
		Message: "User retrieved successfully",
	}

	ctx.JSON(200, response)
}
