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

func (s service) Create(ctx context.Context, post *Post) (int, error) {
	post.Title = policy.Sanitize(post.Title)
	post.Body = policy.Sanitize(post.Body)
	return s.repo.Create(ctx, post)
}

func (s service) Get(postID int) (*Post, error) {
	return s.repo.Get(postID)
}

func (s service) Edit(post *Post) error {
	post.Title = policy.Sanitize(post.Title)
	post.Body = policy.Sanitize(post.Body)
	return s.repo.Edit(post)
}

func (s service) Delete(postID, authorID int) error {
	return s.repo.Delete(postID, authorID)
}

func (s service) Vote(postID, voterID, delta int) error {
	return s.repo.Vote(postID, voterID, delta)
}

func (s service) Unvote(postID, voterID int) error {
	return s.repo.Unvote(postID, voterID)
}

func (s service) Votes(postID int) (int, error) {
	return s.repo.Votes(postID)
}
