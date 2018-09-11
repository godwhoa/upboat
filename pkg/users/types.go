package users

import "errors"

var (
	// ErrUserAlreadyExists is returned if an user is already registred with the given email or username.
	ErrUserAlreadyExists = errors.New("User already exists")
	// ErrInvalidCredentials is returned if login credentials are invalid.
	ErrInvalidCredentials = errors.New("Invalid login credentials")
	// ErrUserNotFound is returned if user in not found in the database
	ErrUserNotFound = errors.New("User not found")
)

// User models an user
type User struct {
	ID       int
	Email    string
	Username string
	Hash     string
}

// Repository handles storing/retrieving an user
type Repository interface {
	// Create creates a new user, returns ErrUserAlreadyExists if user already exists with same email/username.
	Create(user *User) error
	// FindByEmail finds an user by ID, returns ErrUserNotFound if no user is found
	Find(id int) (*User, error)
	// FindByEmail finds an user by email, returns ErrUserNotFound if no user is found
	FindByEmail(email string) (*User, error)
}

// Service handles creation and authentication of a user
type Service interface {
	Register(u *User, password string) (*User, error)
	Login(email string, password string) (*User, error)
}
