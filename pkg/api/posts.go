package api

import (
	"encoding/json"
	"net/http"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/godwhoa/upboat/pkg/posts"
	R "github.com/godwhoa/upboat/pkg/response"
	"go.uber.org/zap"
)

// DecodeValidate decodes JSON, validates it and responds accordingly
func DecodeValidate(w http.ResponseWriter, r *http.Request, obj v.Validatable) (ok bool) {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		R.Respond(w, R.JSONError())
		return false
	}
	if err := obj.Validate(); err != nil {
		R.Respond(w, R.ValidationError(err))
		return false
	}
	return true
}

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
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	req := &createRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	post := &posts.Post{
		AuthorID: userID,
		Title:    req.Title,
		Body:     req.Body,
	}
	postID, err := p.service.Create(ctx, post)
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Created("Post R.Created!", map[string]int{"post_id": postID}))
}

// Get fetches a specific post by its ID
func (p *PostsAPI) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)

	post, err := p.service.Get(ctx, postID)
	if err == posts.ErrPostNotFound {
		R.Respond(w, R.NotFound("No such post"))
		return
	}
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.OkData("Post found", post))
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
	err := p.service.Edit(ctx, post)
	if err == posts.ErrUnauthorized {
		R.Respond(w, R.Unauthorized(err.Error()))
		return
	}
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Ok("Post updated!"))
}

// Delete deletes a specific post of the user
func (p *PostsAPI) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	err := p.service.Delete(ctx, userID, postID)
	if err == posts.ErrUnauthorized {
		R.Respond(w, R.Unauthorized(err.Error()))
		return
	}
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Ok("Deleted!"))
}

// Votes fetches votes for a specific post
func (p *PostsAPI) Votes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)

	votes, err := p.service.Votes(ctx, postID)
	if err == posts.ErrPostNotFound {
		R.Respond(w, R.NotFound("No such post"))
		return
	}
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Created("Votes for the post", map[string]int{"votes": votes}))
}

// Vote votes on a specific post
func (p *PostsAPI) Vote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	req := &voteRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	if err := p.service.Vote(ctx, postID, userID, req.Delta); err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Ok("Voted!"))
}

// Unvote deletes user's vote on a specific post
func (p *PostsAPI) Unvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := ctx.Value("post_id").(int)
	userID := ctx.Value("user_id").(int)

	if err := p.service.Unvote(ctx, postID, userID); err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	R.Respond(w, R.Ok("Vote removed!"))
}
