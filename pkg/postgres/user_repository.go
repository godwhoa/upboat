package postgres

import (
	"database/sql"

	"github.com/godwhoa/upboat/pkg/users"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

// UserRepository implements users.UserRepository interface
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository is a constructor
func NewUserRepository(db *sql.DB) users.Repository {
	return &UserRepository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

// Create creates a new user
func (repo *UserRepository) Create(user *users.User) error {
	uid := uuid.Must(uuid.NewV4()).String()
	_, err := repo.db.Exec(`INSERT INTO users(uid, username, email, hash) VALUES($1, $2, $3, $4)`,
		uid, user.Username, user.Email, user.Hash)
	if IsUniqueKeyViolation(err) {
		return users.ErrUserAlreadyExists
	}

	return err
}

// Find finds an user by id
func (repo *UserRepository) Find(id int) (*users.User, error) {
	user := &users.User{}
	err := repo.db.QueryRow(`SELECT id, username, email, hash FROM users WHERE id = $1;`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Hash)
	if err == sql.ErrNoRows {
		return nil, users.ErrUserNotFound
	}
	return user, err
}

// FindByEmail finds by email
func (repo *UserRepository) FindByEmail(email string) (*users.User, error) {
	user := &users.User{}
	err := repo.db.QueryRow(`SELECT id, username, email, hash FROM users WHERE email = $1;`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Hash)
	if err == sql.ErrNoRows {
		return nil, users.ErrUserNotFound
	}
	return user, err
}
