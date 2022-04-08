package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/generate", JsonRoute(DoGenerate))
	r.Post("/test/{tag}", JsonRoute(doTest))

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, GetRouter())
}
