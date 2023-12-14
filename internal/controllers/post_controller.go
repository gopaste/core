package controllers

import (
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService entity.PostRepository
	Env         *config.Config
}

// @Summary	Create a post
// @Schemes
// @Description	create a post on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Param			request	body		entity.Post	true	"Post"
// @Success		200		{object}	entity.Response
// @Router			/post/create [post]
func (ps *PostController) Post(ctx *gin.Context) {
	var payload entity.Post
	userID := ctx.GetString("x-user-id")

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(entity.BadRequest)
		return
	}

	payload.UserID = &userID

	err = ps.PostService.Create(ctx, &payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusCreated,
		Message: "Post created successfully",
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary	Get all posts of the logged-in user
// @Schemes		http
// @Description	Get all posts of the logged-in user on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200		{object}	[]entity.Response
// @Router			/post/all [get]
func (ps *PostController) GetPosts(ctx *gin.Context) {
	userID := ctx.GetString("x-user-id")

	posts, err := ps.PostService.GetPosts(ctx, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Posts retrieved successfully",
		Data:    posts,
	}

	ctx.JSON(http.StatusOK, response)
}
