package main

import (
	"fmt"
	"log"

	"github.com/Caixetadev/snippet/config"
	_ "github.com/Caixetadev/snippet/docs"
	"github.com/Caixetadev/snippet/internal/app"
	"github.com/Caixetadev/snippet/internal/infra/db/postgres"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/middleware/http"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cfg, err := config.NewConfig(".")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	validation := validation.NewValidator(validatorv10.New())

	tokenMaker, err := token.NewPasetoMaker(cfg.AccessTokenSecret)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - token.NewPasetoMaker: %w", err))
	}

	db, err := postgres.New(cfg.DBURL)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	defer db.Close()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(http.ErrorHandler())

	app.Run(cfg, db, router, validation, tokenMaker)

	router.Run(":8080")
}
