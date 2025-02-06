package main

import (
	"errors"
	"net/http"
	"nikolai/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostPayload struct{
	Title string `json:"title" validator:"required,max=100"`
	Content string `json:"content" validator:"required,max=1000"`
	Tags []string `json:"tags"`
}

func(app *application) createPostHandler(res http.ResponseWriter, req *http.Request){
	var payload PostPayload

	if err := readJSON(res,req,&payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	userID := 1 // TODO Format to take user_id in URL params
	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		UserID: int64(userID),
		Tags: payload.Tags,
	}
	ctx := req.Context()

	if err := app.store.Posts.Create(ctx,post); err != nil {
		app.internalServerError(res, req, err)
		return
	}

	if err := writeJSON(res, http.StatusCreated, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}

func(app *application) getPostHandler(res http.ResponseWriter, req *http.Request){
	postID :=  chi.URLParam(req,"postID")
	id, err := strconv.ParseInt(postID,10,64)
	if err != nil {
		app.internalServerError(res, req, err)
		return
	}
	ctx := req.Context()
	post, err := app.store.Posts.GetByID(ctx,id)
	if err != nil {
		switch{
			case errors.Is(err,store.ErrNotFound):
				app.resourceNotFound(res, req, err)
			default:
				app.internalServerError(res, req, err)
		}
		return
	}

	comments,err := app.store.Comments.GetPostByID(ctx,id)
	if err != nil {
		app.internalServerError(res, req, err)
		return
	}
	post.Comments = comments

	if err := writeJSON(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}

func(app *application) deletePostHandler(res http.ResponseWriter, req *http.Request){
	postID :=  chi.URLParam(req,"postID")
	id, err := strconv.ParseInt(postID,10,64)
	if err != nil {
		app.internalServerError(res, req, err)
		return
	}
	ctx := req.Context()
	err = app.store.Posts.Delete(ctx,id)
	if err != nil {
		switch{
			case errors.Is(err,store.ErrNotFound):
				app.resourceNotFound(res, req, err)
			default:
				app.internalServerError(res, req, err)
		}
		return
	}
	if err := writeJSON(res, http.StatusNoContent, nil); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}

func(app *application) updatePostHandler(res http.ResponseWriter, req *http.Request){
	postID :=  chi.URLParam(req,"postID")
	id, err := strconv.ParseInt(postID,10,64)
	if err != nil {
		app.internalServerError(res, req, err)
		return
	}

	var payload PostPayload

	if err := readJSON(res,req,&payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	ctx := req.Context()
	post, err := app.store.Posts.GetByID(ctx,id)
	if err != nil {
		switch{
			case errors.Is(err,store.ErrNotFound):
				app.resourceNotFound(res, req, err)
			default:
				app.internalServerError(res, req, err)
		}
		return
	}

	if err := writeJSON(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}