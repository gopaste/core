package handlers

import (
	"context"
	"net/http"

	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type PostService interface {
	Create(ctx context.Context, post *entity.PostInput) error
	GetPosts(ctx context.Context, id uuid.UUID, page string) ([]*entity.PostOutput, *entity.PaginationInfo, error)
	GetAllPublics(ctx context.Context, page string) ([]*entity.PostOutput, *entity.PaginationInfo, error)
	DeletePost(ctx context.Context, id string, userID uuid.UUID) error
	UpdatePost(ctx context.Context, post *entity.PostUpdateInput, userID uuid.UUID, id string) error
	SearchPost(ctx context.Context, query string, page string) ([]*entity.PostOutput, *entity.PaginationInfo, error)
	GetPost(ctx context.Context, id string, userID string, password string) (*entity.PostOutput, error)
}

type PostHandler struct {
	PostService PostService
	Env         *config.Config
}

// @Summary				Create a post
// @Schemes
// @Description	create a post on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Param			request	body		entity.Post	true	"Post"
// @Success		200		{object}	entity.Response
// @Router			/post/create [post]
func (ps *PostHandler) Post(ctx *gin.Context) {
	var payload entity.PostInput
	userID := ctx.GetString("x-user-id")

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(typesystem.BadRequest)
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

// @Summary		Get all posts of the logged-in user
// @Schemes		http
// @Description	Get all posts of the logged-in user on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	[]entity.Response
// @Router			/post/all [get]
func (ps *PostHandler) GetPosts(ctx *gin.Context) {
	userID := ctx.GetString("x-user-id")
	pageStr := ctx.Query("page")

	id, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(typesystem.Unauthorized)
		return
	}

	posts, paginationInfo, err := ps.PostService.GetPosts(ctx, id, pageStr)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Posts retrieved successfully",
		Info:    paginationInfo,
		Data:    posts,
	}

	ctx.JSON(http.StatusOK, response)
}

func (ps *PostHandler) GetAllPublics(ctx *gin.Context) {
	pageStr := ctx.Query("page")

	posts, paginationInfo, err := ps.PostService.GetAllPublics(ctx, pageStr)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Posts retrieved successfully",
		Info:    paginationInfo,
		Data:    posts,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Delete a post by ID
// @Schemes		http
// @Description	Delete a post belonging to the logged-in user on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string			true	"Post ID"
// @Success		200	{object}	entity.Response	"Post deleted successfully"
// @Router			/post/{id} [delete]
func (ps *PostHandler) DeletePost(ctx *gin.Context) {
	userID := ctx.GetString("x-user-id")
	postID := ctx.Param("id")

	id, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(typesystem.ServerError)
		return
	}

	err = ps.PostService.DeletePost(ctx, postID, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, entity.Response{
		Status:  http.StatusOK,
		Message: "Post deleted successfully",
	})
}

// @Summary		Update a post by ID
// @Schemes		http
// @Description	Update a post belonging to the logged-in user on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string					true	"Post ID"
// @Param			request	body		entity.PostUpdateInput	true	"Post"
// @Success		200		{object}	entity.Response			"Post updated successfully"
// @Router			/post/{id} [patch]
func (ps *PostHandler) UpdatePost(ctx *gin.Context) {
	var payload entity.PostUpdateInput

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(typesystem.BadRequest)
		return
	}

	userID := ctx.GetString("x-user-id")
	postID := ctx.Param("id")

	id, err := uuid.Parse(userID)
	if err != nil {
		ctx.Error(typesystem.ServerError)
		return
	}

	err = ps.PostService.UpdatePost(ctx, &payload, id, postID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Post updated successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Search a post
// @Schemes		http
// @Description	Search a post on the platform
// @Tags			Post
// @Accept			json
// @Produce		json
// @Param			q	query		string			true	"Query"
// @Success		200	{object}	entity.Response	"Post updated successfully"
// @Router			/post/search   [get]
func (ps *PostHandler) SearchPost(ctx *gin.Context) {
	query := ctx.Query("q")
	page := ctx.Query("page")

	post, paginationInfo, err := ps.PostService.SearchPost(ctx, query, page)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Post retrieved successfully",
		Data:    post,
		Info:    paginationInfo,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Get post by ID
// @Schemes		http
// @Description	Get post by ID
// @Tags			Post
// @Accept			json
// @Produce		json
// @Param			id	path		string			true	"Post ID"
// @Success		200	{object}	entity.Response	"Post retrieved successfully"
// @Failure		400	{object}	typesystem.Http	"Bad Request"
// @Failure		404	{object}	typesystem.Http	"Post not found"
// @Router			/post/{id}   [get]
func (ps *PostHandler) GetPost(ctx *gin.Context) {
	var payload entity.GetPostInput

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.Error(typesystem.BadRequest)
		return
	}

	userID := ctx.GetString("x-user-id")

	id := ctx.Param("id")

	if err != nil {
		ctx.Error(typesystem.BadRequest)
		return
	}

	post, err := ps.PostService.GetPost(ctx, id, userID, payload.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := entity.Response{
		Status:  http.StatusOK,
		Message: "Post retrieved successfully",
		Data:    post,
	}

	ctx.JSON(http.StatusOK, response)
}
