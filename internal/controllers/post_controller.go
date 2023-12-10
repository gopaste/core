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

	ctx.JSON(http.StatusCreated, "ok")
}
