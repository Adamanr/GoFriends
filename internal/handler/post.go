package api

import (
	entity "accessCloude/internal/storage/entity"
	"context"
	"net/http"
)

// GetPosts implements ServerInterface.
func (ac *AccessCloude) GetPosts(w http.ResponseWriter, r *http.Request, params GetPostsParams) {
	posts, err := entity.GetAllPosts(context.Background(), ac.DB.Conn)
	if err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, posts, 200)
}

// CreatePost implements ServerInterface.
func (ac *AccessCloude) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post entity.Post
	if err := UnmarshalObject(r, &post); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	if err := post.Create(context.Background(), ac.DB.Conn); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 201)
}

// DeletePost implements ServerInterface.
func (ac *AccessCloude) DeletePost(w http.ResponseWriter, r *http.Request, id int) {
	var post entity.Post
	if err := post.Delete(context.Background(), ac.DB.Conn, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

// GetPost implements ServerInterface.
func (ac *AccessCloude) GetPost(w http.ResponseWriter, r *http.Request, id int, params GetPostParams) {
	var post entity.Post
	if err := post.Get(context.Background(), ac.DB.Conn, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

// LikePost implements ServerInterface.
func (ac *AccessCloude) LikePost(w http.ResponseWriter, r *http.Request, params LikePostParams) {

	var post entity.Post
	if err := post.LikePost(context.Background(), ac.DB.Conn, params.UserId, params.PostId); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)
}

func (ac *AccessCloude) GetLikes(w http.ResponseWriter, r *http.Request, params GetLikesParams) {
	var post entity.Post
	if err := post.GetLikesPost(context.Background(), ac.DB.Conn, params.PostId); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post.Likes, 200)
}

// UpdatePost implements ServerInterface.
func (ac *AccessCloude) UpdatePost(w http.ResponseWriter, r *http.Request, id int) {
	var post entity.Post
	if err := UnmarshalObject(r, &post); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	if err := post.Update(context.Background(), ac.DB.Conn, id); err != nil {
		Response(w, err.Error(), 500)
		return
	}

	Response(w, post, 200)

}
