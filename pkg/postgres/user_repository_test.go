package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/godwhoa/upboat/pkg/users"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

func setupDB() (*sql.DB, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	database := "upboat"
	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + database})
	if err != nil {
		return nil, nil, err
	}
	purge := func() {
		pool.Purge(resource)
	}

	var db *sql.DB

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), database))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := Migrate(db); err != nil {
		return nil, nil, err
	}

	return db, purge, err
}

func TestUserRepository(t *testing.T) {
	c := qt.New(t)

	db, purge, err := setupDB()
	c.Assert(err, qt.IsNil)
	defer purge()

	userrepo := NewUserRepository(db)
	// Create OK
	err = userrepo.Create(&users.User{
		Username: "pacninja",
		Email:    "pac@pac.com",
		Hash:     "bcrypt_hash",
	})
	c.Assert(err, qt.IsNil)
	// Create AlreadyExists
	err = userrepo.Create(&users.User{
		Username: "pacninja",
		Email:    "pac@pac.com",
		Hash:     "bcrypt_hash",
	})
	c.Assert(err, qt.Equals, users.ErrUserAlreadyExists)
	// FindByEmail OK
	user, err := userrepo.FindByEmail("pac@pac.com")
	c.Assert(err, qt.IsNil)
	c.Assert(user, qt.Not(qt.IsNil))
	c.Assert(user.Username, qt.Equals, "pacninja")
	c.Assert(user.Hash, qt.Equals, "bcrypt_hash")
	id := user.ID
	// FindByEmail User not Found
	_, err = userrepo.FindByEmail("none@none.com")
	c.Assert(err, qt.Equals, users.ErrUserNotFound)
	// Find OK
	user, err = userrepo.Find(id)
	c.Assert(err, qt.IsNil)
	c.Assert(user, qt.Not(qt.IsNil))
	c.Assert(user.Username, qt.Equals, "pacninja")
	c.Assert(user.Hash, qt.Equals, "bcrypt_hash")
	// Find User Not Found
	_, err = userrepo.Find(666)
	c.Assert(err, qt.Equals, users.ErrUserNotFound)
}
