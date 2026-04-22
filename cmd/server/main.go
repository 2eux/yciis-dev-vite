package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/edusyspro/edusys/internal/config"
	"github.com/edusyspro/edusys/internal/database"
	"github.com/edusyspro/edusys/internal/middleware"
	"github.com/edusyspro/edusys/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg := config.Load()

	// Validate configuration before proceeding
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Edusys Pro",
		// Do NOT expose server software/version in production
		ServerHeader: "",
		// Sanitize error messages — never expose internal details to clients
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal server error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			// In production, never leak raw error messages
			if cfg.IsProduction() && code == fiber.StatusInternalServerError {
				message = "An unexpected error occurred"
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": message,
			})
		},
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		BodyLimit:         4 * 1024 * 1024, // 4MB max request body
		DisableKeepalive:  false,
		EnablePrintRoutes: cfg.Debug,
	})

	// ─── Global Middleware (order matters) ───────────────────────────

	// 1. Panic recovery
	app.Use(recover.New())

	// 2. Request ID for tracing
	app.Use(requestid.New())

	// 3. Security headers (HSTS, X-Frame-Options, X-Content-Type-Options, etc.)
	app.Use(helmet.New())

	// 4. CORS — configured per environment
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.AllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// 5. Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:               cfg.RateLimitRequests,
		Expiration:        time.Duration(cfg.RateLimitWindowMs) * time.Millisecond,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	}))

	// 6. Request logging
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${reqHeader:X-Request-Id}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// 7. Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// ─── Static files ───────────────────────────────────────────────
	app.Static("/public", "./web/public")

	// ─── Health check (no auth required) ────────────────────────────
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	// Root info endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "Edusys Pro API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// ─── API Routes ─────────────────────────────────────────────────
	api := app.Group("/api/v1")

	// Tenant middleware now trusts JWT tenant_id, not client headers
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

	// 404 handler for unknown API routes
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Endpoint not found",
		})
	})

	// ─── Graceful Shutdown ──────────────────────────────────────────
	go func() {
		addr := ":" + cfg.ServerPort
		log.Printf("Server starting on %s (env=%s, debug=%v)", addr, cfg.Environment, cfg.Debug)
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