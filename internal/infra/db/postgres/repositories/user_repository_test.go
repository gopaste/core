package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/infra/db/postgres"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/pkg/validation"

	validatorv10 "github.com/go-playground/validator/v10"
)

func BenchmarkUserExistsByEmail(b *testing.B) {
	db, err := postgres.New("postgres://root:password@localhost:5432/pastebin")
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	validation := validation.NewValidator(validatorv10.New())
	repo := NewUserRepository(db)

	passwordHasher := &entity.BcryptPasswordHasher{}
	o := services.NewUserService(repo, validation, passwordHasher)

	for i := 0; i < b.N; i++ {
		_, err := o.UserExistsByEmail(context.TODO(), "caixetadev@gmail.com")
		if err != nil {
			b.Fatalf("Erro durante o benchmark: %v", err)
		}
	}
}
