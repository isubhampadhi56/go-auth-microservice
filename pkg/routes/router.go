package router

import (
	"net/http"

	v1router "github.com/go-auth-microservice/pkg/routes/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api", registerRouterVersions())
	return r
}

func registerRouterVersions() http.Handler {
	r := chi.NewRouter()

	r.Mount("/v1", v1router.V1Router())
	return r

}
