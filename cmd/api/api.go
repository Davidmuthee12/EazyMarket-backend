package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Davidmuthee12/eazymarket/internals/auth"
	"github.com/Davidmuthee12/eazymarket/internals/env"
	"github.com/Davidmuthee12/eazymarket/internals/mailer"
	"github.com/Davidmuthee12/eazymarket/internals/store"
	cache "github.com/Davidmuthee12/eazymarket/internals/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	auth        authConfig
	mail        mailConfig
	frontendURL string
	redisCfg    redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type dbConfig struct {
	addr          string
	maxOpenConn   int
	maxIddleConns int
	maxIddleTime  string
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Basic CORS
	// Be careful where you place the cors middleware. e.g. place before the RateLimiter.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.AuthTokenMiddleware).Get("/health", app.healthCheckHandler)
		docsURL := "/v1/swagger/doc.json"
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL),
			httpSwagger.UIConfig(map[string]string{
				"tagsSorter":       "\"alpha\"",
				"operationsSorter": "\"alpha\"",
			}),
		))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/users", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userUUID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
				r.Post("/upgrade-to-vendor", app.updateRoleHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getAllUsersHandlers)
			})
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Use(app.RequireRole("admin"))
			r.Get("/", app.getAllUsersHandlers)
			r.Get("/vendor-request", app.vendorRequestHandler)
			r.Put("/vendor-request/{userUUID}/approve", app.approveVendorHandler)
			r.Put("/vendor-request/{userUUID}/reject", app.rejectVendorHandler)
		})

		r.Route("/vendor", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Use(app.RequireRole("vendor"))
			r.Post("/profile", app.vendorProfileHandler)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Graceful shutdown

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("Signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.Infow("Server has stopped", "addr", app.config.env)

	return nil
}
