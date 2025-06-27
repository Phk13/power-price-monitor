package main

import (
	"log"

	"power-price-monitor/config"
	"power-price-monitor/handler"
	"power-price-monitor/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	if cfg.Token == "" {
		log.Fatal("Error: TOKEN is required. Please set it in .env file or environment variable")
	}

	// Initialize services
	powerService := service.NewPowerPriceService(cfg)

	// Initialize handlers
	powerHandler := handler.NewPowerHandler(powerService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/health", powerHandler.HealthCheck)
	e.GET("/optimal-hours", powerHandler.GetOptimalHours)

	// Start server
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
