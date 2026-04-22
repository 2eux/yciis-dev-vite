package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID                 uuid.UUID `json:"id"`
	Name               string   `json:"name"`
	Code               string   `json:"code"`
	Timezone           string   `json:"timezone"`
	CurrencyCode       string   `json:"currency_code"`
	LogoURL            string   `json:"logo_url"`
	Address            string   `json:"address"`
	Phone              string   `json:"phone"`
	Email              string   `json:"email"`
	IsActive           bool     `json:"is_active"`
	SubscriptionPlan  string   `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	Settings          string   `json:"settings"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type User struct {
	ID               uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	Email           string   `json:"email"`
	PasswordHash    string   `json:"-"`
	Role            string   `json:"role"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Phone           string   `json:"phone"`
	AvatarURL       string   `json:"avatar_url"`
	IsActive        bool     `json:"is_active"`
	IsVerified      bool     `json:"is_verified"`
	TwoSecretEnabled bool    `json:"two_secret_enabled"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP    string    `json:"last_login_ip"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	FailedLoginAttempts int    `json:"failed_login_attempts"`
	LockedUntil     *time.Time `json:"locked_until"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserSession struct {
	ID          uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	TokenHash  string   `json:"-"`
	DeviceInfo string   `json:"device_info"`
	IPAddress string   `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name       string   `json:"name"`
	Module    string   `json:"module"`
	Action    string   `json:"action"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	Role          string    `json:"role"`
	PermissionID uuid.UUID `json:"permission_id"`
}

type AuditLog struct {
	ID          uuid.UUID `json:"id"`
	TenantID    *uuid.UUID `json:"tenant_id"`
	UserID     *uuid.UUID `json:"user_id"`
	Action     string     `json:"action"`
	Module     string     `json:"module"`
	EntityType string    `json:"entity_type"`
	EntityID   *uuid.UUID `json:"entity_id"`
	OldValues  string    `json:"old_values"`
	NewValues  string    `json:"new_values"`
	IPAddress string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt time.Time  `json:"created_at"`
}

type AcademicYear struct {
	ID        uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Name    string   `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool     `json:"is_current"`
	IsActive  bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Section struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	Name            string   `json:"name"`
	GradeLevel     int      `json:"grade_level"`
	Room           string   `json:"room"`
	Capacity       int      `json:"capacity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SectionStudent struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	StudentID      uuid.UUID `json:"student_id"`
	SectionID     uuid.UUID `json:"section_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	RollNumber    string    `json:"roll_number"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Subject struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_name"`
	Name       string   `json:"name"`
	Code      string   `json:"code"`
	Description string  `json:"description"`
	IsOptional bool    `json:"is_optional"`
	IsActive   bool    `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SectionSubject struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	SectionID      uuid.UUID `json:"section_id"`
	SubjectID      uuid.UUID `json:"subject_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Student struct {
	ID                     uuid.UUID `json:"id"`
	TenantID               uuid.UUID `json:"tenant_id"`
	UserID                 *uuid.UUID `json:"user_id"`
	StudentID              string    `json:"student_id"`
	Gender                 string    `json:"gender"`
	DateOfBirth           *time.Time `json:"date_of_birth"`
	PlaceOfBirth          string    `json:"place_of_birth"`
	Nationality           string    `json:"nationality"`
	Religion              string    `json:"religion"`
	BloodType             string    `json:"blood_type"`
	Address               string    `json:"address"`
	City                  string    `json:"city"`
	Province              string    `json:"province"`
	PostalCode            string    `json:"postal_code"`
	EmergencyContactName string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	EmergencyContactRelation string `json:"emergency_contact_relation"`
	Notes                 string    `json:"notes"`
	Documents             string    `json:"documents"`
	HealthInfo           string    `json:"health_info"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type StudentParent struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	StudentID     uuid.UUID `json:"student_id"`
	UserID        *uuid.UUID `json:"user_id"`
	Relation      string    `json:"relation"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Occupation    string    `json:"occupation"`
	Company       string    `json:"company"`
	IncomeBracket string    `json:"income_bracket"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Teacher struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	UserID          *uuid.UUID `json:"user_id"`
	EmployeeID     string    `json:"employee_id"`
	Gender         string    `json:"gender"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Qualification  string    `json:"qualification"`
	Specialization string    `json:"specialization"`
	ExperienceYears int     `json:"experience_years"`
	JoinDate       *time.Time `json:"join_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TeacherSubject struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	TeacherID       uuid.UUID `json:"teacher_id"`
	SectionSubjectID uuid.UUID `json:"section_subject_id"`
	AcademicYearID  uuid.UUID `json:"academic_year_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Timetable struct {
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	SectionID        uuid.UUID `json:"section_id"`
	SectionSubjectID uuid.UUID `json:"section_subject_id"`
	TeacherID        *uuid.UUID `json:"teacher_id"`
	DayOfWeek        int       `json:"day_of_week"`
	PeriodStart      int       `json:"period_start"`
	PeriodEnd        int       `json:"period_end"`
	Room             string    `json:"room"`
	AcademicYearID  uuid.UUID `json:"academic_year_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type StudentAttendance struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	StudentID      uuid.UUID `json:"student_id"`
	SectionID      uuid.UUID `json:"section_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	Date           time.Time `json:"date"`
	Status         string    `json:"status"`
	TimeIn         *time.Time `json:"time_in"`
	TimeOut        *time.Time `json:"time_out"`
	Remarks        string    `json:"remarks"`
	MarkedBy       *uuid.UUID `json:"marked_by"`
	DeviceID       string    `json:"device_id"`
	CreatedAt      time.Time  `json:"created_at"`
}

type StaffAttendance struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	StaffID        uuid.UUID `json:"staff_id"`
	Date           time.Time `json:"date"`
	Status         string    `json:"status"`
	TimeIn         *time.Time `json:"time_in"`
	TimeOut        *time.Time `json:"time_out"`
	Remarks        string    `json:"remarks"`
	MarkedBy       *uuid.UUID `json:"marked_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type Exam struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	SectionID      uuid.UUID `json:"section_id"`
	SubjectID      uuid.UUID `json:"subject_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	DurationMinutes int    `json:"duration_minutes"`
	TotalMarks     float64   `json:"total_marks"`
	PassingMarks    float64   `json:"passing_marks"`
	Instructions  string    `json:"instructions"`
	IsPublished   bool      `json:"is_published"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ExamQuestion struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	ExamID        uuid.UUID `json:"exam_id"`
	QuestionText string   `json:"question_text"`
	QuestionType string   `json:"question_type"`
	OptionA       string   `json:"option_a"`
	OptionB       string   `json:"option_b"`
	OptionC       string   `json:"option_c"`
	OptionD       string   `json:"option_d"`
	CorrectAnswer string  `json:"correct_answer"`
	Marks         float64  `json:"marks"`
	SortOrder     int      `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
}

type StudentExamMark struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	StudentID    uuid.UUID `json:"student_id"`
	ExamID       uuid.UUID `json:"exam_id"`
	QuestionID   uuid.UUID `json:"question_id"`
	AnswerText   string    `json:"answer_text"`
	MarksObtained float64   `json:"marks_obtained"`
	GradedBy     *uuid.UUID `json:"graded_by"`
	GradedAt    *time.Time `json:"graded_at"`
	Remarks      string    `json:"remarks"`
	CreatedAt    time.Time `json:"created_at"`
}

type ReportCard struct {
	ID                    uuid.UUID `json:"id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	StudentID            uuid.UUID `json:"student_id"`
	AcademicYearID       uuid.UUID `json:"academic_year_id"`
	SectionID           uuid.UUID `json:"section_id"`
	Term                 string   `json:"term"`
	TotalMarks           float64  `json:"total_marks"`
	Percentage          float64  `json:"percentage"`
	Grade                string   `json:"grade"`
	Rank                int       `json:"rank"`
	AttendancePercentage float64  `json:"attendance_percentage"`
	TeacherRemarks       string   `json:"teacher_remarks"`
	PrincipalRemarks     string   `json:"principal_remarks"`
	GeneratedAt         *time.Time `json:"generated_at"`
	PDFURL              string   `json:"pdf_url"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type LMSCourse struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	TeacherID      *uuid.UUID `json:"teacher_id"`
	AcademicYearID *uuid.UUID `json:"academic_year_id"`
	Title          string    `json:"title"`
	Description   string    `json:"description"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	IsPublished   bool      `json:"is_published"`
	IsFree        bool      `json:"is_free"`
	Price         float64   `json:"price"`
	Language      string    `json:"language"`
	Difficulty    string    `json:"difficulty"`
	DurationHours float64   `json:"duration_hours"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LMSSection struct {
	ID          uuid.UUID `json:"id"`
	CourseID   uuid.UUID `json:"course_id"`
	Title      string    `json:"title"`
	Description string  `json:"description"`
	SortOrder  int      `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

type LMSContent struct {
	ID              uuid.UUID `json:"id"`
	SectionID      uuid.UUID `json:"section_id"`
	Title          string   `json:"title"`
	ContentType   string   `json:"content_type"`
	ContentURL    string   `json:"content_url"`
	ContentText   string   `json:"content_text"`
	DurationMinutes int   `json:"duration_minutes"`
	SortOrder     int      `json:"sort_order"`
	IsPreview     bool     `json:"is_preview"`
	CreatedAt    time.Time `json:"created_at"`
}

type LMSEnrollment struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	CourseID         uuid.UUID `json:"course_id"`
	StudentID        uuid.UUID `json:"student_id"`
	ProgressPercentage float64 `json:"progress_percentage"`
	CompletedAt      *time.Time `json:"completed_at"`
	EnrolledAt       time.Time `json:"enrolled_at"`
}

type FeeStructure struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	Name          string    `json:"name"`
	Description  string    `json:"description"`
	Amount        float64   `json:"amount"`
	DueDate       *time.Time `json:"due_date"`
	IsRecurring   bool      `json:"is_recurring"`
	Frequency    string    `json:"frequency"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FeeDiscount struct {
	ID            uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string   `json:"name"`
	DiscountType string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	Criteria     string   `json:"criteria"`
	IsActive     bool     `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StudentFee struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	StudentID      uuid.UUID `json:"student_id"`
	FeeStructureID uuid.UUID `json:"fee_structure_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	Amount         float64   `json:"amount"`
	DiscountAmount float64   `json:"discount_amount"`
	FinalAmount    float64   `json:"final_amount"`
	Status         string    `json:"status"`
	DueDate        *time.Time `json:"due_date"`
	PaidAmount     float64   `json:"paid_amount"`
	PaidAt        *time.Time `json:"paid_at"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PaymentTransaction struct {
	ID                    uuid.UUID `json:"id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	StudentFeeID         *uuid.UUID `json:"student_fee_id"`
	Amount               float64   `json:"amount"`
	PaymentMethod        string    `json:"payment_method"`
	Gateway              string    `json:"gateway"`
	GatewayTransactionID string   `json:"gateway_transaction_id"`
	GatewayResponse     string    `json:"gateway_response"`
	Status              string    `json:"status"`
	PaidAt              *time.Time `json:"paid_at"`
	CreatedAt           time.Time `json:"created_at"`
}

type Staff struct {
	ID                     uuid.UUID `json:"id"`
	TenantID               uuid.UUID `json:"tenant_id"`
	UserID                 *uuid.UUID `json:"user_id"`
	EmployeeID             string    `json:"employee_id"`
	Department             string    `json:"department"`
	Position               string    `json:"position"`
	JoinDate               *time.Time `json:"join_date"`
	EmploymentType         string    `json:"employment_type"`
	Status                 string    `json:"status"`
	Salary                 float64   `json:"salary"`
	BankAccount            string    `json:"bank_account"`
	BankName               string    `json:"bank_name"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type LeaveRequest struct {
	ID          uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	StaffID    uuid.UUID `json:"staff_id"`
	LeaveType  string    `json:"leave_type"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalDays int       `json:"total_days"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	ApprovedBy *uuid.UUID `json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at"`
	Remarks   string    `json:"remarks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PayrollRun struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
	Status       string    `json:"status"`
	TotalGross  float64   `json:"total_gross"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalNet   float64   `json:"total_net"`
	ProcessedBy *uuid.UUID `json:"processed_by"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type PayrollItem struct {
	ID          uuid.UUID `json:"id"`
	PayrollRunID uuid.UUID `json:"payroll_run_id"`
	StaffID    uuid.UUID `json:"staff_id"`
	BasicSalary float64 `json:"basic_salary"`
	Allowances string   `json:"allowances"`
	Deductions string  `json:"deductions"`
	GrossSalary float64 `json:"gross_salary"`
	NetSalary  float64 `json:"net_salary"`
	CreatedAt time.Time `json:"created_at"`
}

type AdmissionLead struct {
	ID           uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Gender     string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	GradeApplied int     `json:"grade_applied"`
	Source     string    `json:"source"`
	Status     string    `json:"status"`
	AssignedTo *uuid.UUID `json:"assigned_to"`
	Notes      string    `json:"notes"`
	Documents string    `json:"documents"`
	FollowUpAt *time.Time `json:"follow_up_at"`
	ConvertedAt *time.Time `json:"converted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransportRoute struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	Name              string   `json:"name"`
	StartPoint        string   `json:"start_point"`
	EndPoint          string   `json:"end_point"`
	Waypoints        string   `json:"waypoints"`
	DistanceKm       float64  `json:"distance_km"`
	EstimatedTimeMin int     `json:"estimated_time_minutes"`
	CreatedAt       time.Time `json:"created_at"`
}

type Vehicle struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	VehicleNumber string   `json:"vehicle_number"`
	VehicleType  string   `json:"vehicle_type"`
	Model        string   `json:"model"`
	Capacity     int      `json:"capacity"`
	DriverName   string   `json:"driver_name"`
	DriverPhone  string   `json:"driver_phone"`
	InsuranceExpiry *time.Time `json:"insurance_expiry"`
	FitnessExpiry *time.Time `json:"fitness_expiry"`
	IsActive     bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type RouteAssignment struct {
	ID              uuid.UUID `json:"id"`
	RouteID         uuid.UUID `json:"route_id"`
	VehicleID      uuid.UUID `json:"vehicle_id"`
	DriverID       *uuid.UUID `json:"driver_id"`
	AcademicYearID uuid.UUID `json:"academic_year_id"`
	StartTime      *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	Days          string    `json:"days"`
	CreatedAt     time.Time `json:"created_at"`
}

type BoardingLog struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	StudentID         uuid.UUID `json:"student_id"`
	RouteAssignmentID uuid.UUID `json:"route_assignment_id"`
	Date              time.Time `json:"date"`
	PickupTime       *time.Time `json:"pickup_time"`
	PickupLocation   string    `json:"pickup_location"`
	DropoffTime      *time.Time `json:"dropoff_time"`
	DropoffLocation  string    `json:"dropoff_location"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type Book struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	ISBN           string    `json:"isbn"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Publisher      string    `json:"publisher"`
	PublicationYear int     `json:"publication_year"`
	Category       string    `json:"category"`
	Location       string    `json:"location"`
	TotalCopies    int       `json:"total_copies"`
	AvailableCopies int      `json:"available_copies"`
	CreatedAt     time.Time `json:"created_at"`
}

type BookIssue struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	BookID    uuid.UUID `json:"book_id"`
	StudentID *uuid.UUID `json:"student_id"`
	StaffID   *uuid.UUID `json:"staff_id"`
	IssueDate time.Time `json:"issue_date"`
	DueDate   time.Time `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	SenderID     *uuid.UUID `json:"sender_id"`
	ReceiverID   *uuid.UUID `json:"receiver_id"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	MessageType string    `json:"message_type"`
	Priority     string    `json:"priority"`
	IsRead       bool      `json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	UserID           uuid.UUID `json:"user_id"`
	Title            string   `json:"title"`
	Body             string   `json:"body"`
	NotificationType string   `json:"notification_type"`
	Data             string   `json:"data"`
	IsRead           bool     `json:"is_read"`
	ReadAt          *time.Time `json:"read_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type KPIMetric struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	Name        string   `json:"name"`
	Module      string   `json:"module"`
	MetricType  string   `json:"metric_type"`
	Value      float64  `json:"value"`
	TargetValue float64 `json:"target_value"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd  time.Time `json:"period_end"`
	CalculatedAt time.Time `json:"calculated_at"`
}

type AIPrediction struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	EntityType    string   `json:"entity_type"`
	EntityID     uuid.UUID `json:"entity_id"`
	PredictionType string   `json:"prediction_type"`
	PredictionValue float64 `json:"prediction_value"`
	Confidence   float64  `json:"confidence"`
	ModelVersion string   `json:"model_version"`
	Features     string   `json:"features"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Role       string   `json:"role"`
	TenantID   *uuid.UUID `json:"tenant_id"`
	AvatarURL  string   `json:"avatar_url"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type DashboardOverview struct {
	TotalStudents  int     `json:"total_students"`
	TotalTeachers int     `json:"total_teachers"`
	TotalStaff     int     `json:"total_staff"`
	TotalSections int     `json:"total_sections"`
}

type AttendanceStats struct {
	Present      int     `json:"present"`
	Absent       int     `json:"absent"`
	Late         int     `json:"late"`
	Percentage  float64 `json:"percentage"`
}

type FeeStats struct {
	Collected float64 `json:"collected"`
	Pending  float64 `json:"pending"`
	Overdue  float64 `json:"overdue"`
}

type DashboardResponse struct {
	Overview        *DashboardOverview `json:"overview"`
	TodayAttendance *AttendanceStats   `json:"today_attendance"`
	FeeCollection  *FeeStats         `json:"fee_collection"`
}