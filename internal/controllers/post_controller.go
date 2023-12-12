package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"

	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/core/error"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService domain.PostRepository
	Env         *config.Config
}

// @Summary	Create a post
// @Schemes
// @Description	create a post on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Param			request	body		domain.Post	true	"Post"
// @Success		200		{object}	domain.Response
// @Router			/post/create [post]
func (ps *PostController) Post(ctx *gin.Context) {
	var payload domain.Post
	userID := ctx.GetString("x-user-id")

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(error.BadRequest)
		return
	}

	payload.UserID = &userID

	err = ps.PostService.Create(ctx, &payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := domain.Response{
		Status:  http.StatusCreated,
		Message: "Post created successfully",
	}

	ctx.JSON(http.StatusCreated, response)
}
