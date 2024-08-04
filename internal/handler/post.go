package api

import (
	entity "accessCloude/internal/storage/postgres/entity"
	"context"
	"errors"
	"net/http"
	"strconv"
)

// GetPosts implements ServerInterface.
func (ac *AccessCloude) GetPosts(w http.ResponseWriter, r *http.Request, params GetPostsParams) {
	posts, err := entity.GetAllPosts(context.Background(), ac.DB)
	if err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, posts, 200)
}

// CreatePost implements ServerInterface.
func (ac *AccessCloude) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 30) // 10 MB
	if err != nil {
		http.Error(w, "The uploaded file is too big", http.StatusBadRequest)
		return
	}

	var post entity.Post

	post.Title = r.FormValue("title")
	post.Body = r.FormValue("body")

	if post.Body == "" {
		err := errors.New("body is empty")
		Response(w, err.Error(), 500)
		return
	}

	authorId, err := strconv.Atoi(r.FormValue("author_id"))
	if err != nil {
		Response(w, err.Error(), 500)
		return
	}

	post.AuthorID = authorId

	files := r.MultipartForm.File["images_file"]
	for _, file := range files {
		post.ImagesFile = append(post.ImagesFile, file)
	}

	if err := post.Create(context.Background(), ac.DB); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 201)
}

// DeletePost implements ServerInterface.
func (ac *AccessCloude) DeletePost(w http.ResponseWriter, r *http.Request, id int) {
	var post entity.Post
	if err := post.Delete(context.Background(), ac.DB.PConn, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

// GetPost implements ServerInterface.
func (ac *AccessCloude) GetPost(w http.ResponseWriter, r *http.Request, id int, params GetPostParams) {
	var post entity.Post
	if err := post.Get(context.Background(), ac.DB, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

// LikePost implements ServerInterface.
func (ac *AccessCloude) LikePost(w http.ResponseWriter, r *http.Request, params LikePostParams) {

	var post entity.Post
	if err := post.LikePost(context.Background(), ac.DB.PConn, params.UserId, params.PostId); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

func (ac *AccessCloude) GetLikes(w http.ResponseWriter, r *http.Request, params GetLikesParams) {
	var post entity.Post
	if err := post.GetLikesPost(context.Background(), ac.DB.PConn, params.PostId); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post.Likes, 200)
}

// UpdatePost implements ServerInterface.
func (ac *AccessCloude) UpdatePost(w http.ResponseWriter, r *http.Request, id int) {

	err := r.ParseMultipartForm(20 << 30) // 10 MB
	if err != nil {
		http.Error(w, "The uploaded file is too big", http.StatusBadRequest)
		return
	}

	var post entity.Post

	post.Title = r.FormValue("title")
	post.Body = r.FormValue("body")

	files := r.MultipartForm.File["images_file"]
	for _, file := range files {
		post.ImagesFile = append(post.ImagesFile, file)
	}

	if err := post.Update(context.Background(), ac.DB, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)

}
