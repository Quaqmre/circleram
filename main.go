package main

import (
	"net/http"
	"os"

	"github.com/Quaqmre/circleramkit/auth"
	"github.com/Quaqmre/circleramkit/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	//Make User Service tree
	var userService user.Service
	{
		userService = user.NewUserService()
		userService = user.NewLoggingService(logger, userService)
	}

	//Make Auth Service tree
	var authService auth.Service
	{
		authService = auth.NewAuthService("mysecret", userService)
		authService = auth.NewLoggingService(logger, authService)
	}

	//Middleware
	var authMiddleware = auth.Middleware(authService)

	//Make Handler each endpoint
	ch, chpwless := user.MakeHandler(userService)
	//Make Handler each endpoint
	ach := auth.MakeHandler(authService)

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             3599, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	r.Mount("/auth", ach)
	r.Mount("/", chpwless)
	r.Group(func(r chi.Router) {
		// r.Use(authService.Handler)
		r.Use(authMiddleware)
		r.Use(cors.Handler)
		r.Mount("/user", ch)
	})

	srv := http.Server{
		Addr:    "localhost:3001",
		Handler: r,
	}

	srv.ListenAndServe()
}
