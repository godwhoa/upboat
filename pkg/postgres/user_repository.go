package postgres

import (
	"context"
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
func (repo *UserRepository) Create(ctx context.Context, user *users.User) error {
	stmt := `INSERT INTO users(uid, username, email, hash) VALUES($1, $2, $3, $4)`

	uid := uuid.Must(uuid.NewV4()).String()
	_, err := repo.db.ExecContext(ctx, stmt,
		uid, user.Username, user.Email, user.Hash)
	if IsUniqueKeyViolation(err) {
		return users.ErrUserAlreadyExists
	}

	return err
}

// Find finds an user by id
func (repo *UserRepository) Find(ctx context.Context, id int) (*users.User, error) {
	query := `SELECT id, username, email, hash FROM users WHERE id = $1;`

	user := &users.User{}
	err := repo.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Hash)
	if err == sql.ErrNoRows {
		return nil, users.ErrUserNotFound
	}
	return user, err
}

// FindByEmail finds by email
func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	query := `SELECT id, username, email, hash FROM users WHERE email = $1;`

	user := &users.User{}
	err := repo.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Hash)
	if err == sql.ErrNoRows {
		return nil, users.ErrUserNotFound
	}
	return user, err
}
