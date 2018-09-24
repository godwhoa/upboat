package posts

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
)

var policy = bluemonday.StrictPolicy().AllowElements("br")

type service struct {
	repo Repository
}

// NewService is a constructor for user.Service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, post *Post) (int, error) {
	post.Title = policy.Sanitize(post.Title)
	post.Body = policy.Sanitize(post.Body)
	return s.repo.Create(ctx, post)
}

func (s *service) Get(ctx context.Context, postID int) (*Post, error) {
	return s.repo.Get(ctx, postID)
}

func (s *service) Edit(ctx context.Context, post *Post) error {
	post.Title = policy.Sanitize(post.Title)
	post.Body = policy.Sanitize(post.Body)
	return s.repo.Edit(ctx, post)
}

func (s *service) Delete(ctx context.Context, postID, authorID int) error {
	return s.repo.Delete(ctx, postID, authorID)
}

func (s *service) Vote(ctx context.Context, postID, voterID, delta int) error {
	return s.repo.Vote(ctx, postID, voterID, delta)
}

func (s *service) Unvote(ctx context.Context, postID, voterID int) error {
	return s.repo.Unvote(ctx, postID, voterID)
}

func (s *service) Score(ctx context.Context, postID int) (int, error) {
	return s.repo.Score(ctx, postID)
}
