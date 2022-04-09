package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/generate", JsonRoute(DoGenerate))
	r.Get("/solve/{id}", JsonRoute(DoPuzzleSolveSimple))
	r.Get("/puzzle/{id}", JsonRoute(DoPuzzleGet))
	r.Get("/puzzle", JsonRoute(DoPuzzleGenerateSimple))

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, GetRouter())
}
