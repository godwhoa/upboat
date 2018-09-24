package comments

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
)

var policy = bluemonday.StrictPolicy().AllowElements("br")

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, comment *Comment) error {
	comment.Body = policy.Sanitize(comment.Body)
	return s.repo.Create(ctx, comment)
}

func (s *service) Comments(ctx context.Context, postID int) ([]*Comment, error) {
	return s.repo.Comments(ctx, postID)
}

func (s *service) Delete(ctx context.Context, commentID, authorID int) error {
	return s.repo.Delete(ctx, commentID, authorID)
}

func (s *service) Vote(ctx context.Context, commentID, voterID, delta int) error {
	// TODO: commenter should not be able to self-vote their comment
	return s.repo.Vote(ctx, commentID, voterID, delta)
}

func (s *service) Unvote(ctx context.Context, commentID, voterID int) error {
	return s.repo.Unvote(ctx, commentID, voterID)
}

func (s *service) Score(ctx context.Context, commentID int) (score int, err error) {
	return s.repo.Score(ctx, commentID)
}
