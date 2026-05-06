package middleware

import (
	"github.com/edusyspro/edusys/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TenantMiddleware extracts the tenant_id from the authenticated user's JWT claims.
// It does NOT trust client-supplied headers like X-Tenant-ID, which would allow
// cross-tenant data access.
//
// This middleware should run AFTER AuthMiddleware.Authenticate so that
// tenant_id is already set in locals from the JWT. For unauthenticated routes,
// the tenant is resolved via other means (e.g., subdomain or API key).
func TenantMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// If tenant_id was already set by the auth middleware (from JWT claims), use it.
		if tenantID, ok := c.Locals("tenant_id").(uuid.UUID); ok {
			if tenantID != uuid.Nil {
				return c.Next()
			}
		}

		// For unauthenticated routes (e.g., public registration forms),
		// we could resolve tenant from subdomain or a trusted API key.
		// For now, reject requests without a valid tenant context.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tenant context required. Please authenticate first.",
		})
	}
}

// SuperAdminTenantOverride allows super_admins to specify a target tenant
// via the X-Tenant-ID header for cross-tenant administration.
// This should only be applied to super_admin-restricted routes.
func SuperAdminTenantOverride() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "super_admin" {
			return c.Next()
		}

		// Super admins can optionally override the tenant context
		overrideTenantID := c.Get("X-Tenant-ID")
		if overrideTenantID != "" {
			if id, err := uuid.Parse(overrideTenantID); err == nil {
				c.Locals("tenant_id", id)
			}
		}

		return c.Next()
	}
}
