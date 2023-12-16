package repository

import (
	"context"
	"time"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := ur.db.Exec(
		ctx,
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)
	return err
}

// GetUserByEmail make a query in database and return an user or error
func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	line, err := ur.db.Query(ctx, "SELECT id, name, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var user entity.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	} else {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

func (ur *userRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := ur.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *userRepository) StoreVerificationData(ctx context.Context, verificationData *entity.VerificationData) error {
	_, err := ur.db.Exec(ctx, "INSERT INTO password_reset (id, user_id, reset_token, expiration_datetime) VALUES ($1, $2, $3, $4)", verificationData.ID, verificationData.UserID, verificationData.Code, verificationData.ExpiresAt)
	return err
}

func (ur *userRepository) VerifyCodeToResetPassword(ctx context.Context, code string) (string, bool, error) {
	var user_id string

	query := `
		SELECT
			user_id
		FROM
			password_reset
		WHERE
			reset_token = $1 AND expiration_datetime > $2
	`

	err := ur.db.QueryRow(ctx, query, code, time.Now()).Scan(&user_id)
	if err != nil {
		return "", false, err
	}

	return user_id, true, nil
}

func (ur *userRepository) UpdatePassword(ctx context.Context, password string, id string) error {
	_, err := ur.db.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, id)
	return err
}
