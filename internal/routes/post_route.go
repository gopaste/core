package routes

import (
	"github.com/Caixetadev/snippet/config"
	"github.com/Caixetadev/snippet/internal/handlers"
	repository "github.com/Caixetadev/snippet/internal/infra/db/postgres/repositories"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/pkg/passwordhash"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostRouter(cfg *config.Config, db *pgxpool.Pool, group *gin.RouterGroup, validation validation.Validator) {
	pr := repository.NewPostRepository(db)

	postService := services.NewPostService(pr, validation, &passwordhash.BcryptPasswordHasher{})

	pc := &handlers.PostHandler{
		PostService: postService,
		Env:         cfg,
	}

	group.POST("/post/create", pc.Post)
	group.GET("/post/user/all", pc.GetPosts)
	group.DELETE("/post/:id", pc.DeletePost)
	group.PATCH("/post/:id", pc.UpdatePost)
	group.GET("/post/search", pc.SearchPost)
	group.GET("/post/:id", pc.GetPost)
	group.GET("/post/all", pc.GetAllPublics)
}
