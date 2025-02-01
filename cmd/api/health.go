package main

import "net/http"

func (app *application) healthCheckHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("ok"))
}
