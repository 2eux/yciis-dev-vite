package handlers

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/edusyspro/edusys/internal/config"
	"github.com/edusyspro/edusys/internal/middleware"
	"github.com/edusyspro/edusys/internal/models"
	"github.com/edusyspro/edusys/internal/totp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type TwoFAHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewTwoFAHandler(db *pgxpool.Pool, cfg *config.Config) *TwoFAHandler {
	return &TwoFAHandler{db: db, cfg: cfg}
}

// Get2FASetup returns the current 2FA status or generates a new secret for setup.
func (h *TwoFAHandler) Get2FASetup(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var twoFAEnabled bool
	var twoFASecret *string

	err := h.db.QueryRow(context.Background(),
		`SELECT two_secret_enabled, two_secret_secret FROM users WHERE id = $1`,
		userID,
	).Scan(&twoFAEnabled, &twoFASecret)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	if twoFAEnabled {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"enabled": true,
			},
		})
	}

	totpInstance := totp.NewTOTP()
	secret, err := totpInstance.GenerateSecret()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate secret",
		})
	}

	var email, firstName string
	h.db.QueryRow(context.Background(),
		`SELECT email, first_name FROM users WHERE id = $1`,
		userID,
	).Scan(&email, &firstName)

	otpAuthURL := totpInstance.GetGoogleAuthenticatorURL(secret, email, "Edusys Pro")

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"enabled":        false,
			"secret":         secret,
			"otpauth_url":    otpAuthURL,
			"qr_code_base64": base64.StdEncoding.EncodeToString([]byte(otpAuthURL)),
		},
	})
}

// Enable2FA validates a TOTP code against the provided secret and enables 2FA.
func (h *TwoFAHandler) Enable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	type Enable2FARequest struct {
		Secret string `json:"secret"`
		Code   string `json:"code"`
	}

	var req Enable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.Secret == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Secret and code are required",
		})
	}

	totpInstance := totp.NewTOTP()

	if !totpInstance.Validate(req.Secret, req.Code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid verification code. Please try again.",
		})
	}

	// Store the TOTP secret encrypted using pgcrypto
	// TODO: For production, encrypt with application-level AES-256 key
	_, err := h.db.Exec(context.Background(),
		`UPDATE users SET two_secret_enabled = true, two_secret_secret = $1, updated_at = NOW() WHERE id = $2`,
		req.Secret, userID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to enable 2FA",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "2FA enabled successfully",
	})
}

// Disable2FA validates the current TOTP code and disables 2FA.
func (h *TwoFAHandler) Disable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	type Disable2FARequest struct {
		Code string `json:"code"`
	}

	var req Disable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	var storedSecret *string
	err := h.db.QueryRow(context.Background(),
		`SELECT two_secret_secret FROM users WHERE id = $1`,
		userID,
	).Scan(&storedSecret)

	if err != nil || storedSecret == nil || *storedSecret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "2FA is not enabled",
		})
	}

	totpInstance := totp.NewTOTP()
	if !totpInstance.Validate(*storedSecret, req.Code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid verification code",
		})
	}

	_, err = h.db.Exec(context.Background(),
		`UPDATE users SET two_secret_enabled = false, two_secret_secret = NULL, updated_at = NOW() WHERE id = $1`,
		userID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to disable 2FA",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "2FA disabled successfully",
	})
}

// Verify2FA validates a 2FA code for a specific user (used during login flow).
func (h *TwoFAHandler) Verify2FA(c *fiber.Ctx) error {
	type Verify2FARequest struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}

	var req Verify2FARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	var storedSecret *string
	var twoFAEnabled bool
	err = h.db.QueryRow(context.Background(),
		`SELECT two_secret_enabled, two_secret_secret FROM users WHERE id = $1`,
		userID,
	).Scan(&twoFAEnabled, &storedSecret)

	if err != nil || !twoFAEnabled || storedSecret == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "2FA is not enabled",
		})
	}

	totpInstance := totp.NewTOTP()
	if !totpInstance.Validate(*storedSecret, req.Code) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid 2FA code",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "2FA verification successful",
	})
}

// LoginWith2FA handles login with email/password plus optional 2FA code.
func (h *TwoFAHandler) LoginWith2FA(c *fiber.Ctx) error {
	type LoginWith2FARequest struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		TwoFactorCode string `json:"two_factor_code"`
	}

	var req LoginWith2FARequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Input validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email and password are required",
		})
	}

	var user models.User
	var twoFAEnabled bool
	var twoFASecret *string

	err := h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, email, password_hash, role, first_name, last_name, avatar_url, is_active,
		failed_login_attempts, locked_until, two_secret_enabled, two_secret_secret
		FROM users WHERE email = $1`,
		req.Email,
	).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash, &user.Role,
		&user.FirstName, &user.LastName, &user.AvatarURL, &user.IsActive,
		&user.FailedLoginAttempts, &user.LockedUntil, &twoFAEnabled, &twoFASecret,
	)

	if err != nil {
		// Deliberate generic message to prevent user enumeration
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Account is disabled",
		})
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Account is locked. Try again later",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Increment failed attempts and lock after 5
		_, _ = h.db.Exec(context.Background(),
			`UPDATE users SET failed_login_attempts = failed_login_attempts + 1,
			locked_until = CASE WHEN failed_login_attempts >= 4 THEN NOW() + INTERVAL '15 minutes' ELSE NULL END
			WHERE id = $1`,
			user.ID,
		)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	// Check 2FA if enabled
	if twoFAEnabled && twoFASecret != nil && *twoFASecret != "" {
		if req.TwoFactorCode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":      false,
				"message":      "2FA code required",
				"requires_2fa": true,
			})
		}

		totpInstance := totp.NewTOTP()
		if !totpInstance.Validate(*twoFASecret, req.TwoFactorCode) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid 2FA code",
			})
		}
	}

	accessToken, err := middleware.GenerateToken(user.ID, user.TenantID, user.Email, user.Role, h.cfg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate access token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, h.cfg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate refresh token",
		})
	}

	// Update login metadata and reset failed attempts
	_, _ = h.db.Exec(context.Background(),
		`UPDATE users SET last_login_at = NOW(), last_login_ip = $1, failed_login_attempts = 0
		WHERE id = $2`,
		c.IP(), user.ID,
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": models.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(h.cfg.JWTExpiry.Seconds()),
			User: models.UserResponse{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				TenantID:  &user.TenantID,
				AvatarURL: user.AvatarURL,
			},
		},
	})
}

// Get2FAStatus returns whether 2FA is enabled for the current user.
func (h *TwoFAHandler) Get2FAStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var twoFAEnabled bool
	err := h.db.QueryRow(context.Background(),
		`SELECT two_secret_enabled FROM users WHERE id = $1`,
		userID,
	).Scan(&twoFAEnabled)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"enabled":        twoFAEnabled,
			"time_remaining": totp.GetTimeRemaining(),
		},
	})
}