package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int           `json:"id,omitempty"`
	Title     string        `json:"title"`
	Body      string        `json:"body"`
	UserID    int           `json:"user_id"`
	Likes     int           `json:"likes,omitempty"`
	ImagesID  []string      `json:"images,omitempty"`
	CreatedAt *sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt *sql.NullTime `json:"updated_at,omitempty"`
}

type LikePost struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Poster interface {
	Create(ctx context.Context, db *pgxpool.Pool) error
	Get(ctx context.Context, db *pgxpool.Pool, id int) error
	Update(ctx context.Context, db *pgxpool.Pool, id int) error
	Delete(ctx context.Context, db *pgxpool.Pool, id int) error
	LikePost(ctx context.Context, db *pgxpool.Pool, user_id, post_id int) error
	GetLikesPost(ctx context.Context, db *pgxpool.Pool, post_id int) error
}

var _ Poster = &Post{}

func GetAllPosts(ctx context.Context, db *pgxpool.Pool) ([]*Post, error) {
	query := `SELECT * FROM posts`
	rows, err := db.Query(ctx, query)
	if err != nil {
		comment := fmt.Errorf("failed to get posts: %s", err)
		slog.Error(comment.Error())
		return nil, comment
	}

	var posts []*Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.UserID, &post.ImagesID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			comment := fmt.Errorf("failed to scan post: %s", err)
			slog.Error(comment.Error())
			return nil, comment
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *Post) Create(ctx context.Context, db *pgxpool.Pool) error {
	query := `INSERT INTO posts (title, body, user_id, images_id) VALUES ($1, $2, $3, $4)`

	if _, err := db.Exec(ctx, query, p.Title, p.Body, p.UserID, p.ImagesID); err != nil {
		comment := fmt.Errorf("failed to insert post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}

func (p *Post) Delete(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	if _, err := db.Exec(ctx, query, id); err != nil {
		comment := fmt.Errorf("failed to delete post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}

func (p *Post) Get(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `SELECT * FROM posts WHERE id = $1`

	if err := db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Body, &p.UserID, &p.ImagesID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		comment := fmt.Errorf("failed to get post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}

func (p *Post) LikePost(ctx context.Context, db *pgxpool.Pool, user_id int, post_id int) error {
	query := `INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)`
	if _, err := db.Exec(ctx, query, user_id, post_id); err != nil {
		comment := fmt.Errorf("failed to like post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}

func (p *Post) GetLikesPost(ctx context.Context, db *pgxpool.Pool, post_id int) error {
	query := `SELECT count(*) FROM post_likes WHERE post_id = $1`
	if err := db.QueryRow(ctx, query, post_id).Scan(&p.Likes); err != nil {
		comment := fmt.Errorf("failed to get likes post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}

func (p *Post) Update(ctx context.Context, db *pgxpool.Pool, id int) error {
	updated_at := time.Now()

	query := `UPDATE posts SET title = $1, body = $2, user_id = $3, images_id = $4, updated_at = $5 WHERE id = $6`
	if _, err := db.Exec(ctx, query, p.Title, p.Body, p.UserID, p.ImagesID, updated_at, id); err != nil {
		comment := fmt.Errorf("failed to update post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}
