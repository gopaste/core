package domain

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	// Create(user *User) error
}
