package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/julimonteiro/cupcake-store/internal/handler"
	"github.com/julimonteiro/cupcake-store/internal/repository"
	"github.com/julimonteiro/cupcake-store/internal/service"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	cupcakeRepo := repository.NewCupcakeRepository(db)
	cupcakeService := service.NewCupcakeService(cupcakeRepo)
	cupcakeHandler := handler.NewCupcakeHandler(cupcakeService)

	r.Get("/health", cupcakeHandler.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/cupcakes", func(r chi.Router) {
			r.Get("/", cupcakeHandler.GetAllCupcakes)
			r.Post("/", cupcakeHandler.CreateCupcake)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", cupcakeHandler.GetCupcake)
				r.Put("/", cupcakeHandler.UpdateCupcake)
				r.Delete("/", cupcakeHandler.DeleteCupcake)
			})
		})
	})

	r.Handle("/", http.FileServer(http.Dir("web")))

	return r
}
