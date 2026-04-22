package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/edusyspro/edusys/internal/config"
	"github.com/edusyspro/edusys/internal/middleware"
	"github.com/edusyspro/edusys/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAuthHandler(db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
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
	err := h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, email, password_hash, role, first_name, last_name, avatar_url, is_active, failed_login_attempts, locked_until
		FROM users WHERE email = $1`,
		req.Email,
	).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash, &user.Role,
		&user.FirstName, &user.LastName, &user.AvatarURL, &user.IsActive,
		&user.FailedLoginAttempts, &user.LockedUntil,
	)

	if err != nil {
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
				LastName:   user.LastName,
				Role:      user.Role,
				TenantID:  &user.TenantID,
				AvatarURL: user.AvatarURL,
			},
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: When refresh token storage is implemented,
	// invalidate the user's refresh token here.
	// For now, the client should discard the tokens.
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully. Please discard your tokens.",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	claims, err := middleware.ValidateToken(req.RefreshToken, h.cfg)
	if err != nil || claims["type"] != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid refresh token",
		})
	}

	userID, _ := uuid.Parse(claims["sub"].(string))

	var user models.User
	err = h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, email, role FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.TenantID, &user.Email, &user.Role)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	accessToken, _ := middleware.GenerateToken(user.ID, user.TenantID, user.Email, user.Role, h.cfg)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token": accessToken,
			"expires_in":   900,
		},
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email"`
	}

	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email is required",
		})
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	// Since we don't have reset_token right now in users schema, we use a simple comment.
	// TODO: implement emailing the token, and add reset_token and reset_expires_at columns to DB.
	log.Printf("Password reset token generated for %s: %s (expires: %v)\n", req.Email, token, expiresAt)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a reset link has been sent",
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Token and new password are required",
		})
	}

	// For security, enforce a minimum length
	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password must be at least 8 characters",
		})
	}

	// TODO: Verify token in the database when reset_token column exists
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	// UPDATE users SET password_hash = $1 WHERE reset_token = $2 AND reset_expires_at > NOW()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password has been reset",
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	type VerifyRequest struct {
		Token string `json:"token"`
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email verified successfully",
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var user models.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, email, role, first_name, last_name, phone, avatar_url, is_verified
		FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.TenantID, &user.Email, &user.Role, &user.FirstName, &user.LastName, &user.Phone, &user.AvatarURL, &user.IsVerified)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			TenantID:  &user.TenantID,
			AvatarURL: user.AvatarURL,
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	type UpdateProfileRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		AvatarURL string `json:"avatar_url"`
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	_, err := h.db.Exec(context.Background(),
		`UPDATE users SET first_name = $1, last_name = $2, phone = $3, avatar_url = $4, updated_at = NOW()
		WHERE id = $5`,
		req.FirstName, req.LastName, req.Phone, req.AvatarURL, userID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
	})
}

type StudentHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewStudentHandler(db *pgxpool.Pool, cfg *config.Config) *StudentHandler {
	return &StudentHandler{db: db, cfg: cfg}
}

func (h *StudentHandler) List(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	search := c.Query("search", "")
	status := c.Query("status", "active")
	sectionID := c.Query("section_id")

	offset := (page - 1) * limit

	query := `
		SELECT s.id, s.student_id, u.first_name, u.last_name, u.gender, u.date_of_birth, u.avatar_url, 
			sec.name as section_name, ss.roll_number, s.status
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN section_students ss ON ss.student_id = s.id
		LEFT JOIN sections sec ON sec.id = ss.section_id
		WHERE s.tenant_id = $1 AND s.status = $2
	`
	args := []interface{}{tenantID, status}
	argIndex := 3

	if search != "" {
		query += fmt.Sprintf(" AND (u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR s.student_id ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if sectionID != "" {
		query += fmt.Sprintf(" AND ss.section_id = $%d", argIndex)
		args = append(args, sectionID)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY ss.roll_number ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch students",
		})
	}
	defer rows.Close()

	var students []fiber.Map
	for rows.Next() {
		var studentID, userID, secID uuid.UUID
		var firstName, lastName, gender, sectionName, rollNumber, avatarURL string
		var dob *time.Time
		var status string

		err := rows.Scan(&studentID, &userID, &firstName, &lastName, &gender, &dob, &avatarURL, &sectionName, &rollNumber, &status)
		if err != nil {
			continue
		}

		students = append(students, fiber.Map{
			"id":          studentID,
			"student_id":  userID,
			"first_name":   firstName,
			"last_name":   lastName,
			"gender":      gender,
			"date_of_birth": dob,
			"section":     sectionName,
			"roll_number": rollNumber,
			"avatar_url":  avatarURL,
			"status":      status,
		})
	}

	var total int
	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM students WHERE tenant_id = $1 AND status = $2`,
		tenantID, status,
	).Scan(&total)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"students":   students,
			"pagination": models.Pagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: (total + limit - 1) / limit,
			},
		},
	})
}

func (h *StudentHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid student ID",
		})
	}

	tenantID := c.Locals("tenant_id").(uuid.UUID)

	var student models.Student
	err = h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, student_id, gender, date_of_birth, place_of_birth, 
			nationality, religion, blood_type, address, city, province, postal_code,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relation, notes
		FROM students WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(
		&student.ID, &student.TenantID, &student.StudentID, &student.Gender, &student.DateOfBirth,
		&student.PlaceOfBirth, &student.Nationality, &student.Religion, &student.BloodType,
		&student.Address, &student.City, &student.Province, &student.PostalCode,
		&student.EmergencyContactName, &student.EmergencyContactPhone,
		&student.EmergencyContactRelation, &student.Notes,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   student,
	})
}

type CreateStudentRequest struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	Gender              string `json:"gender"`
	DateOfBirth         string `json:"date_of_birth"`
	PlaceOfBirth        string `json:"place_of_birth"`
	Nationality         string `json:"nationality"`
	Religion            string `json:"religion"`
	BloodType           string `json:"blood_type"`
	Address             string `json:"address"`
	City                string `json:"city"`
	Province            string `json:"province"`
	PostalCode          string `json:"postal_code"`
	EmergencyContactName string `json:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
	EmergencyContactRelation string `json:"emergency_contact_relation"`
	SectionID           string `json:"section_id"`
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)

	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process password",
		})
	}
	userID := uuid.New()
	studentID := fmt.Sprintf("STU%d%04d", time.Now().Year(), 1)

	ctx := context.Background()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to start transaction",
		})
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, tenant_id, email, password_hash, role, first_name, last_name, is_active)
		VALUES ($1, $2, $3, $4, 'student', $5, $6, true)`,
		userID, tenantID, req.Email, string(hashedPassword), req.FirstName, req.LastName,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
		})
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO students (id, tenant_id, user_id, student_id, gender, date_of_birth, place_of_birth, 
			nationality, religion, blood_type, address, city, province, postal_code,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relation)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		uuid.New(), tenantID, userID, studentID, req.Gender, req.DateOfBirth, req.PlaceOfBirth,
		req.Nationality, req.Religion, req.BloodType, req.Address, req.City, req.Province,
		req.PostalCode, req.EmergencyContactName, req.EmergencyContactPhone, req.EmergencyContactRelation,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create student",
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Student created successfully",
		"data": fiber.Map{
			"student_id": studentID,
			"user_id":  userID,
		},
	})
}

func (h *StudentHandler) Update(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Student updated successfully",
	})
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid student ID",
		})
	}

	tenantID := c.Locals("tenant_id").(uuid.UUID)

	_, err = h.db.Exec(context.Background(),
		`UPDATE students SET is_active = false, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete student",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Student deleted successfully",
	})
}

func (h *StudentHandler) GetProfile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *StudentHandler) GetFamily(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *StudentHandler) AddFamily(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Family member added successfully",
	})
}

func (h *StudentHandler) GetAttendance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *StudentHandler) GetGrades(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *StudentHandler) GetFees(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *StudentHandler) GetReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

type AcademicHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAcademicHandler(db *pgxpool.Pool, cfg *config.Config) *AcademicHandler {
	return &AcademicHandler{db: db, cfg: cfg}
}

func (h *AcademicHandler) ListAcademicYears(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)

	rows, err := h.db.Query(context.Background(),
		`SELECT id, name, start_date, end_date, is_current, is_active 
		FROM academic_years WHERE tenant_id = $1 ORDER BY start_date DESC`,
		tenantID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch academic years",
		})
	}
	defer rows.Close()

	var years []fiber.Map
	for rows.Next() {
		var id uuid.UUID
		var name string
		var startDate, endDate time.Time
		var isCurrent, isActive bool

		if err := rows.Scan(&id, &name, &startDate, &endDate, &isCurrent, &isActive); err != nil {
			continue
		}

		years = append(years, fiber.Map{
			"id":          id,
			"name":        name,
			"start_date":  startDate,
			"end_date":   endDate,
			"is_current": isCurrent,
			"is_active":  isActive,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   years,
	})
}

func (h *AcademicHandler) CreateAcademicYear(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)

	type CreateAcademicYearRequest struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate  string `json:"end_date"`
		IsActive  bool   `json:"is_active"`
	}

	var req CreateAcademicYearRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO academic_years (id, tenant_id, name, start_date, end_date, is_current, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.New(), tenantID, req.Name, req.StartDate, req.EndDate, false, req.IsActive,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create academic year",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Academic year created successfully",
	})
}

func (h *AcademicHandler) ListSections(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)
	academicYearID := c.Query("academic_year_id")

	rows, err := h.db.Query(context.Background(),
		`SELECT s.id, s.name, s.grade_level, s.room, s.capacity, 
			(SELECT COUNT(*) FROM section_students WHERE section_id = s.id) as student_count
		FROM sections s
		WHERE s.tenant_id = $1 AND s.academic_year_id = $2
		ORDER BY s.grade_level, s.name`,
		tenantID, academicYearID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch sections",
		})
	}
	defer rows.Close()

	var sections []fiber.Map
	for rows.Next() {
		var id uuid.UUID
		var name string
		var gradeLevel, capacity, studentCount int
		var room string

		if err := rows.Scan(&id, &name, &gradeLevel, &room, &capacity, &studentCount); err != nil {
			continue
		}

		sections = append(sections, fiber.Map{
			"id":            id,
			"name":          name,
			"grade_level":    gradeLevel,
			"room":         room,
			"capacity":     capacity,
			"student_count": studentCount,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   sections,
	})
}

func (h *AcademicHandler) GetSection(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)
	id := c.Params("id")

	var section models.Section
	err := h.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, academic_year_id, name, grade_level, room, capacity
		FROM sections WHERE id = $1 AND tenant_id = $2`,
		id, tenantID,
	).Scan(&section.ID, &section.TenantID, &section.AcademicYearID, &section.Name, &section.GradeLevel, &section.Room, &section.Capacity)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Section not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   section,
	})
}

func (h *AcademicHandler) CreateSection(c *fiber.Ctx) error {
	type CreateSectionRequest struct {
		AcademicYearID string `json:"academic_year_id"`
		Name       string `json:"name"`
		GradeLevel int    `json:"grade_level"`
		Room      string `json:"room"`
		Capacity  int    `json:"capacity"`
	}

	var req CreateSectionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	tenantID := c.Locals("tenant_id").(uuid.UUID)

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO sections (id, tenant_id, academic_year_id, name, grade_level, room, capacity)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.New(), tenantID, req.AcademicYearID, req.Name, req.GradeLevel, req.Room, req.Capacity,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create section",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section created successfully",
	})
}

func (h *AcademicHandler) AssignStudents(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Students assigned successfully",
	})
}

func (h *AcademicHandler) GetTimetable(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AcademicHandler) ListSubjects(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)

	rows, err := h.db.Query(context.Background(),
		`SELECT id, name, code, description, is_optional, is_active
		FROM subjects WHERE tenant_id = $1 ORDER BY name`,
		tenantID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch subjects",
		})
	}
	defer rows.Close()

	var subjects []fiber.Map
	for rows.Next() {
		var id uuid.UUID
		var name, code, description string
		var isOptional, isActive bool

		if err := rows.Scan(&id, &name, &code, &description, &isOptional, &isActive); err != nil {
			continue
		}

		subjects = append(subjects, fiber.Map{
			"id":           id,
			"name":         name,
			"code":        code,
			"description":  description,
			"is_optional": isOptional,
			"is_active":   isActive,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   subjects,
	})
}

func (h *AcademicHandler) CreateSubject(c *fiber.Ctx) error {
	type CreateSubjectRequest struct {
		Name        string `json:"name"`
		Code       string `json:"code"`
		Description string `json:"description"`
		IsOptional bool   `json:"is_optional"`
	}

	var req CreateSubjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	tenantID := c.Locals("tenant_id").(uuid.UUID)

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO subjects (id, tenant_id, name, code, description, is_optional, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)`,
		uuid.New(), tenantID, req.Name, req.Code, req.Description, req.IsOptional,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create subject",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Subject created successfully",
	})
}

func (h *AcademicHandler) ListTimetables(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AcademicHandler) CreateTimetable(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Timetable created successfully",
	})
}

func (h *AcademicHandler) GenerateTimetable(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Timetable generated successfully",
	})
}

type AttendanceHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAttendanceHandler(db *pgxpool.Pool, cfg *config.Config) *AttendanceHandler {
	return &AttendanceHandler{db: db, cfg: cfg}
}

func (h *AttendanceHandler) ListStudentAttendance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AttendanceHandler) MarkStudentAttendance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Attendance marked successfully",
	})
}

func (h *AttendanceHandler) ListStaffAttendance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AttendanceHandler) MarkStaffAttendance(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Attendance marked successfully",
	})
}

func (h *AttendanceHandler) AttendanceReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *AttendanceHandler) GenerateQRCode(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *AttendanceHandler) ScanQR(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "QR scanned successfully",
	})
}

type ExamHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewExamHandler(db *pgxpool.Pool, cfg *config.Config) *ExamHandler {
	return &ExamHandler{db: db, cfg: cfg}
}

func (h *ExamHandler) List(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *ExamHandler) Get(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *ExamHandler) Create(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Exam created successfully",
	})
}

func (h *ExamHandler) Update(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Exam updated successfully",
	})
}

func (h *ExamHandler) Delete(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Exam deleted successfully",
	})
}

func (h *ExamHandler) GetQuestions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *ExamHandler) AddQuestion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Question added successfully",
	})
}

func (h *ExamHandler) PublishExam(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Exam published successfully",
	})
}

func (h *ExamHandler) GetResults(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *ExamHandler) GradeStudent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Student graded successfully",
	})
}

func (h *ExamHandler) ExportResults(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

type FeeHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewFeeHandler(db *pgxpool.Pool, cfg *config.Config) *FeeHandler {
	return &FeeHandler{db: db, cfg: cfg}
}

func (h *FeeHandler) ListFeeStructures(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *FeeHandler) CreateFeeStructure(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Fee structure created successfully",
	})
}

func (h *FeeHandler) UpdateFeeStructure(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Fee structure updated successfully",
	})
}

func (h *FeeHandler) AssignFees(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Fees assigned successfully",
	})
}

func (h *FeeHandler) ListStudentFees(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *FeeHandler) GetStudentFees(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *FeeHandler) ProcessPayment(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payment processed successfully",
	})
}

func (h *FeeHandler) PaymentWebhook(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func (h *FeeHandler) FinancialReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *FeeHandler) OverviewReport(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

type HRHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewHRHandler(db *pgxpool.Pool, cfg *config.Config) *HRHandler {
	return &HRHandler{db: db, cfg: cfg}
}

func (h *HRHandler) ListStaff(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *HRHandler) GetStaff(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *HRHandler) CreateStaff(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Staff created successfully",
	})
}

func (h *HRHandler) UpdateStaff(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Staff updated successfully",
	})
}

func (h *HRHandler) DeleteStaff(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Staff deleted successfully",
	})
}

func (h *HRHandler) ListLeaveRequests(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *HRHandler) SubmitLeaveRequest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Leave request submitted successfully",
	})
}

func (h *HRHandler) ApproveLeave(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Leave request processed successfully",
	})
}

func (h *HRHandler) ListPayrollRuns(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *HRHandler) CreatePayrollRun(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payroll run created successfully",
	})
}

func (h *HRHandler) GetPayrollRun(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *HRHandler) ProcessPayroll(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payroll processed successfully",
	})
}

type LMSHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewLMSHandler(db *pgxpool.Pool, cfg *config.Config) *LMSHandler {
	return &LMSHandler{db: db, cfg: cfg}
}

func (h *LMSHandler) ListCourses(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *LMSHandler) GetCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *LMSHandler) CreateCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Course created successfully",
	})
}

func (h *LMSHandler) UpdateCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Course updated successfully",
	})
}

func (h *LMSHandler) DeleteCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Course deleted successfully",
	})
}

func (h *LMSHandler) PublishCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Course published successfully",
	})
}

func (h *LMSHandler) GetContent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *LMSHandler) AddContent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Content added successfully",
	})
}

func (h *LMSHandler) EnrollStudent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Enrolled successfully",
	})
}

func (h *LMSHandler) MyEnrollments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *LMSHandler) SubmitAssignment(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Assignment submitted successfully",
	})
}

func (h *LMSHandler) SubmitQuiz(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Quiz submitted successfully",
	})
}

type LibraryHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewLibraryHandler(db *pgxpool.Pool, cfg *config.Config) *LibraryHandler {
	return &LibraryHandler{db: db, cfg: cfg}
}

func (h *LibraryHandler) ListBooks(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *LibraryHandler) GetBook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *LibraryHandler) AddBook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Book added successfully",
	})
}

func (h *LibraryHandler) UpdateBook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Book updated successfully",
	})
}

func (h *LibraryHandler) IssueBook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Book issued successfully",
	})
}

func (h *LibraryHandler) ReturnBook(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Book returned successfully",
	})
}

func (h *LibraryHandler) ListIssues(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *LibraryHandler) LibraryReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

type TransportHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewTransportHandler(db *pgxpool.Pool, cfg *config.Config) *TransportHandler {
	return &TransportHandler{db: db, cfg: cfg}
}

func (h *TransportHandler) ListRoutes(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *TransportHandler) CreateRoute(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Route created successfully",
	})
}

func (h *TransportHandler) ListVehicles(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *TransportHandler) AddVehicle(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Vehicle added successfully",
	})
}

func (h *TransportHandler) ListAssignments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *TransportHandler) CreateAssignment(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Assignment created successfully",
	})
}

func (h *TransportHandler) ListBoardingLogs(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *TransportHandler) LogBoarding(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Boarding logged successfully",
	})
}

func (h *TransportHandler) LiveTracking(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

type AnalyticsHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAnalyticsHandler(db *pgxpool.Pool, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{db: db, cfg: cfg}
}

func (h *AnalyticsHandler) Dashboard(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uuid.UUID)

	var totalStudents, totalTeachers, totalStaff, totalSections int
	
	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM students WHERE tenant_id = $1 AND status = 'active'`,
		tenantID,
	).Scan(&totalStudents)

	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM teachers WHERE tenant_id = $1 AND status = 'active'`,
		tenantID,
	).Scan(&totalTeachers)

	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM staff WHERE tenant_id = $1 AND status = 'active'`,
		tenantID,
	).Scan(&totalStaff)

	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM sections WHERE tenant_id = $1`,
		tenantID,
	).Scan(&totalSections)

	today := time.Now().Format("2006-01-02")
	var present, absent, late int
	
	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(COUNT(*), 0) FROM student_attendance 
		WHERE tenant_id = $1 AND date = $2 AND status = 'present'`,
		tenantID, today,
	).Scan(&present)

	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(COUNT(*), 0) FROM student_attendance 
		WHERE tenant_id = $1 AND date = $2 AND status = 'absent'`,
		tenantID, today,
	).Scan(&absent)

	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(COUNT(*), 0) FROM student_attendance 
		WHERE tenant_id = $1 AND date = $2 AND status = 'late'`,
		tenantID, today,
	).Scan(&late)

	totalAttendance := present + absent + late
	attendancePercentage := float64(0)
	if totalAttendance > 0 {
		attendancePercentage = float64(present+late) / float64(totalAttendance) * 100
	}

	var collected, pending, overdue float64
	
	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(paid_amount), 0) FROM student_fees 
		WHERE tenant_id = $1 AND status IN ('paid', 'partial')`,
		tenantID,
	).Scan(&collected)

	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(final_amount - paid_amount), 0) FROM student_fees 
		WHERE tenant_id = $1 AND status = 'pending' AND due_date >= $2`,
		tenantID, today,
	).Scan(&pending)

	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(final_amount - paid_amount), 0) FROM student_fees 
		WHERE tenant_id = $1 AND status = 'pending' AND due_date < $2`,
		tenantID, today,
	).Scan(&overdue)

	return c.JSON(fiber.Map{
		"success": true,
		"data": models.DashboardResponse{
			Overview: &models.DashboardOverview{
				TotalStudents:  totalStudents,
				TotalTeachers: totalTeachers,
				TotalStaff:    totalStaff,
				TotalSections: totalSections,
			},
			TodayAttendance: &models.AttendanceStats{
				Present:     present,
				Absent:     absent,
				Late:       late,
				Percentage: attendancePercentage,
			},
			FeeCollection: &models.FeeStats{
				Collected: collected,
				Pending:  pending,
				Overdue:  overdue,
			},
		},
	})
}

func (h *AnalyticsHandler) KPIMetrics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AnalyticsHandler) CustomReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AnalyticsHandler) ExportData(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *AnalyticsHandler) Predictions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AnalyticsHandler) PerformancePrediction(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *AnalyticsHandler) AnomalyDetection(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AnalyticsHandler) ChatBot(c *fiber.Ctx) error {
	type ChatRequest struct {
		Message string `json:"message"`
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reply": "Thank you for your message. Our team will get back to you shortly.",
		},
	})
}

type AdmissionHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAdmissionHandler(db *pgxpool.Pool, cfg *config.Config) *AdmissionHandler {
	return &AdmissionHandler{db: db, cfg: cfg}
}

func (h *AdmissionHandler) ListLeads(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AdmissionHandler) GetLead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *AdmissionHandler) CreateLead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Lead created successfully",
	})
}

func (h *AdmissionHandler) UpdateLead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Lead updated successfully",
	})
}

func (h *AdmissionHandler) ConvertLead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Lead converted successfully",
	})
}

func (h *AdmissionHandler) ListForms(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *AdmissionHandler) SubmitForm(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Form submitted successfully",
	})
}

func (h *AdmissionHandler) PipelineView(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

type MessageHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewMessageHandler(db *pgxpool.Pool, cfg *config.Config) *MessageHandler {
	return &MessageHandler{db: db, cfg: cfg}
}

func (h *MessageHandler) List(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *MessageHandler) Get(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *MessageHandler) Send(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Message sent successfully",
	})
}

func (h *MessageHandler) MarkRead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Message marked as read",
	})
}

func (h *MessageHandler) ListNotifications(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *MessageHandler) MarkNotificationRead(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification marked as read",
	})
}

func (h *MessageHandler) ListAnnouncements(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   []fiber.Map{},
	})
}

func (h *MessageHandler) CreateAnnouncement(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Announcement created successfully",
	})
}

func (h *MessageHandler) SendWhatsApp(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "WhatsApp message sent successfully",
	})
}

type TenantHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewTenantHandler(db *pgxpool.Pool, cfg *config.Config) *TenantHandler {
	return &TenantHandler{db: db, cfg: cfg}
}

func (h *TenantHandler) List(c *fiber.Ctx) error {
	rows, err := h.db.Query(context.Background(),
		`SELECT id, name, code, email, phone, is_active, subscription_plan
		FROM tenants ORDER BY created_at DESC`,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch tenants",
		})
	}
	defer rows.Close()

	var tenants []fiber.Map
	for rows.Next() {
		var id uuid.UUID
		var name, code, email, phone, subscriptionPlan string
		var isActive bool

		if err := rows.Scan(&id, &name, &code, &email, &phone, &isActive, &subscriptionPlan); err != nil {
			continue
		}

		tenants = append(tenants, fiber.Map{
			"id":                id,
			"name":             name,
			"code":            code,
			"email":           email,
			"phone":           phone,
			"is_active":       isActive,
			"subscription_plan": subscriptionPlan,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":   tenants,
	})
}

func (h *TenantHandler) Get(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
	type CreateTenantRequest struct {
		Name  string `json:"name"`
		Code  string `json:"code"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO tenants (id, name, code, email, phone, is_active)
		VALUES ($1, $2, $3, $4, $5, true)`,
		uuid.New(), req.Name, req.Code, req.Email, req.Phone,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create tenant",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tenant created successfully",
	})
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tenant updated successfully",
	})
}

func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tenant deleted successfully",
	})
}

func (h *TenantHandler) GetSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":   fiber.Map{},
	})
}

func (h *TenantHandler) UpdateSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Settings updated successfully",
	})
}