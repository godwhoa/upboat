package postgres

import (
	"database/sql"
	"fmt"

	"github.com/godwhoa/upboat/pkg/postgres/migrations"
	"github.com/godwhoa/upboat/pkg/posts"
	"github.com/godwhoa/upboat/pkg/users"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/lib/pq"
	"github.com/basvanbeek/ocsql"
)

// Options holds information for connecting to a postgres instance
type Options struct {
	User             string
	Pass             string
	Host             string
	Port             int
	DBName           string
	SSLMode          string
	StatementTimeout int
}

// ConnectionInfo returns the postgresql connection string from an options struct.
func (o Options) ConnectionInfo() string {
	return fmt.Sprintf("host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s' statement_timeout=%d", o.Host, o.Port, o.User, o.Pass, o.DBName, o.SSLMode, o.StatementTimeout)
}

// NewFromOptions will connect to a postgresql server with given options
func NewFromOptions(options Options) (*Repositories, error) {
	driverName, err := ocsql.Register("postgres", ocsql.WithAllTraceOptions())
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName, options.ConnectionInfo())
	if err != nil {
		return nil, err
	}
	return New(db)
}

// Repositories is a container for multiple setup repositories (eg. User, Posts etc.)
type Repositories struct {
	UserRepo users.Repository
	PostRepo posts.Repository
}

// New runs migrations and returns wired-up Repositories
func New(db *sql.DB) (*Repositories, error) {
	if err := Migrate(db); err != nil {
		return nil, err
	}
	return &Repositories{
		UserRepo: NewUserRepository(db),
		PostRepo: NewPostRepository(db),
	}, nil
}

// Migrate runs migrations on the database
func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "migrations", DatabaseName: "upboat"})
	if err != nil {
		return err
	}

	assetsrc := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	srcdriver, err := bindata.WithInstance(assetsrc)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", srcdriver, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != migrate.ErrNoChange && err != nil {
		return err
	}
	return nil
}

// IsUniqueKeyViolation checks if an error was caused by an unique key violation
func IsUniqueKeyViolation(err error) bool {
	pqerr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqerr.Code.Name() == "unique_violation"
}

// IsForeignKeyViolation checks if an error was caused by a foreign key violation
func IsForeignKeyViolation(err error) bool {
	pqerr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqerr.Code.Name() == "foreign_key_violation"
}
