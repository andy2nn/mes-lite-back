package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "mes-lite-back/docs" // swagger docs

	httpSwagger "github.com/swaggo/http-swagger"

	"mes-lite-back/internal/db"
	"mes-lite-back/internal/features/role"
	"mes-lite-back/internal/features/user"

	config "mes-lite-back/cmd/config"

	"mes-lite-back/pkg/logger"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init(slog.LevelDebug)

	dbConn, err := db.ConnectDB(&db.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	userRepo := user.NewGormRepository(dbConn)
	refreshRepo := user.NewRefreshTokenRepository(dbConn)
	roleRepo := role.NewGormRepository(dbConn)

	userService := user.NewService(userRepo)

	authService := user.NewAuthService(
		userRepo,
		refreshRepo,
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.TTL)*time.Second,
	)

	roleService := role.NewService(roleRepo)

	userHandler := user.NewHandler(userService)
	authHandler := user.NewAuthHandler(authService)
	roleHandler := role.NewHandler(roleService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/auth", func(r chi.Router) {
		r.Mount("/", authHandler.Routes())
	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Mount("/", userHandler.Routes())
	})

	apiRouter.Route("/roles", func(r chi.Router) {
		r.Mount("/", roleHandler.Routes())
	})

	r.Mount("/api/v1", apiRouter)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("MES Lite API v1.0"))
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 server started on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
