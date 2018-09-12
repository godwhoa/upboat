package posts

import (
	"context"
	"errors"
)

// Post models a post
type Post struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

var (
	// ErrPostNotFound for when post is not found.
	ErrPostNotFound = errors.New("Post not found")
	// ErrUnauthorized for when a user tries to delete or edit of another user
	ErrUnauthorized = errors.New("Unauthorized to delete/edit the post")
)

// Repository handles storing posts and their votes
type Repository interface {
	Create(ctx context.Context, post *Post) (id int, err error)
	Get(ctx context.Context, postID int) (post *Post, err error)
	Edit(ctx context.Context, post *Post) error
	Delete(ctx context.Context, authorID, postID int) error
	// Vote upserts a vote for a specific post/user
	Vote(ctx context.Context, postID, voterID, delta int) error
	Unvote(ctx context.Context, postID, voterID int) error
	// Votes fetches votes on a specific post
	Votes(ctx context.Context, postID int) (votes int, err error)
}

// Service is a thin layer around Repository which sanitizes Title and Body
type Service interface {
	Repository
}
