package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)

	r.Mount("/debug", middleware.Profiler())

	// r.Post("/generate", JsonRoute(DoGenerate))
	// r.Get("/solve/{id}", JsonRoute(DoPuzzleSolveSimple))
	// r.Get("/puzzle/{id}", JsonRoute(DoPuzzleGet))
	// r.Get("/puzzle", JsonRoute(DoPuzzleGenerateSimple))
	// r.Get("/pdf", JsonRoute(DoPuzzlePDFSimple))

	// new routes
	r.Get("/puzzle/{format}/{id}", JsonRoute(DoPuzzleFormatSingle))
	r.Get("/solve/{format}/{id}", JsonRoute(DoSolveFormatSingle))
	r.Get("/generate/{format}", JsonRoute(DoGenerateFormatSingle))

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, GetRouter())
}
