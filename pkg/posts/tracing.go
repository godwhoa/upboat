package posts

import (
	"context"

	"go.opencensus.io/trace"
)

type tracingMiddleware struct {
	service Service
}

// Tracing is a middleware that provides tracing to Service
func Tracing(service Service) Service {
	return &tracingMiddleware{service}
}

func (m *tracingMiddleware) Create(ctx context.Context, post *Post) (id int, err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Create")
	defer span.End()
	return m.service.Create(ctx, post)
}

func (m *tracingMiddleware) Get(ctx context.Context, postID int) (*Post, error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Get")
	defer span.End()
	return m.service.Get(ctx, postID)
}

func (m *tracingMiddleware) Edit(ctx context.Context, post *Post) (err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Edit")
	defer span.End()
	return m.service.Edit(ctx, post)
}

func (m *tracingMiddleware) Delete(ctx context.Context, postID, authorID int) (err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Delete")
	defer span.End()
	return m.service.Delete(ctx, postID, authorID)
}

func (m *tracingMiddleware) Vote(ctx context.Context, postID, voterID, delta int) (err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Vote")
	defer span.End()
	return m.service.Vote(ctx, postID, voterID, delta)
}

func (m *tracingMiddleware) Unvote(ctx context.Context, postID, voterID int) (err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Unvote")
	defer span.End()
	return m.service.Unvote(ctx, postID, voterID)
}

func (m *tracingMiddleware) Votes(ctx context.Context, postID int) (votes int, err error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Votes")
	defer span.End()
	return m.service.Votes(ctx, postID)
}
