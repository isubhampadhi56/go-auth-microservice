package v1router

import (
	"net/http"

	"github.com/go-auth-microservice/pkg/controller"
	authMiddleware "github.com/go-auth-microservice/pkg/middleware/auth"
	"github.com/go-chi/chi/v5"
)

func V1Router() http.Handler {
	r := chi.NewRouter()
	r.Mount("/auth", authRouter())
	r.Mount("/", protectedRouter())
	return r
}

func authRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/signup", controller.Signup)
	r.Post("/login", controller.Login)
	r.Get("/token", controller.RefreshAccessToken)
	return r
}

func protectedRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(authMiddleware.AccessTokenVerify)
		r.Get("/me", controller.CheckIfSessionValid)
		r.Get("/user", controller.GetUserData)
		r.Patch("/deactivate", controller.DeActivateUser)
		r.Patch("/changePassword", controller.ChangePassword)
	})
	return r
}
