package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/controllers"
	"github.com/Caixetadev/snippet/internal/core/services"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
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

	group.POST("/post", pc.Post)
}
