package comments

import (
	"context"

	"github.com/godwhoa/upboat/pkg/errors"
)

type Comment struct {
	ID          int    `json:"id"`
	PostID      int    `json:"post_id"`
	ParentID    *int   `json:"parent_id"`
	CommenterID int    `json:"author_id"`
	Body        string `json:"body"`
}

var (
	// ErrCommentNotFound for when comment is not found.
	ErrCommentNotFound = errors.E(errors.NotFound, "Comment(s) not found")
	// ErrPostNotFound for when post is not found.
	ErrPostNotFound = errors.E(errors.NotFound, "Post not found")
	// ErrUnauthorized for when a user tries to delete or edit of another user
	ErrUnauthorized = errors.E(errors.Unauthorized, "Unauthorized to delete/edit the comment.")
)

// TODO: think about depth and limit of nested comments
type Repository interface {
	Create(ctx context.Context, comment *Comment) (id int, err error)
	Comments(ctx context.Context, postID int) ([]*Comment, error)
	Delete(ctx context.Context, commentID, authorID int) error
	Vote(ctx context.Context, commentID, voterID, delta int) error
	Unvote(ctx context.Context, commentID, voterID int) error
	Score(ctx context.Context, commentID int) (score int, err error)
}

type Service interface {
	Repository
}
