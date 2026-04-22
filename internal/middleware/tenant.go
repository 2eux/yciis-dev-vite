package middleware

import (
	"github.com/edusyspro/edusys/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TenantMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default"
		}

		id, err := uuid.Parse(tenantID)
		if err != nil {
			id = uuid.MustParse("00000000-0000-0000-0000-000000000000")
		}

		c.Locals("tenant_id", id)
		return c.Next()
	}
}