package main

import (
	"net/http"
	"os"

	"github.com/Quaqmre/circleramkit/customer"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	//Make Customer Service tree
	var s customer.Service
	{
		s = customer.NewCustomerService()
		s = customer.NewLoggingService(logger, s)
	}

	//Make Handler each endpoint
	ch := customer.MakeHandler(s)

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
	r.Use(cors.Handler)
	r.Mount("/v1", ch)

	srv := http.Server{
		Addr:    "localhost:3001",
		Handler: r,
	}

	srv.ListenAndServe()
}
