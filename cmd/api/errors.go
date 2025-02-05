package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(res http.ResponseWriter, req *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	writeJSONError(res, http.StatusInternalServerError, "The Server encountered a problem")
}

func (app *application) badRequestResponse(res http.ResponseWriter, req *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	writeJSONError(res, http.StatusBadRequest, err.Error())
}

func (app *application) resourceNotFound(res http.ResponseWriter, req *http.Request, err error) {
	log.Printf("Resource not found error: %s path: %s error: %s", req.Method, req.URL.Path, err)
	writeJSONError(res, http.StatusNotFound, "Resource not Found Error")
}