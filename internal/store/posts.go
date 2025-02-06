package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Comments []Comment `json:"comments"`
	User    User     	`json:"user"`
}

type PostgresPostsStore struct {
	db *sql.DB
}

func (store *PostgresPostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content,title,user_id,tags)
	VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at
	`

	err := store.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresPostsStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, content, title, user_id, tags, created_at, updated_at
	FROM posts
	WHERE id = $1
	`

	var post Post
	err := store.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNotFound
		default: return nil, err
		}
	}

	return &post, nil
}

func (store *PostgresPostsStore) Delete(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts
	WHERE id = $1
	`

	res, err := store.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Rows affected: %d", rowsAffected)
	return nil
}

func (store *PostgresPostsStore) UpdatePostByID(ctx context.Context, postID int64, post *Post) (*Post, error){
	query := `
	UPDATE posts
	SET content = $1, title = $2, tags = $3
	WHERE id = $4
	RETURNING updated_at
	`

	err := store.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), postID).Scan(&post.UpdatedAt)
	if err != nil {
		return nil,err
	}
	return post, nil
}
