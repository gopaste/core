package passwordhash

import "golang.org/x/crypto/bcrypt"

// PasswordHasher is a password hashing interface
type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

// BcryptPasswordHasher is a password hashing implementation
type BcryptPasswordHasher struct{}

// GenerateFromPassword generates a bcrypt hash from password
func (h *BcryptPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHashAndPassword compares a bcrypt hash with a plaintext password.
func (h *BcryptPasswordHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
