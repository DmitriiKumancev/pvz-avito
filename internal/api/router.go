package api

import (
	"net/http"
	"time"

	"github.com/dkumancev/avito-pvz/internal/api/v1/handlers"
	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	mux       *http.ServeMux
	db        *sqlx.DB
	jwtSecret []byte
}

func NewRouter(db *sqlx.DB, jwtSecret []byte) *Router {
	return &Router{
		mux:       http.NewServeMux(),
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (r *Router) Setup() http.Handler {
	// health check endpoint
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "db": "connected"}`))
	})

	userRepo := postgres.NewUserRepository(r.db)

	userService := services.NewUserService(userRepo, r.jwtSecret, 24*time.Hour)

	userHandler := handlers.NewUserHandler(userService)

	// Register routes
	r.mux.HandleFunc("/register", userHandler.Register)
	r.mux.HandleFunc("/login", userHandler.Login)
	r.mux.HandleFunc("/dummyLogin", userHandler.DummyLogin)

	// Protected routes (добавить позже)
	// For example:
	// r.mux.Handle("/pvz", middleware.AuthMiddleware(r.jwtSecret)(
	//     middleware.RoleMiddleware(domain.ModeratorRole)(
	//         http.HandlerFunc(pvzHandler.CreatePVZ)))))

	return r.mux
}
