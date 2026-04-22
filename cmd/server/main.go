package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edusyspro/edusys/internal/config"
	"github.com/edusyspro/edusys/internal/database"
	"github.com/edusyspro/edusys/internal/middleware"
	"github.com/edusyspro/edusys/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/websocket/v2"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Edusys Pro",
		ServerHeader: "EdusysPro/1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) *fiber.Error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${method} | ${path} | ${latency} | ${cid}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.BestSpeed,
	}))

	app.Static("/public", "./web/public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "Edusys Pro API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	app.Get("/ws", func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrBadRequest
		}
		return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
			"message": "WebSocket upgrade required",
		})
	})

	api := app.Group("/api/v1")
	api.Use(middleware.TenantMiddleware(cfg))

	routes.RegisterAuthRoutes(api, db, cfg)
	routes.RegisterStudentRoutes(api, db, cfg)
	routes.RegisterAcademicRoutes(api, db, cfg)
	routes.RegisterAttendanceRoutes(api, db, cfg)
	routes.RegisterExamRoutes(api, db, cfg)
	routes.RegisterFeeRoutes(api, db, cfg)
	routes.RegisterHRRoutes(api, db, cfg)
	routes.RegisterLMSRoutes(api, db, cfg)
	routes.RegisterLibraryRoutes(api, db, cfg)
	routes.RegisterTransportRoutes(api, db, cfg)
	routes.RegisterAnalyticsRoutes(api, db, cfg)
	routes.RegisterAdmissionRoutes(api, db, cfg)
	routes.RegisterMessageRoutes(api, db, cfg)
	routes.RegisterTenantRoutes(api, db, cfg)

	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Endpoint not found",
		})
	})

	go func() {
		addr := ":" + cfg.Port
		log.Printf("Server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}