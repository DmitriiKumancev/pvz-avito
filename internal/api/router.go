package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/dkumancev/avito-pvz/internal/api/middleware"
	"github.com/dkumancev/avito-pvz/internal/api/v1/handlers"
	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/logger"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/product"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/pvz"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/reception"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/user"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	router    *mux.Router
	db        *sqlx.DB
	jwtSecret []byte
	logger    *slog.Logger
}

func NewRouter(db *sqlx.DB, jwtSecret []byte) *Router {
	apiLogger := logger.NewLogger(logger.Config{
		Level:  logger.LevelInfo,
		Format: "text",
	})

	return &Router{
		router:    mux.NewRouter(),
		db:        db,
		jwtSecret: jwtSecret,
		logger:    apiLogger,
	}
}

func (r *Router) Setup() http.Handler {
	// Репозитории
	userRepo := user.New(r.db)
	pvzRepo := pvz.New(r.db)
	receptionRepo := reception.New(r.db)
	productRepo := product.New(r.db)

	// Сервисы
	userService := services.NewUserService(userRepo, r.jwtSecret, 24*time.Hour)
	pvzService := services.NewPVZService(pvzRepo)
	receptionService := services.NewReceptionService(pvzRepo, receptionRepo, productRepo)

	// Хендлеры
	userHandler := handlers.NewUserHandler(userService)
	pvzHandler := handlers.NewPVZHandler(pvzService, receptionService)
	receptionHandler := handlers.NewReceptionHandler(receptionService)
	productHandler := handlers.NewProductHandler(receptionService)

	// Глобальные middleware
	r.router.Use(middleware.RecoveryMiddleware(r.logger))
	r.router.Use(middleware.LoggerMiddleware(r.logger))

	// Health check
	r.router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "db": "connected"}`))
	}).Methods(http.MethodGet)

	// Публичные маршруты
	r.router.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	r.router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	r.router.HandleFunc("/dummyLogin", userHandler.DummyLogin).Methods(http.MethodPost)

	// Защищенные маршруты

	// ПВЗ - модератор может создавать, все авторизованные могут просматривать
	r.router.Handle("/pvz", middleware.AuthMiddleware(r.jwtSecret,
		middleware.RoleMiddleware([]domain.UserRole{domain.ModeratorRole},
			http.HandlerFunc(pvzHandler.CreatePVZ)))).Methods(http.MethodPost)

	r.router.Handle("/pvz", middleware.AuthMiddleware(r.jwtSecret,
		http.HandlerFunc(pvzHandler.ListPVZ))).Methods(http.MethodGet)

	// Закрытие приемки - только сотрудник
	r.router.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(r.jwtSecret,
		middleware.RoleMiddleware([]domain.UserRole{domain.EmployeeRole},
			http.HandlerFunc(pvzHandler.CloseLastReception)))).Methods(http.MethodPost)

	// Удаление последнего товара - только сотрудник
	r.router.Handle("/pvz/{pvzId}/delete_last_product", middleware.AuthMiddleware(r.jwtSecret,
		middleware.RoleMiddleware([]domain.UserRole{domain.EmployeeRole},
			http.HandlerFunc(pvzHandler.DeleteLastProduct)))).Methods(http.MethodPost)

	// Создание приемки - только сотрудник
	r.router.Handle("/receptions", middleware.AuthMiddleware(r.jwtSecret,
		middleware.RoleMiddleware([]domain.UserRole{domain.EmployeeRole},
			http.HandlerFunc(receptionHandler.CreateReception)))).Methods(http.MethodPost)

	// Добавление товара - только сотрудник
	r.router.Handle("/products", middleware.AuthMiddleware(r.jwtSecret,
		middleware.RoleMiddleware([]domain.UserRole{domain.EmployeeRole},
			http.HandlerFunc(productHandler.AddProduct)))).Methods(http.MethodPost)

	r.logger.Info("API маршрутизатор настроен", "routes_count", muxRoutesCount(r.router))
	return r.router
}

// количество зарегистрированных маршрутов
func muxRoutesCount(router *mux.Router) int {
	count := 0
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		count++
		return nil
	})
	return count
}
