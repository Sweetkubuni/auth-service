package main

import (
	"fmt"
	"net/http"
	"os"

	"auth-service/login"
	"auth-service/register"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	router := router(logger)

	loggedRouter := hlog.NewHandler(logger)(router)

	if err := http.ListenAndServe(addr, loggedRouter); err != nil {
		logger.Fatal().Err(err).Msg("Server error")
	}
}

func router(logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			logger.Info().Msgf("Protected area. User ID: %v", claims["user_id"])
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/register", register.RegisterHandler)
		r.Post("/login", login.LoginHandler)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Msg("Welcome anonymous")
			w.Write([]byte("welcome anonymous"))
		})
	})

	return r
}
