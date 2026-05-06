package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/edusyspro/edusys/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthMiddleware struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// NewAuthMiddlewareWithDB creates an auth middleware with database access
// for permission checking against the role_permissions table.
func NewAuthMiddlewareWithDB(cfg *config.Config, db *pgxpool.Pool) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg, db: db}
}

// Authenticate validates the JWT from the Authorization header and
// extracts user identity into fiber locals.
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Authorization header required",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid authorization format. Use: Bearer <token>",
		})
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC signing method to prevent algorithm confusion attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token claims",
		})
	}

	// Safely extract claims with type checking — no panics on malformed tokens
	subStr, ok := claims["sub"].(string)
	if !ok || subStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token missing subject claim",
		})
	}
	userID, err := uuid.Parse(subStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID in token",
		})
	}

	tenantStr, ok := claims["tenant_id"].(string)
	if !ok || tenantStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token missing tenant claim",
		})
	}
	tenantID, err := uuid.Parse(tenantStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid tenant ID in token",
		})
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token missing role claim",
		})
	}

	// Reject refresh tokens used as access tokens
	if tokenType, exists := claims["type"].(string); exists && tokenType == "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Refresh tokens cannot be used for API access",
		})
	}

	c.Locals("user_id", userID)
	c.Locals("tenant_id", tenantID)
	c.Locals("role", role)

	return c.Next()
}

// RequireRole checks if the authenticated user has one of the allowed roles.
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Unable to determine user role",
			})
		}
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Insufficient permissions",
		})
	}
}

// RequirePermission checks if the user's role has the specified module/action
// permission in the role_permissions table. Super admins bypass this check.
func (m *AuthMiddleware) RequirePermission(module, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Unable to determine user role",
			})
		}

		// Super admins always have full access
		if role == "super_admin" {
			return c.Next()
		}

		// Check against the role_permissions table
		if m.db == nil {
			// Fallback: if no DB connection, only allow admin+
			if role == "admin" {
				return c.Next()
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions",
			})
		}

		var count int
		err := m.db.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM role_permissions rp
			 JOIN permissions p ON p.id = rp.permission_id
			 WHERE rp.role = $1 AND p.module = $2 AND p.action = $3`,
			role, module, action,
		).Scan(&count)

		if err != nil || count == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions for this action",
			})
		}

		return c.Next()
	}
}

// GenerateToken creates a short-lived access JWT with user identity claims.
func GenerateToken(userID, tenantID uuid.UUID, email, role string, cfg *config.Config) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":       userID.String(),
		"tenant_id": tenantID.String(),
		"email":     email,
		"role":      role,
		"type":      "access",
		"exp":       now.Add(cfg.JWTExpiry).Unix(),
		"iat":       now.Unix(),
		"nbf":       now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// GenerateRefreshToken creates a long-lived refresh token containing only the user ID.
func GenerateRefreshToken(userID uuid.UUID, cfg *config.Config) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"type": "refresh",
		"exp":  now.Add(cfg.RefreshExpiry).Unix(),
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateToken parses and validates a JWT string, returning the claims.
func ValidateToken(tokenString string, cfg *config.Config) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
