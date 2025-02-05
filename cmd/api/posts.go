package main

import (
	"errors"
	"net/http"
	"nikolai/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct{
	Title string `json:"title"`
	Content string `json:"content"`
	Tags []string `json:"tags"`
}

func(app *application) createPostHandler(res http.ResponseWriter, req *http.Request){
	var payload CreatePostPayload
	if err := readJSON(res,req,&payload); err != nil {
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
	if err := writeJSON(res, http.StatusOK, post); err != nil {
		app.internalServerError(res, req, err)
		return
	}
}