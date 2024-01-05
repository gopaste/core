package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BenchmarkSearch(b *testing.B) {
	poolConfig, err := pgxpool.ParseConfig("postgres://root:password@localhost:5432/pastebin")
	if err != nil {
		b.Fatal(err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		b.Fatal(err)
	}

	repo := &postRepository{
		db: db,
	}

	// Substitua esta string de consulta conforme necessário
	query := "CAIXETA"

	limit := 10
	offset := 0

	for i := 0; i < b.N; i++ {
		_, err := repo.Search(context.Background(), query, limit, offset)
		if err != nil {
			b.Fatal(err)
		}
	}
}
