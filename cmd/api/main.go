package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mohdrashid9678/tirush/config"
	"github.com/mohdrashid9678/tirush/internal/database"
	"github.com/mohdrashid9678/tirush/internal/handlers"
	"github.com/mohdrashid9678/tirush/internal/repository"
	"github.com/mohdrashid9678/tirush/internal/service"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Database Connection
	dbService, err := database.New(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbService.Close()

	// 3. Dependency Injection, Service and Handler Setup
	repo := repository.NewRepository(dbService.Db)
	svc := service.NewService(repo)
	handler := handlers.NewHandler(svc)

	// 4. Router Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// 5. Register Routes
	handler.RegisterRoutes(r)

	// 6. Start Server
	log.Printf("tirush is running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
