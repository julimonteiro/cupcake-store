package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "Link")
			w.Header().Set("Access-Control-Max-Age", "300")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

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
