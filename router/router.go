package router

import (
	"projetpostgre/internal/handler"

	"github.com/go-chi/chi"
)

func SetupRouter(userHandler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/users", userHandler.CreateHandler)
	router.Get("/api/users/{id}", userHandler.GetByIDHandler)
	router.Put("/api/users/{id}", userHandler.UpdateHandler)
	router.Delete("/api/users/{id}", userHandler.DeleteHandler)

	router.Get("/api/users", userHandler.ListHandler)

	return router
}
