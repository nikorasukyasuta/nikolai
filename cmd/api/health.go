package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"status":"ok",
		"env": app.config.env,
		"version": app.config.version,
	}
	if err := writeJSON(res,http.StatusOK,data); err!=nil {
		app.internalServerError(res,req,err)
	}
}
