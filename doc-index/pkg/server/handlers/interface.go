package handlers

import "net/http"

type HandlersInterface interface {
	Healthcheck(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}
