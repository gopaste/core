package repository

import (
	"context"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ services.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Insert(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)"

	_, err := ur.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

// GetUserByEmail make a query in database and return an user or error
func (ur *userRepository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = $1"

	line, err := ur.db.Query(ctx, query, email)
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

func (ur *userRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"

	line, err := ur.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var user entity.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
	} else {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

func (ur *userRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	err := ur.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *userRepository) StoreVerificationData(ctx context.Context, verificationData *entity.VerificationData) error {
	query := "INSERT INTO password_reset (id, user_id, reset_token, expiration_datetime) VALUES ($1, $2, $3, $4)"

	_, err := ur.db.Exec(
		ctx,
		query,
		verificationData.ID,
		verificationData.UserID,
		verificationData.Code,
		verificationData.ExpiresAt,
	)

	return err
}

func (ur *userRepository) VerifyCodeToResetPassword(ctx context.Context, code string) (entity.VerificationData, error) {
	query := "SELECT id, user_id, expiration_datetime FROM password_reset WHERE reset_token = $1"

	line, err := ur.db.Query(
		ctx,
		query,
		code,
	)
	if err != nil {
		return entity.VerificationData{}, err
	}

	defer line.Close()

	var verificationData entity.VerificationData
	if line.Next() {
		if err = line.Scan(verificationData.ID, verificationData.UserID, verificationData.ExpiresAt); err != nil {
			return entity.VerificationData{}, err
		}
	} else {
		return entity.VerificationData{}, pgx.ErrNoRows
	}

	return verificationData, nil
}

func (ur *userRepository) UpdatePassword(ctx context.Context, password string, id uuid.UUID) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"

	_, err := ur.db.Exec(ctx, query, password, id)

	return err
}

func (ur *userRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	query := "INSERT INTO sessions (id, name, refresh_token, user_agent, client_ip, is_blocked, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := ur.db.Exec(
		ctx,
		query,
		session.ID,
		session.Name,
		session.RefreshToken,
		session.UserAgent,
		session.ClientIp,
		session.IsBlocked,
		session.ExpiresAt,
	)

	return err
}

func (ur *userRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	query := "SELECT id, token, user_id, expiration_datetime FROM refresh_tokens WHERE token = $1"

	line, err := ur.db.Query(
		ctx,
		query,
		token,
	)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var refreshToken entity.RefreshToken
	if line.Next() {
		if err = line.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.UserID, &refreshToken.ExpiresAt); err != nil {
			return nil, err
		}
	} else {
		return nil, pgx.ErrNoRows
	}

	return &refreshToken, nil
}

func (ur *userRepository) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	query := "SELECT id, name, refresh_token, is_blocked, is_revoked, expires_at FROM sessions WHERE id = $1 LIMIT 1"

	line, err := ur.db.Query(
		ctx,
		query,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var session entity.Session

	if line.Next() {
		if err = line.Scan(&session.ID, &session.Name, &session.RefreshToken, &session.IsBlocked, &session.IsRevoked, &session.ExpiresAt); err != nil {
			return nil, err
		}
	} else {
		return nil, pgx.ErrNoRows
	}

	return &session, nil
}

func (ur *userRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := "UPDATE sessions SET is_revoked = $1 WHERE refresh_token = $2"

	_, err := ur.db.Exec(ctx, query, true, token)
	if err != nil {
		return err
	}

	return nil
}
