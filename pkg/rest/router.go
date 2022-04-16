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

	r.Get("/puzzle/{format}/{id}", JsonRoute(DoPuzzleFormatSingle))

	r.Get("/solve/{format}/{id}", JsonRoute(DoSolveFormatSingle))
	r.Post("/solve/{format}", JsonRoute(DoSolveFormatComplex))

	r.Get("/generate/{format}", JsonRoute(DoGenerateFormatSingle))
	r.Post("/generate/{format}", JsonRoute(DoGenerateFormatMany))

	r.Post("/solutions/{format}", JsonRoute(DoSolutionsFormatComplex))

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, GetRouter())
}
