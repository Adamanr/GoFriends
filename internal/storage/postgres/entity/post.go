package storage

import (
	"accessCloude/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type Post struct {
	ID         int64                   `json:"id"`
	AuthorID   int                     `json:"author_id"`
	Title      string                  `json:"title,omitempty"`
	Body       string                  `json:"body,omitempty"`
	ImagesUrl  []string                `json:"images_url,omitempty"`
	ImagesFile []*multipart.FileHeader `json:"images_file,omitempty"`
	ImagesName []string                `json:"images_name,omitempty"`
	Likes      *int                    `json:"likes,omitempty"`
	CreatedAt  time.Time               `json:"created_at,omitempty"`
	UpdatedAt  *time.Time              `json:"updated_at,omitempty"`
}

type LikePost struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Poster interface {
	Create(ctx context.Context, db *storage.Database) error
	Get(ctx context.Context, db *storage.Database, id int) error
	Update(ctx context.Context, db *storage.Database, id int) error
	Delete(ctx context.Context, db *pgxpool.Pool, id int) error
	LikePost(ctx context.Context, db *pgxpool.Pool, user_id, post_id int) error
	GetLikesPost(ctx context.Context, db *pgxpool.Pool, post_id int) error
}

var _ Poster = &Post{}

func GetAllPosts(ctx context.Context, db *storage.Database) ([]*Post, error) {
	query := `SELECT * FROM posts`
	rows, err := db.PConn.Query(ctx, query)
	if err != nil {
		comment := fmt.Errorf("failed to get posts: %s", err)
		slog.Error(comment.Error())
		return nil, comment
	}

	var posts []*Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.AuthorID, &p.ImagesName, &p.Likes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			comment := fmt.Errorf("failed to scan post: %s", err)
			slog.Error(comment.Error())
			return nil, comment
		}

		fmt.Println("Post", p)
		for k, v := range p.ImagesName {
			fmt.Println(v)
			reqParams := make(url.Values)
			reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%q", v))

			presignedURL, err := db.MConn.PresignedGetObject(context.Background(), "images", v, time.Second*24*60*60, reqParams)
			if err != nil {
				return nil, err
			}

			print("Params", reqParams)
			// Нужно удалить имя файла
			p.ImagesName = append(p.ImagesName[:k], p.ImagesName[k+1:]...)

			url := strings.ReplaceAll(presignedURL.String(), "127.0.0.1:9000", "10.0.2.2:9000")
			url, _, _ = strings.Cut(url, "?")
			p.ImagesUrl = append(p.ImagesUrl, url)
			fmt.Println("Successfully generated presigned URL", url)
		}

		posts = append(posts, &p)
		fmt.Println("Author", p.AuthorID)
	}

	return posts, nil
}

func (p *Post) Create(ctx context.Context, db *storage.Database) error {
	for _, v := range p.ImagesFile {
		uploadFile, err := v.Open()
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%s_%s", p.CreatedAt.String(), v.Filename)
		if _, putErr := db.MConn.PutObject(context.Background(), "images", filename, uploadFile, v.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"}); putErr != nil {
			return err
		}

		uploadFile.Close()

		p.ImagesName = append(p.ImagesName, filename)
	}

	query := `INSERT INTO posts (title, body, user_id, images_name) VALUES ($1, $2, $3, $4)`

	if _, err := db.PConn.Exec(ctx, query, p.Title, p.Body, p.AuthorID, p.ImagesName); err != nil {
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

func (p *Post) Get(ctx context.Context, db *storage.Database, id int) error {
	query := `SELECT * FROM posts WHERE id = $1`

	if err := db.PConn.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Body, &p.AuthorID, &p.ImagesName, &p.Likes, &p.CreatedAt, &p.UpdatedAt); err != nil {
		comment := fmt.Errorf("failed to get post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	for _, v := range p.ImagesName {
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%q", v))

		presignedURL, err := db.MConn.PresignedGetObject(context.Background(), "images", v, time.Second*24*60*60, reqParams)
		if err != nil {
			return err
		}

		p.ImagesUrl = append(p.ImagesUrl, presignedURL.String())
		fmt.Println("Successfully generated presigned URL", presignedURL)
	}

	return nil
}

func (p *Post) LikePost(ctx context.Context, db *pgxpool.Pool, user_id int, post_id int) error {
	// query := `INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)`
	// if _, err := db.Exec(ctx, query, user_id, post_id); err != nil {
	// 	comment := fmt.Errorf("failed to like post: %s", err)
	// 	slog.Error(comment.Error())
	// 	return comment
	// }

	return nil
}

func (p *Post) GetLikesPost(ctx context.Context, db *pgxpool.Pool, post_id int) error {
	// query := `SELECT count(*) FROM post_likes WHERE post_id = $1`
	// if err := db.QueryRow(ctx, query, post_id).Scan(&p.Likes); err != nil {
	// 	comment := fmt.Errorf("failed to get likes post: %s", err)
	// 	slog.Error(comment.Error())
	// 	return comment
	// }

	return nil
}

func (p *Post) Update(ctx context.Context, db *storage.Database, id int) error {
	fmt.Printf("\nID: %v\n\n", id)

	if p.ImagesFile != nil {
		for _, v := range p.ImagesFile {
			uploadFile, err := v.Open()
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("%s_%s", p.CreatedAt.String(), v.Filename)
			if _, putErr := db.MConn.PutObject(context.Background(), "images", filename, uploadFile, v.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"}); putErr != nil {
				return err
			}

			uploadFile.Close()

			p.ImagesName = append(p.ImagesName, filename)
		}
	}

	var (
		query string
		args  pgx.NamedArgs
	)
	if len(p.ImagesName) == 0 {
		query = `UPDATE posts SET title = @title, body = @body, updated_at = current_timestamp WHERE id = @id`
		args = pgx.NamedArgs{
			"title": p.Title,
			"body":  p.Body,
			"id":    id,
		}
	} else {
		query = `UPDATE posts SET title = @title, body = @body, images_name = @images_name, updated_at = current_timestamp WHERE id = @id`
		args = pgx.NamedArgs{
			"title":       p.Title,
			"body":        p.Body,
			"images_name": p.ImagesName,
			"id":          id,
		}
	}

	if _, err := db.PConn.Exec(ctx, query, args); err != nil {
		comment := fmt.Errorf("\nfailed to update post: %s", err)
		slog.Error(comment.Error())
		return comment
	}

	return nil
}
