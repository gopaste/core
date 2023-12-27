package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator) {
	pr := repository.NewPostRepository(db)

	postService := services.NewPostService(pr, validation)

	pc := &controllers.PostController{
		PostService: postService,
		Env:         cfg,
	}

	group.POST("/post/create", pc.Post)
	group.GET("/post/all", pc.GetPosts)
	group.DELETE("/post/:id", pc.DeletePost)
	group.PATCH("/post/:id", pc.UpdatePost)
}
