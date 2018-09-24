package postgres

import (
	"context"
	"database/sql"

	"github.com/godwhoa/upboat/pkg/posts"
	"github.com/jmoiron/sqlx"
)

// PostRepository implements `posts.Repository` interface
type PostRepository struct {
	db *sqlx.DB
}

// NewPostRepository is a constructor
func NewPostRepository(db *sql.DB) posts.Repository {
	return &PostRepository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (repo *PostRepository) Create(ctx context.Context, post *posts.Post) (id int, err error) {
	stmt := `INSERT INTO posts(author_id, title, body) VALUES($1, $2, $3) RETURNING id`

	err = repo.db.QueryRowContext(ctx, stmt,
		post.AuthorID, post.Title, post.Body).
		Scan(&id)
	return
}

func (repo *PostRepository) Get(ctx context.Context, postID int) (*posts.Post, error) {
	query := `SELECT id, author_id, body, title FROM posts WHERE id = $1 AND deleted IS NULL;`

	post := &posts.Post{}
	err := repo.db.QueryRowContext(ctx, query, postID).
		Scan(&post.ID, &post.AuthorID, &post.Title, &post.Body)
	if err == sql.ErrNoRows {
		return nil, posts.ErrPostNotFound
	}
	return post, err
}

func (repo *PostRepository) Edit(ctx context.Context, post *posts.Post) error {
	stmt := `UPDATE posts SET title = $1, body = $2 WHERE id = $3 AND author_id = $4 AND deleted IS NULL;`

	result, err := repo.db.Exec(stmt,
		post.Title, post.Body, post.ID, post.AuthorID)
	affected, _ := result.RowsAffected()
	if IsForeignKeyViolation(err) {
		return posts.ErrPostNotFound
	}
	if err == nil && affected < 1 {
		return posts.ErrUnauthorized
	}
	return err
}

func (repo *PostRepository) Delete(ctx context.Context, authorID, postID int) error {
	stmt := `UPDATE posts SET deleted = now() WHERE id = $1 AND author_id = $2`

	result, err := repo.db.ExecContext(ctx, stmt, postID, authorID)
	affected, _ := result.RowsAffected()
	if IsForeignKeyViolation(err) {
		return posts.ErrPostNotFound
	}
	if err == nil && affected < 1 {
		return posts.ErrUnauthorized
	}
	return err
}

func (repo *PostRepository) Vote(ctx context.Context, postID int, voterID int, delta int) error {
	stmt := `
	INSERT INTO post_votes(post_id, voter_id, delta) 
	VALUES 
	($1, $2, $3)
	ON CONFLICT (voter_id, post_id) DO
	UPDATE SET delta = $3 WHERE post_votes.post_id = $1 AND post_votes.voter_id = $2;
	`

	_, err := repo.db.ExecContext(ctx, stmt, postID, voterID, delta)
	if IsForeignKeyViolation(err) {
		return posts.ErrPostNotFound
	}
	return err
}

func (repo *PostRepository) Unvote(ctx context.Context, postID int, voterID int) error {
	stmt := `DELETE FROM post_votes WHERE post_id = $1 AND voter_id = $2`

	_, err := repo.db.ExecContext(ctx, stmt, postID, voterID)
	if IsForeignKeyViolation(err) {
		return posts.ErrPostNotFound
	}
	return err
}

func (repo *PostRepository) Score(ctx context.Context, postID int) (score int, err error) {
	query := `SELECT COALESCE(SUM(delta), 0) FROM post_votes WHERE post_id = $1`

	err = repo.db.QueryRowContext(ctx, query, postID).
		Scan(&score)
	if err == sql.ErrNoRows {
		err = posts.ErrPostNotFound
	}
	return
}
