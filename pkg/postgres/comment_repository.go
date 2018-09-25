package postgres

import (
	"context"
	"database/sql"

	"github.com/godwhoa/upboat/pkg/comments"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

// NewCommentRepository is a constructor
func NewCommentRepository(db *sql.DB) comments.Repository {
	return &CommentRepository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *comments.Comment) (id int, err error) {
	stmt := `
	INSERT INTO comments(post_id, parent_id, commenter_id, depth, body) 
	VALUES($1, $2, $3, calculate_depth($2), $4) RETURNING id`
	err = r.db.QueryRowContext(ctx, stmt,
		comment.PostID, comment.ParentID, comment.CommenterID, comment.Body).
		Scan(&id)
	return
}

func (r *CommentRepository) Comments(ctx context.Context, postID int) (c []*comments.Comment, err error) {
	query := `SELECT id, post_id, parent_id, commenter_id, body 
	FROM comments WHERE post_id = $1 AND deleted IS NULL;`
	err = r.db.SelectContext(ctx, &c, query, postID)
	if err == sql.ErrNoRows {
		err = comments.ErrCommentNotFound
	}
	return
}

func (r *CommentRepository) Delete(ctx context.Context, commentID int, commenterID int) error {
	stmt := `UPDATE comments SET deleted = now() WHERE id = $1 AND commenter_id = $2`

	result, err := r.db.ExecContext(ctx, stmt, commentID, commenterID)
	affected, _ := result.RowsAffected()
	if IsForeignKeyViolation(err) {
		return comments.ErrCommentNotFound
	}
	if err == nil && affected < 1 {
		return comments.ErrUnauthorized
	}
	return err
}

func (r *CommentRepository) Vote(ctx context.Context, commentID int, voterID int, delta int) error {
	stmt := `
	INSERT INTO comment_votes(comment_id, voter_id, delta) 
	VALUES 
	($1, $2, $3)
	ON CONFLICT (voter_id, comment_id) DO
	UPDATE SET delta = $3 WHERE comment_votes.comment_id = $1 AND comment_votes.voter_id = $2;
	`
	_, err := r.db.ExecContext(ctx, stmt, commentID, voterID, delta)
	if IsForeignKeyViolation(err) {
		return comments.ErrCommentNotFound
	}
	return err
}

func (r *CommentRepository) Unvote(ctx context.Context, commentID int, voterID int) error {
	stmt := `DELETE FROM comment_votes WHERE comment_id = $1 AND voter_id = $2`

	_, err := r.db.ExecContext(ctx, stmt, commentID, voterID)
	if IsForeignKeyViolation(err) {
		return comments.ErrCommentNotFound
	}
	return err
}

func (r *CommentRepository) Score(ctx context.Context, commentID int) (score int, err error) {
	query := `SELECT COALESCE(SUM(delta), 0) FROM comment_votes WHERE comment_id = $1`

	err = r.db.QueryRowContext(ctx, query, commentID).
		Scan(&score)
	if err == sql.ErrNoRows {
		err = comments.ErrCommentNotFound
	}
	return
}
