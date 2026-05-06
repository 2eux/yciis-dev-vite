package routes

import (
	"github.com/edusyspro/edusys/internal/config"
	"github.com/edusyspro/edusys/internal/handlers"
	"github.com/edusyspro/edusys/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewAuthHandler(db, cfg)
	twoFAHandler := handlers.NewTwoFAHandler(db, cfg)
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/login-2fa", twoFAHandler.LoginWith2FA)
	auth.Post("/logout", handler.Logout)
	auth.Post("/refresh", handler.RefreshToken)
	auth.Post("/forgot-password", handler.ForgotPassword)
	auth.Post("/reset-password", handler.ResetPassword)
	auth.Post("/verify-email", handler.VerifyEmail)
	auth.Get("/me", middleware.NewAuthMiddleware(cfg).Authenticate, handler.Me)
	auth.Put("/profile", middleware.NewAuthMiddleware(cfg).Authenticate, handler.UpdateProfile)

	twoFA := api.Group("/2fa", middleware.NewAuthMiddleware(cfg).Authenticate)
	twoFA.Get("/status", twoFAHandler.Get2FAStatus)
	twoFA.Get("/setup", twoFAHandler.Get2FASetup)
	twoFA.Post("/enable", twoFAHandler.Enable2FA)
	twoFA.Post("/disable", twoFAHandler.Disable2FA)
	twoFA.Post("/verify", twoFAHandler.Verify2FA)
}

func RegisterStudentRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewStudentHandler(db, cfg)
	students := api.Group("/students", middleware.NewAuthMiddleware(cfg).Authenticate)

	students.Get("", handler.List)
	students.Get("/:id", handler.Get)
	students.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.Create)
	students.Put("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.Update)
	students.Delete("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.Delete)
	students.Get("/:id/profile", handler.GetProfile)
	students.Get("/:id/family", handler.GetFamily)
	students.Post("/:id/family", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.AddFamily)
	students.Get("/:id/attendance", handler.GetAttendance)
	students.Get("/:id/grades", handler.GetGrades)
	students.Get("/:id/fees", handler.GetFees)
	students.Get("/:id/reports", handler.GetReports)
}

func RegisterAcademicRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewAcademicHandler(db, cfg)
	academic := api.Group("/academic", middleware.NewAuthMiddleware(cfg).Authenticate)

	academic.Get("/years", handler.ListAcademicYears)
	academic.Post("/years", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.CreateAcademicYear)

	sections := api.Group("/sections", middleware.NewAuthMiddleware(cfg).Authenticate)
	sections.Get("", handler.ListSections)
	sections.Get("/:id", handler.GetSection)
	sections.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.CreateSection)
	sections.Put("/:id/students", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.AssignStudents)
	sections.Get("/:id/timetable", handler.GetTimetable)

	subjects := api.Group("/subjects", middleware.NewAuthMiddleware(cfg).Authenticate)
	subjects.Get("", handler.ListSubjects)
	subjects.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.CreateSubject)

	timetable := api.Group("/timetables", middleware.NewAuthMiddleware(cfg).Authenticate)
	timetable.Get("", handler.ListTimetables)
	timetable.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.CreateTimetable)
	timetable.Get("/generate", handler.GenerateTimetable)
}

func RegisterAttendanceRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewAttendanceHandler(db, cfg)
	attendance := api.Group("/attendance", middleware.NewAuthMiddleware(cfg).Authenticate)

	attendance.Get("/students", handler.ListStudentAttendance)
	attendance.Post("/students", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.MarkStudentAttendance)
	attendance.Get("/staff", handler.ListStaffAttendance)
	attendance.Post("/staff", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.MarkStaffAttendance)
	attendance.Get("/reports", handler.AttendanceReports)
	attendance.Get("/qr/generate", handler.GenerateQRCode)
	attendance.Post("/qr/scan", handler.ScanQR)
}

func RegisterExamRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewExamHandler(db, cfg)
	exams := api.Group("/exams", middleware.NewAuthMiddleware(cfg).Authenticate)

	exams.Get("", handler.List)
	exams.Get("/:id", handler.Get)
	exams.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.Create)
	exams.Put("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.Update)
	exams.Delete("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.Delete)
	exams.Get("/:id/questions", handler.GetQuestions)
	exams.Post("/:id/questions", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.AddQuestion)
	exams.Put("/:id/publish", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.PublishExam)
	exams.Get("/:id/results", handler.GetResults)
	exams.Post("/:id/grade", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.GradeStudent)
	exams.Get("/:id/export", handler.ExportResults)
}

func RegisterFeeRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewFeeHandler(db, cfg)
	fees := api.Group("/fees", middleware.NewAuthMiddleware(cfg).Authenticate)

	fees.Get("/structures", handler.ListFeeStructures)
	fees.Post("/structures", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "finance"), handler.CreateFeeStructure)
	fees.Put("/structures/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "finance"), handler.UpdateFeeStructure)
	fees.Post("/assign", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "finance"), handler.AssignFees)
	fees.Get("/students", handler.ListStudentFees)
	fees.Get("/students/:id", handler.GetStudentFees)
	fees.Post("/students/:id/pay", handler.ProcessPayment)
	fees.Post("/webhook", handler.PaymentWebhook)
	fees.Get("/reports", handler.FinancialReports)
	fees.Get("/reports/overview", handler.OverviewReport)
}

func RegisterHRRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewHRHandler(db, cfg)
	hr := api.Group("/hr", middleware.NewAuthMiddleware(cfg).Authenticate)

	hr.Get("/staff", handler.ListStaff)
	hr.Get("/staff/:id", handler.GetStaff)
	hr.Post("/staff", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.CreateStaff)
	hr.Put("/staff/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.UpdateStaff)
	hr.Delete("/staff/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.DeleteStaff)

	hr.Get("/leave", handler.ListLeaveRequests)
	hr.Post("/leave", handler.SubmitLeaveRequest)
	hr.Put("/leave/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.ApproveLeave)

	hr.Get("/payroll", handler.ListPayrollRuns)
	hr.Post("/payroll", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.CreatePayrollRun)
	hr.Get("/payroll/:id", handler.GetPayrollRun)
	hr.Post("/payroll/:id/process", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "hr"), handler.ProcessPayroll)
}

func RegisterLMSRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewLMSHandler(db, cfg)
	lms := api.Group("/lms", middleware.NewAuthMiddleware(cfg).Authenticate)

	lms.Get("/courses", handler.ListCourses)
	lms.Get("/courses/:id", handler.GetCourse)
	lms.Post("/courses", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.CreateCourse)
	lms.Put("/courses/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.UpdateCourse)
	lms.Delete("/courses/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "super_admin"), handler.DeleteCourse)
	lms.Post("/courses/:id/publish", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.PublishCourse)
	lms.Get("/courses/:id/content", handler.GetContent)
	lms.Post("/courses/:id/content", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.AddContent)
	lms.Post("/courses/:id/enroll", handler.EnrollStudent)
	lms.Get("/enrollments", handler.MyEnrollments)
	lms.Post("/assignments/:id/submit", handler.SubmitAssignment)
	lms.Post("/quizzes/:id/submit", handler.SubmitQuiz)
}

func RegisterLibraryRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewLibraryHandler(db, cfg)
	library := api.Group("/library", middleware.NewAuthMiddleware(cfg).Authenticate)

	library.Get("/books", handler.ListBooks)
	library.Get("/books/:id", handler.GetBook)
	library.Post("/books", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "library"), handler.AddBook)
	library.Put("/books/:id", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "library"), handler.UpdateBook)
	library.Post("/issue", handler.IssueBook)
	library.Post("/return", handler.ReturnBook)
	library.Get("/issues", handler.ListIssues)
	library.Get("/reports", handler.LibraryReports)
}

func RegisterTransportRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewTransportHandler(db, cfg)
	transport := api.Group("/transport", middleware.NewAuthMiddleware(cfg).Authenticate)

	transport.Get("/routes", handler.ListRoutes)
	transport.Post("/routes", middleware.NewAuthMiddleware(cfg).RequireRole("admin"), handler.CreateRoute)
	transport.Get("/vehicles", handler.ListVehicles)
	transport.Post("/vehicles", middleware.NewAuthMiddleware(cfg).RequireRole("admin"), handler.AddVehicle)
	transport.Get("/assignments", handler.ListAssignments)
	transport.Post("/assignments", middleware.NewAuthMiddleware(cfg).RequireRole("admin"), handler.CreateAssignment)
	transport.Get("/boarding", handler.ListBoardingLogs)
	transport.Post("/boarding", handler.LogBoarding)
	transport.Get("/tracking", handler.LiveTracking)
}

func RegisterAnalyticsRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewAnalyticsHandler(db, cfg)
	analytics := api.Group("/analytics", middleware.NewAuthMiddleware(cfg).Authenticate)

	analytics.Get("/dashboard", handler.Dashboard)
	analytics.Get("/kpi", handler.KPIMetrics)
	analytics.Get("/reports", handler.CustomReports)
	analytics.Get("/export", handler.ExportData)

	ai := api.Group("/ai", middleware.NewAuthMiddleware(cfg).Authenticate)
	ai.Get("/predictions", handler.Predictions)
	ai.Get("/performance", handler.PerformancePrediction)
	ai.Get("/anomalies", handler.AnomalyDetection)
	ai.Post("/chat", handler.ChatBot)
}

func RegisterAdmissionRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewAdmissionHandler(db, cfg)
	admission := api.Group("/admission", middleware.NewAuthMiddleware(cfg).Authenticate)

	admission.Get("/leads", handler.ListLeads)
	admission.Get("/leads/:id", handler.GetLead)
	admission.Post("/leads", handler.CreateLead)
	admission.Put("/leads/:id", handler.UpdateLead)
	admission.Post("/leads/:id/convert", handler.ConvertLead)
	admission.Get("/forms", handler.ListForms)
	admission.Post("/forms/submit", handler.SubmitForm)
	admission.Get("/pipeline", handler.PipelineView)
}

func RegisterMessageRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewMessageHandler(db, cfg)
	messages := api.Group("/messages", middleware.NewAuthMiddleware(cfg).Authenticate)

	messages.Get("", handler.List)
	messages.Get("/:id", handler.Get)
	messages.Post("", handler.Send)
	messages.Put("/:id/read", handler.MarkRead)
	
	notifs := api.Group("/notifications", middleware.NewAuthMiddleware(cfg).Authenticate)
	notifs.Get("", handler.ListNotifications)
	notifs.Put("/:id/read", handler.MarkNotificationRead)
	
	announce := api.Group("/announcements", middleware.NewAuthMiddleware(cfg).Authenticate)
	announce.Get("", handler.ListAnnouncements)
	announce.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("admin", "teacher"), handler.CreateAnnouncement)

	whatsapp := api.Group("/whatsapp", middleware.NewAuthMiddleware(cfg).Authenticate)
	whatsapp.Post("/send", handler.SendWhatsApp)
}

func RegisterTenantRoutes(api fiber.Router, db *pgxpool.Pool, cfg *config.Config) {
	handler := handlers.NewTenantHandler(db, cfg)
	tenants := api.Group("/tenants", middleware.NewAuthMiddleware(cfg).Authenticate)
	
	tenants.Get("", middleware.NewAuthMiddleware(cfg).RequireRole("super_admin"), handler.List)
	tenants.Get("/:id", handler.Get)
	tenants.Post("", middleware.NewAuthMiddleware(cfg).RequireRole("super_admin"), handler.Create)
	tenants.Put("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("super_admin", "admin"), handler.Update)
	tenants.Delete("/:id", middleware.NewAuthMiddleware(cfg).RequireRole("super_admin"), handler.Delete)
	tenants.Get("/:id/settings", handler.GetSettings)
	tenants.Put("/:id/settings", middleware.NewAuthMiddleware(cfg).RequireRole("super_admin", "admin"), handler.UpdateSettings)
}
