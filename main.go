package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hrms/admin"
	"hrms/auth"
	db "hrms/database"
	"hrms/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Configs, logger, database initialization
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	cfg := MustLoadConfig()
	log := SetupLogger(cfg)
	db := db.MustInitDB(cfg.DBUrl, log)

	// Access data layer
	userRepo := users.NewUserRepository(db, log)
	adminRepo := admin.NewAdminRepository(db, log)

	// App initialization and routes setup
	app := fiber.New(fiber.Config{AppName: "HRMS"})
	app.Use(cors.New())

	auth.SetupRoutes(app.Group("/auth"), userRepo, log, cfg.JWTSecret)

	userGroup := app.Group("/user")
	userGroup.Use(auth.Middleware(db, cfg.JWTSecret))
	users.SetupRoutes(userGroup, userRepo)

	adminGroup := app.Group("/admin")
	adminGroup.Use(admin.Middleware())
	admin.SetupRoutes(adminGroup, adminRepo)

	// Starting the server with graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Starting server...", "port", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Error("Failed to start server", "error", err.Error())
		}
	}()

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err.Error())
	}

	log.Info("Server exited gracefully")
}
