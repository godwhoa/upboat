package postgres

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/godwhoa/upboat/pkg/posts"
	"github.com/godwhoa/upboat/pkg/users"
)

func TestPostRepository(t *testing.T) {
	c := qt.New(t)
	ctx := context.Background()

	db, purge, err := setupDB()
	c.Assert(err, qt.IsNil)
	defer purge()

	userrepo := NewUserRepository(db)
	// setup user
	err = userrepo.Create(&users.User{
		Username: "pacninja",
		Email:    "pac@pac.com",
		Hash:     "bcrypt_hash",
	})
	c.Assert(err, qt.IsNil)

	user, err := userrepo.FindByEmail("pac@pac.com")
	c.Assert(err, qt.IsNil)
	userID := user.ID

	err = userrepo.Create(&users.User{
		Username: "lala",
		Email:    "lala@lala.com",
		Hash:     "bcrypt_hash",
	})
	c.Assert(err, qt.IsNil)

	user2, err := userrepo.FindByEmail("lala@lala.com")
	c.Assert(err, qt.IsNil)
	user2ID := user2.ID

	postrepo := NewPostRepository(db)

	// Create OK
	postID, err := postrepo.Create(ctx, &posts.Post{
		AuthorID: userID,
		Title:    "Testing!",
		Body:     "Testing body",
	})
	c.Assert(err, qt.IsNil)

	// Get
	post, err := postrepo.Get(ctx, postID)
	c.Assert(err, qt.IsNil)
	c.Assert(post.AuthorID, qt.Equals, userID)

	// Update
	err = postrepo.Edit(ctx, &posts.Post{
		ID:       post.ID,
		AuthorID: post.AuthorID,
		Body:     "Updated",
		Title:    "Updated",
	})
	c.Assert(err, qt.IsNil)

	// Verify update
	upost, err := postrepo.Get(ctx, post.ID)
	c.Assert(err, qt.IsNil)
	c.Assert(upost.Title, qt.Equals, "Updated")
	c.Assert(upost.Body, qt.Equals, "Updated")

	// Vote
	err = postrepo.Vote(ctx, post.ID, userID, +1)
	c.Assert(err, qt.IsNil)
	err = postrepo.Vote(ctx, post.ID, user2ID, -1)
	c.Assert(err, qt.IsNil)
	err = postrepo.Vote(ctx, post.ID, user2ID, +1)
	c.Assert(err, qt.IsNil)

	// Get vote
	votes, err := postrepo.Votes(ctx, post.ID)
	c.Assert(err, qt.IsNil)
	c.Assert(votes, qt.Equals, 2)

	// Unvote
	err = postrepo.Unvote(ctx, post.ID, userID)
	c.Assert(err, qt.IsNil)
	err = postrepo.Unvote(ctx, post.ID, user2ID)
	c.Assert(err, qt.IsNil)
	// Verify
	votes, err = postrepo.Votes(ctx, post.ID)
	c.Assert(err, qt.IsNil)
	c.Assert(votes, qt.Equals, 0)

	// Delete post
	err = postrepo.Delete(ctx, userID, postID)
	c.Assert(err, qt.IsNil)
	_, err = postrepo.Get(ctx, postID)
	c.Assert(err, qt.Equals, posts.ErrPostNotFound)
}
