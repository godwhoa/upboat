package api

import (
	"net/http"

	"github.com/godwhoa/upboat/pkg/posts"
	"go.uber.org/zap"
)

// PostsAPI contains all the handlers releated to posts
type PostsAPI struct {
	service posts.Service
	log     *zap.Logger
}

// NewPostsAPI takes in all the deps. and constructs a type with all the handlers
func NewPostsAPI(service posts.Service, log *zap.Logger) *PostsAPI {
	return &PostsAPI{
		service: service,
		log:     log,
	}
}

// Create creates a new post
func (p *PostsAPI) Create(w http.ResponseWriter, r *http.Request) {
	req := &createRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	post := &posts.Post{
		AuthorID: userID,
		Title:    req.Title,
		Body:     req.Body,
	}
	postID, err := p.service.Create(ctx, post)
	if err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Created("Post created!", map[string]int{"post_id": postID}))
}

// Get fetches a specific post by its ID
func (p *PostsAPI) Get(w http.ResponseWriter, r *http.Request) {
	postID := r.Context().Value("post_id").(int)

	post, err := p.service.Get(postID)
	if err == posts.ErrPostNotFound {
		Respond(w, NotFound("No such post"))
		return
	}
	if err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, OkData("Post found", post))
}

// Update updates a specific post of the user
func (p *PostsAPI) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	req := &updateRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	post := &posts.Post{
		ID:       postID,
		AuthorID: userID,
		Title:    req.Title,
		Body:     req.Body,
	}
	err := p.service.Edit(post)
	if err == posts.ErrUnauthorized {
		Respond(w, Unauthorized(err.Error()))
		return
	}
	if err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Ok("Post updated!"))
}

// Delete deletes a specific post of the user
func (p *PostsAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	err := p.service.Delete(userID, postID)
	if err == posts.ErrUnauthorized {
		Respond(w, Unauthorized(err.Error()))
		return
	}
	if err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Ok("Deleted!"))
}

// Votes fetches votes for a specific post
func (p *PostsAPI) Votes(w http.ResponseWriter, r *http.Request) {
	postID := r.Context().Value("post_id").(int)

	votes, err := p.service.Votes(postID)
	if err == posts.ErrPostNotFound {
		Respond(w, NotFound("No such post"))
		return
	}
	if err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Created("Votes for the post", map[string]int{"votes": votes}))
}

// Vote votes on a specific post
func (p *PostsAPI) Vote(w http.ResponseWriter, r *http.Request) {
	postID := r.Context().Value("post_id").(int)

	req := &voteRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	userID := r.Context().Value("user_id").(int)

	if err := p.service.Vote(postID, userID, req.Delta); err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Ok("Voted!"))
}

// Unvote deletes user's vote on a specific post
func (p *PostsAPI) Unvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	if err := p.service.Unvote(postID, userID); err != nil {
		Respond(w, InternalError())
		return
	}

	Respond(w, Ok("Vote removed!"))
}
