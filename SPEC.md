# EDUSYS PRO - Enterprise School ERP System

## System Architecture Document

### Version: 1.0.0
### Architecture: Microservices-ready HMVC / Domain-Driven Design
### Tech Stack: Go (Fiber) + React + TypeScript + PostgreSQL

---

## 1. SYSTEM OVERVIEW

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         EDUSYS PRO ARCHITECTURE                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    API GATEWAY (Caddy/Nginx)                │   │
│  │                  Rate Limit | Auth | Routing                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│         ┌────────────────────┼────────────────────┐                │
│         ▼                    ▼                    ▼                │
│  ┌────────────┐      ┌────────────┐      ┌────────────┐          │
│  │  AUTH SVC  │      │  CORE API  │      │  ANALYTICS │          │
│  │  (Auth)   │      │  (Mixed)  │      │   (OLAP)   │          │
│  └────────────┘      └────────────┘      └────────────┘          │
│         │                    │                    │                │
│         └────────────────────┼────────────────────┘                │
│                              ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              MESSAGE BROKER (Kafka / NATS)                    │   │
│  │            Event Streaming | Async Processing                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                     │
│         ┌────────────────────┼────────────────────┐                │
│         ▼                    ▼                    ▼                │
│  ┌────────────┐      ┌────────────┐      ┌────────────┐          │
│  │ PostgreSQL │      │   Redis    │      │  S3/Minio  │          │
│  │  (OLTP)    │      │  (Cache)   │      │  (Files)   │          │
│  └────────────┘      └────────────┘      └────────────┘          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. MODULE ARCHITECTURE (HMVC/DDD)

```
internal/
├── core/                          # Core system modules
│   ├── auth/                      # Authentication & RBAC
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── tenant/                    # Multi-tenant system
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   ├── models.go
│   │   └── services.go
│   ├── audit/                     # Audit logging
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   └── notification/              # Push/Email/WhatsApp
│       ├── handlers.go
│       ├── models.go
│       ├── services.go
│       └── providers/
│           ├── whatsapp.go
│           ├── email.go
│           └── push.go
│
├── academic/                      # Academic modules
│   ├── student/                   # Student management
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── services.go
│   │   └── validators.go
│   ├── admission/                 # CRM & admissions
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── services.go
│   │   └── workflows.go
│   ├── academic/                  # Curriculum & subjects
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── timetable/                 # Timetable engine
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── generator.go
│   │   └── services.go
│   ├── attendance/                 # Attendance system
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── exam/                       # Exams & grading
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── services.go
│   │   └── grading.go
│   └── lms/                        # Learning Management
│       ├── handlers.go
│       ├── models.go
│       ├── repository.go
│       ├── services.go
│       └── content/
│           ├── video.go
│           ├── quiz.go
│           └── assignment.go
│
├── finance/                       # Finance modules
│   ├── fee/                        # Fee structure
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── invoice/                   # Billing & invoicing
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── services.go
│   │   └── payment/
│   │       ├── midtrans.go
│   │       └── xendit.go
│   └── report/                     # Financial reports
│       ├── handlers.go
│       ├── models.go
│       ���─�� services.go
│
├── hr/                           # HR & Payroll
│   ├── staff/                     # Staff management
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── leave/                     # Leave management
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   ├── payroll/                   # Payroll automation
│   │   ├── handlers.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── services.go
│   └── recruitment/               # Recruitment (AI)
│       ├── handlers.go
│       ├── models.go
│       ├── services.go
│       └── ai_screening.go
│
├── transport/                     # Transport & GPS
│   ├── route/                     # Route planning
│   │   ├── handlers.go
│   │   ├── models.go
│   │   └── services.go
│   ├── vehicle/                   # Vehicle tracking
│   │   ├── handlers.go
│   │   ├── models.go
│   │   └── services.go
│   └── boarding/                  # Student boarding
│       ├── handlers.go
│       ├── models.go
│       └── services.go
│
├── library/                      # Library & Inventory
│   ├── book/                      # Book tracking
│   │   ├── handlers.go
│   │   ├── models.go
│   │   └── services.go
│   ├── asset/                     # Asset management
│   │   ├── handlers.go
│   │   ├── models.go
│   │   └── services.go
│   └── stock/                     # Stock movement
│       ├── handlers.go
│       ├── models.go
│       └── services.go
│
└── analytics/                    # Analytics (AI-ready)
    ├── dashboard/                 # Real-time dashboards
    │   ├── handlers.go
    │   ├── models.go
    │   └── services.go
    ├── kpi/                       # KPI monitoring
    │   ├── handlers.go
    │   ├── models.go
    │   └── services.go
    ├── prediction/                # Predictive analytics
    │   ├── handlers.go
    │   ├── models.go
    │   ├── services.go
    │   └── student_performance.go
    └── anomaly/                   # Anomaly detection
        ├── handlers.go
        ├── models.go
        ├── services.go
        └── detector.go
```

---

## 3. DATABASE SCHEMA (PostgreSQL)

```sql
-- =====================================================
-- EDUSYS PRO - DATABASE SCHEMA (PostgreSQL)
-- =====================================================

-- Enable required extensions
-- Use pgcrypto for UUID generation in Supabase (gen_random_uuid())
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- =====================================================
-- CORE TABLES (Multi-tenant, Auth, Audit)
-- =====================================================

-- Tenants (Multi-school)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    currency_code CHAR(3) DEFAULT 'IDR',
    logo_url VARCHAR(500),
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    subscription_plan VARCHAR(50) DEFAULT 'basic',
    subscription_expires_at TIMESTAMP,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Users (All roles)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    email CITEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'teacher', 'student', 'parent', 'finance', 'hr')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    two_secret_enabled BOOLEAN DEFAULT false,
    two_secret_secret VARCHAR(255),
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    password_changed_at TIMESTAMP,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- User sessions (JWT refresh tokens)
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    device_info VARCHAR(255),
    ip_address VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Permissions (RBAC)
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    module VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(module, action)
);

-- Role permissions
CREATE TABLE role_permissions (
    role VARCHAR(20) NOT NULL,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(role, permission_id)
);

-- Audit logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    module VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ACADEMIC TABLES
-- =====================================================

-- Academic years
CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Sections (Classes)
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    grade_level INT NOT NULL CHECK (grade_level BETWEEN 0 AND 12),
    room VARCHAR(50),
    capacity INT DEFAULT 35,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, academic_year_id, name)
);

-- Section mappings (Student-class)
CREATE TABLE section_students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    roll_number VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'promoted', 'transferred', 'dropped')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, academic_year_id)
);

-- Subjects
CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    is_optional BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- Subject to section mapping
CREATE TABLE section_subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(section_id, subject_id, academic_year_id)
);

-- Student details
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    student_id VARCHAR(50) UNIQUE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    date_of_birth DATE,
    place_of_birth VARCHAR(100),
    nationality VARCHAR(50),
    religion VARCHAR(50),
    blood_type VARCHAR(5),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relation VARCHAR(50),
    notes TEXT,
    documents JSONB DEFAULT '[]',
    health_info JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Student family (Parents/Guardians)
CREATE TABLE student_parents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    relation VARCHAR(50) NOT NULL CHECK (relation IN ('father', 'mother', 'guardian', 'other')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    occupation VARCHAR(100),
    company VARCHAR(200),
    income_bracket VARCHAR(50),
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Teachers
CREATE TABLE teachers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    gender VARCHAR(10),
    date_of_birth DATE,
    qualification VARCHAR(200),
    specialization VARCHAR(200),
    experience_years INT,
    join_date DATE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'on_leave', 'resigned')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Teacher subject assignments
CREATE TABLE teacher_subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES teachers(id) ON DELETE CASCADE,
    section_subject_id UUID REFERENCES section_subjects(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(teacher_id, section_subject_id, academic_year_id)
);

-- Timetables
CREATE TABLE timetables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    section_subject_id UUID REFERENCES section_subjects(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    period_start INT NOT NULL,
    period_end INT NOT NULL,
    room VARCHAR(50),
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ATTENDANCE TABLES
-- =====================================================

-- Student attendance
CREATE TABLE student_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'late', 'excused')),
    time_in TIME,
    time_out TIME,
    remarks TEXT,
    marked_by UUID REFERENCES users(id) ON DELETE SET NULL,
    device_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, date)
);

-- Staff attendance
CREATE TABLE staff_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL,
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'late', 'on_leave')),
    time_in TIME,
    time_out TIME,
    remarks TEXT,
    marked_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(staff_id, date)
);

-- =====================================================
-- EXAM & GRADING TABLES
-- =====================================================

-- Exams
CREATE TABLE exams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('quiz', 'midterm', 'final', 'semester', 'annual')),
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    duration_minutes INT,
    total_marks DECIMAL(10,2),
    passing_marks DECIMAL(10,2),
    instructions TEXT,
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Exam questions
CREATE TABLE exam_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(20) NOT NULL CHECK (question_type IN ('mcq', 'true_false', 'short', 'long')),
    option_a TEXT,
    option_b TEXT,
    option_c TEXT,
    option_d TEXT,
    correct_answer TEXT NOT NULL,
    marks DECIMAL(10,2) NOT NULL,
    sort_order INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Student exam marks
CREATE TABLE student_exam_marks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
   question_id UUID REFERENCES exam_questions(id) ON DELETE CASCADE,
    answer_text TEXT,
    marks_obtained DECIMAL(10,2),
    graded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    graded_at TIMESTAMP,
    remarks TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, exam_id, question_id)
);

-- Report cards
CREATE TABLE report_cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    term VARCHAR(50) NOT NULL,
    total_marks DECIMAL(10,2),
    percentage DECIMAL(5,2),
    grade VARCHAR(5),
    rank INT,
    attendance_percentage DECIMAL(5,2),
    teacher_remarks TEXT,
    principal_remarks TEXT,
    generated_at TIMESTAMP,
    pdf_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, academic_year_id, term)
);

-- =====================================================
-- LMS TABLES
-- =====================================================

-- LMS Courses
CREATE TABLE lms_courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url VARCHAR(500),
    is_published BOOLEAN DEFAULT false,
    is_free BOOLEAN DEFAULT false,
    price DECIMAL(10,2) DEFAULT 0,
    language VARCHAR(20) DEFAULT 'en',
    difficulty VARCHAR(20) DEFAULT 'beginner',
    duration_hours DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Course sections/modules
CREATE TABLE lms_sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID REFERENCES lms_courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Course content
CREATE TABLE lms_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    section_id UUID REFERENCES lms_sections(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('video', 'document', 'quiz', 'assignment', 'link')),
    content_url VARCHAR(500),
    content_text TEXT,
    duration_minutes INT,
    sort_order INT,
    is_preview BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Course enrollments
CREATE TABLE lms_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    course_id UUID REFERENCES lms_courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    progress_percentage DECIMAL(5,2) DEFAULT 0,
    completed_at TIMESTAMP,
    enrolled_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, student_id)
);

-- =====================================================
-- FINANCE TABLES
-- =====================================================

-- Fee structures
CREATE TABLE fee_structures (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    amount DECIMAL(15,2) NOT NULL,
    due_date DATE,
    is_recurring BOOLEAN DEFAULT false,
    frequency VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Fee discounts
CREATE TABLE fee_discounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
    discount_value DECIMAL(15,2) NOT NULL,
    criteria JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Student fee assignments
CREATE TABLE student_fees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    fee_structure_id UUID REFERENCES fee_structures(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    final_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'partial', 'paid', 'overdue', 'waived')),
    due_date DATE,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    paid_at TIMESTAMP,
    payment_method VARCHAR(50),
    transaction_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payment transactions
CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_fee_id UUID REFERENCES student_fees(id) ON DELETE SET NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    gateway VARCHAR(50),
    gateway_transaction_id VARCHAR(100),
    gateway_response JSONB,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'refunded')),
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- HR TABLES
-- =====================================================

-- Staff details
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    department VARCHAR(100),
    position VARCHAR(200),
    join_date DATE,
    employment_type VARCHAR(20) DEFAULT 'full_time' CHECK (employment_type IN ('full_time', 'part_time', 'contract', 'intern')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'on_leave', 'terminated')),
    salary DECIMAL(15,2),
    bank_account VARCHAR(50),
    bank_name VARCHAR(100),
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Leave requests
CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    staff_id UUID REFERENCES staff(id) ON DELETE CASCADE,
    leave_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INT NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    remarks TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payroll runs
CREATE TABLE payroll_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'processing', 'completed', 'cancelled')),
    total_gross DECIMAL(15,2),
    total_deductions DECIMAL(15,2),
    total_net DECIMAL(15,2),
    processed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Payroll items
CREATE TABLE payroll_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payroll_run_id UUID REFERENCES payroll_runs(id) ON DELETE CASCADE,
    staff_id UUID REFERENCES staff(id) ON DELETE CASCADE,
    basic_salary DECIMAL(15,2) NOT NULL,
    allowances JSONB DEFAULT '{}',
    deductions JSONB DEFAULT '{}',
    gross_salary DECIMAL(15,2),
    net_salary DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ADMISSION CRM TABLES
-- =====================================================

-- Admission pipeline
CREATE TABLE admission_leads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    gender VARCHAR(10),
    date_of_birth DATE,
    grade_applied INT,
    source VARCHAR(100),
    status VARCHAR(20) DEFAULT 'lead' CHECK (status IN ('lead', 'prospect', 'registered', 'rejected', 'waitlisted')),
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    documents JSONB DEFAULT '[]',
    follow_up_at TIMESTAMP,
    converted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- TRANSPORT TABLES
-- =====================================================

-- Routes
CREATE TABLE transport_routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    start_point VARCHAR(200) NOT NULL,
    end_point VARCHAR(200) NOT NULL,
    waypoints JSONB DEFAULT '[]',
    distance_km DECIMAL(10,2),
    estimated_time_minutes INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Vehicles
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_number VARCHAR(50) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
    model VARCHAR(100),
    capacity INT,
    driver_name VARCHAR(200),
    driver_phone VARCHAR(20),
    insurance_expiry DATE,
    fitness_expiry DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Route assignments
CREATE TABLE route_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id UUID REFERENCES transport_routes(id) ON DELETE CASCADE,
    vehicle_id UUID REFERENCES vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    start_time TIME,
    end_time TIME,
    days JSONB DEFAULT '[1,2,3,4,5]',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Student boarding logs
CREATE TABLE boarding_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    route_assignment_id UUID REFERENCES route_assignments(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    pickup_time TIME,
    pickup_location VARCHAR(200),
    dropoff_time TIME,
    dropoff_location VARCHAR(200),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'boarded', 'dropped', 'noshow')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, route_assignment_id, date)
);

-- =====================================================
-- LIBRARY TABLES
-- =====================================================

-- Books
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    isbn VARCHAR(20) UNIQUE,
    title VARCHAR(500) NOT NULL,
    author VARCHAR(200),
    publisher VARCHAR(200),
    publication_year INT,
    category VARCHAR(100),
    location VARCHAR(100),
    total_copies INT DEFAULT 1,
    available_copies INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Book issues
CREATE TABLE book_issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    student_id UUID,
    staff_id UUID,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE,
    status VARCHAR(20) DEFAULT 'issued' CHECK (status IN ('issued', 'returned', 'lost', 'overdue')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- COMMUNICATION TABLES
-- =====================================================

-- Messages
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    receiver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'general' CHECK (message_type IN ('general', 'announcement', 'notice')),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    notification_type VARCHAR(50) NOT NULL,
    data JSONB DEFAULT '{}',
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ANALYTICS TABLES
-- =====================================================

-- KPI metrics
CREATE TABLE kpi_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    module VARCHAR(50) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    value DECIMAL(15,4),
    target_value DECIMAL(15,4),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    calculated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name, period_start, period_end)
);

-- AI Predictions cache
CREATE TABLE ai_predictions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    prediction_type VARCHAR(50) NOT NULL,
    prediction_value DECIMAL(10,4),
    confidence DECIMAL(5,4),
    model_version VARCHAR(50),
    features JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_students_tenant ON students(tenant_id);
CREATE INDEX idx_student_attendance_date ON student_attendance(date);
CREATE INDEX idx_student_attendance_student_date ON student_attendance(student_id, date);
CREATE INDEX idx_exams_section ON exams(section_id);
CREATE INDEX idx_exams_date ON exams(start_date);
CREATE INDEX idx_student_fees_status ON student_fees(status);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_notifications_user ON notifications(user_id, is_read);
CREATE INDEX idx_messages_thread ON messages(sender_id, receiver_id);

-- =====================================================
-- FUNCTIONS & TRIGGERS
-- =====================================================

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Soft delete trigger
CREATE OR REPLACE FUNCTION set_deleted_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = false AND OLD.is_active = true THEN
        NEW.deleted_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- SEED DATA
-- =====================================================

-- Default permissions
INSERT INTO permissions (name, module, action, description) VALUES
('View Dashboard', 'dashboard', 'view', 'View dashboard'),
('View Students', 'student', 'view', 'View student list'),
('Create Students', 'student', 'create', 'Create new student'),
('Update Students', 'student', 'update', 'Update student details'),
('Delete Students', 'student', 'delete', 'Delete student'),
('View Attendance', 'attendance', 'view', 'View attendance'),
('Mark Attendance', 'attendance', 'mark', 'Mark attendance'),
('View Exams', 'exam', 'view', 'View exams'),
('Create Exams', 'exam', 'create', 'Create exams'),
('Grade Exams', 'exam', 'grade', 'Grade student exams'),
('View Fees', 'fee', 'view', 'View fee information'),
('Create Fees', 'fee', 'create', 'Create fee structures'),
('View Reports', 'report', 'view', 'View reports'),
('Manage Users', 'user', 'manage', 'Manage system users'),
('Manage Settings', 'settings', 'manage', 'Manage system settings');

-- Admin role permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions;

-- Teacher role permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'teacher', id FROM permissions WHERE module IN ('dashboard', 'student', 'attendance', 'exam', 'lms');

-- Student role permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'student', id FROM permissions WHERE module IN ('dashboard', 'lms', 'attendance');

-- Parent role permissions
INSERT INTO role_permissions (role, permission_id)
SELECT 'parent', id FROM permissions WHERE module IN ('dashboard', 'student', 'attendance', 'fee');
```

---

## 4. API ENDPOINTS

### 4.1 Authentication
```
POST   /api/v1/auth/login              - Login with email/password
POST   /api/v1/auth/logout             - Logout (invalidate session)
POST   /api/v1/auth/refresh            - Refresh access token
POST   /api/v1/auth/forgot-password    - Request password reset
POST   /api/v1/auth/reset-password      - Reset password
POST   /api/v1/auth/verify-email       - Verify email address
POST   /api/v1/auth/enable-2fa          - Enable 2FA
POST   /api/v1/auth/verify-2fa          - Verify 2FA code
GET    /api/v1/auth/me                 - Get current user
PUT    /api/v1/auth/profile            - Update profile
```

### 4.2 Tenants (Schools)
```
GET    /api/v1/tenants                  - List tenants
POST   /api/v1/tenants                 - Create tenant
GET    /api/v1/tenants/:id              - Get tenant
PUT    /api/v1/tenants/:id              - Update tenant
DELETE /api/v1/tenants/:id              - Delete tenant
GET    /api/v1/tenants/:id/settings     - Get settings
PUT    /api/v1/tenants/:id/settings     - Update settings
```

### 4.3 Students
```
GET    /api/v1/students                - List students
POST   /api/v1/students                - Create student
GET    /api/v1/students/:id            - Get student
PUT    /api/v1/students/:id           - Update student
DELETE /api/v1/students/:id           - Delete student
GET    /api/v1/students/:id/profile   - Get full profile
GET    /api/v1/students/:id/family     - Get family members
POST   /api/v1/students/:id/family     - Add family member
PUT    /api/v1/students/:id/documents - Update documents
GET    /api/v1/students/:id/attendance - Get attendance
GET    /api/v1/students/:id/grades     - Get grades
GET    /api/v1/students/:id/fees       - Get fees
GET    /api/v1/students/:id/reports    - Get reports
```

### 4.4 Academic
```
GET    /api/v1/academic-years         - List academic years
POST   /api/v1/academic-years        - Create academic year
GET    /api/v1/sections              - List sections/classes
POST   /api/v1/sections             - Create section
GET    /api/v1/sections/:id          - Get section
PUT    /api/v1/sections/:id/students - Assign students
GET    /api/v1/sections/:id/timetable - Get timetable
GET    /api/v1/subjects             - List subjects
POST   /api/v1/subjects             - Create subject
GET    /api/v1/subjects/:id          - Get subject
GET    /api/v1/timetables            - List timetables
POST   /api/v1/timetables            - Create timetable
GET    /api/v1/timetables/generate   - Auto-generate timetable
```

### 4.5 Attendance
```
GET    /api/v1/attendance/students   - Student attendance
POST   /api/v1/attendance/students   - Mark attendance
GET    /api/v1/attendance/staff       - Staff attendance
POST   /api/v1/attendance/staff     - Mark staff attendance
GET    /api/v1/attendance/reports   - Attendance reports
GET    /api/v1/attendance/qr/generate - Generate QR code
POST   /api/v1/attendance/qr/scan   - Scan QR attendance
```

### 4.6 Exams
```
GET    /api/v1/exams                  - List exams
POST   /api/v1/exams                 - Create exam
GET    /api/v1/exams/:id             - Get exam
PUT    /api/v1/exams/:id             - Update exam
DELETE /api/v1/exams/:id            - Delete exam
GET    /api/v1/exams/:id/questions - Get questions
POST   /api/v1/exams/:id/questions  - Add question
PUT    /api/v1/exams/:id/publish     - Publish exam
GET    /api/v1/exams/:id/results     - Get results
POST   /api/v1/exams/:id/grade      - Grade student
GET    /api/v1/exams/:id/export     - Export results
```

### 4.7 LMS
```
GET    /api/v1/lms/courses           - List courses
POST   /api/v1/lms/courses          - Create course
GET    /api/v1/lms/courses/:id      - Get course
PUT    /api/v1/lms/courses/:id     - Update course
DELETE /api/v1/lms/courses/:id     - Delete course
POST   /api/v1/lms/courses/:id/publish - Publish course
GET    /api/v1/lms/courses/:id/content - Get content
POST   /api/v1/lms/courses/:id/content - Add content
POST   /api/v1/lms/courses/:id/enroll - Enroll student
GET    /api/v1/lms/enrollments      - My enrollments
POST   /api/v1/lms/assignments/:id/submit - Submit assignment
POST   /api/v1/lms/quizzes/:id/submit - Submit quiz
```

### 4.8 Finance
```
GET    /api/v1/fees/structures       - List fee structures
POST   /api/v1/fees/structures      - Create fee structure
PUT    /api/v1/fees/structures/:id - Update fee structure
POST   /api/v1/fees/assign          - Assign fees to students
GET    /api/v1/fees/students        - Student fees
GET    /api/v1/fees/students/:id    - Student fee details
POST   /api/v1/fees/students/:id/pay - Process payment
POST   /api/v1/fees/webhook         - Payment webhook
GET    /api/v1/fees/reports         - Financial reports
GET    /api/v1/fees/reports/overview - Overview report
```

### 4.9 HR
```
GET    /api/v1/hr/staff              - List staff
POST   /api/v1/hr/staff             - Create staff
GET    /api/v1/hr/staff/:id         - Get staff
PUT    /api/v1/hr/staff/:id         - Update staff
DELETE /api/v1/hr/staff/:id        - Delete staff
GET    /api/v1/hr/leave             - Leave requests
POST   /api/v1/hr/leave            - Submit leave request
PUT    /api/v1/hr/leave/:id        - Approve/reject leave
GET    /api/v1/hr/payroll           - Payroll runs
POST   /api/v1/hr/payroll          - Run payroll
GET    /api/v1/hr/payroll/:id       - Get payroll run
POST   /api/v1/hr/payroll/:id/process - Process payroll
```

### 4.10 Communication
```
GET    /api/v1/messages             - List messages
POST   /api/v1/messages            - Send message
GET    /api/v1/messages/:id        - Get message
PUT    /api/v1/messages/:id/read   - Mark as read
GET    /api/v1/notifications        - My notifications
PUT    /api/v1/notifications/:id/read - Mark as read
POST   /api/v1/announcements       - Create announcement
GET    /api/v1/announcements      - List announcements
POST   /api/v1/whatsapp/send       - Send WhatsApp message
```

### 4.11 Transport
```
GET    /api/v1/transport/routes     - List routes
POST   /api/v1/transport/routes   - Create route
GET    /api/v1/transport/vehicles  - List vehicles
POST   /api/v1/transport/vehicles  - Add vehicle
GET    /api/v1/transport/assignments - Route assignments
POST   /api/v1/transport/assignments - Assign route
GET    /api/v1/transport/boarding - Boarding logs
POST   /api/v1/transport/boarding - Log pickup/dropoff
GET    /api/v1/transport/tracking  - Live tracking
```

### 4.12 Library
```
GET    /api/v1/library/books      - List books
POST   /api/v1/library/books     - Add book
GET    /api/v1/library/books/:id - Get book
PUT    /api/v1/library/books/:id - Update book
POST   /api/v1/library/issue     - Issue book
POST   /api/v1/library/return    - Return book
GET    /api/v1/library/issues   - Issue records
GET    /api/v1/library/reports   - Library reports
```

### 4.13 Analytics
```
GET    /api/v1/analytics/dashboard  - Dashboard data
GET    /api/v1/analytics/kpi       - KPI metrics
GET    /api/v1/analytics/reports   - Custom reports
GET    /api/v1/analytics/export   - Export data
GET    /api/v1/ai/predictions      - AI predictions
GET    /api/v1/ai/performance     - Performance prediction
GET    /api/v1/ai/anomalies       - Anomaly detection
POST   /api/v1/ai/chat             - Chat with AI
GET    /api/v1/ai/screening       - CV screening results
```

### 4.14 Admission CRM
```
GET    /api/v1/admission/leads     - List leads
POST   /api/v1/admission/leads    - Create lead
GET    /api/v1/admission/leads/:id - Get lead
PUT    /api/v1/admission/leads/:id - Update lead
POST   /api/v1/admission/leads/:id/convert - Convert to student
GET    /api/v1/admission/forms   - Registration forms
POST   /api/v1/admission/forms/submit - Submit form
GET    /api/v1/admission/pipeline - Pipeline view
```

---

## 5. SAMPLE API RESPONSES

### 5.1 Login
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "admin@school.edu",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "tenant_id": "550e8400-e29b-41d4-a716-446655440001",
      "avatar_url": "https://..."
    }
  }
}
```

### 5.2 Student List
```json
{
  "success": true,
  "data": {
    "students": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "student_id": "STU2024001",
        "first_name": "Jane",
        "last_name": "Smith",
        "gender": "female",
        "date_of_birth": "2012-05-15",
        "section": {
          "id": "...",
          "name": "Class 7-A"
        },
        "roll_number": "001",
        "avatar_url": "https://...",
        "status": "active"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 450,
      "total_pages": 23
    }
  }
}
```

### 5.3 Dashboard
```json
{
  "success": true,
  "data": {
    "overview": {
      "total_students": 1250,
      "total_teachers": 85,
      "total_staff": 42,
      "total_sections": 35
    },
    "today_attendance": {
      "present": 1180,
      "absent": 45,
      "late": 25,
      "percentage": 94.4
    },
    "fee_collection": {
      "collected": 125000000,
      "pending": 25000000,
      "overdue": 5000000
    },
    "recent_activities": [...],
    "upcoming_events": [...],
    "alerts": [...]
  }
}
```

---

## 6. FOLDER STRUCTURE

```
yciis-dev-vite/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   └── tenant.go
│   ├── core/
│   │   ├── auth/
│   │   ├── tenant/
│   │   ├── audit/
│   │   └── notification/
│   ├── academic/
│   │   ├── student/
│   │   ├── admission/
│   │   ├── academic/
│   │   ├── timetable/
│   │   ├── attendance/
│   │   ├── exam/
│   │   └── lms/
│   ├── finance/
│   │   ├── fee/
│   │   ├── invoice/
│   │   └── report/
│   ├── hr/
│   │   ├── staff/
│   │   ├── leave/
│   │   └── payroll/
│   ├── transport/
│   ├── library/
│   ├── analytics/
│   └── utils/
├── migrations/
├── pkg/
│   ├── response/
│   ├── errors/
│   ├── models/
│   └── validators/
├── web/
│   ├── public/
│   └── templates/
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── nginx.conf
├── .env.example
├── go.mod
├── go.sum
├── Makefile
└── README.md

src/
├── components/
│   ├── common/
│   │   ├── Button/
│   │   ├── Card/
│   │   ├── Modal/
│   │   ├── Table/
│   │   ├── Input/
│   │   ├── Select/
│   │   └── Pagination/
│   ├── layout/
│   │   ├── Header/
│   │   ├── Footer/
│   │   └── Layout/
│   └── modules/
│       ├── Dashboard/
│       ├── Students/
│       ├── Academic/
│       ├── Attendance/
│       ├── Exams/
│       ├── LMS/
│       ├── Finance/
│       ├── HR/
│       └── Settings/
├── pages/
├── hooks/
├── services/
├── stores/
├── styles/
├── types/
├── utils/
├── App.tsx
├── main.tsx
└── index.html
```

---

## 7. SECURITY REQUIREMENTS

- [x] JWT with short-lived access tokens (15 min)
- [x] Refresh tokens with rotation
- [x] RBAC with granular permissions
- [x] CSRF protection
- [x] XSS protection (sanitization)
- [x] SQL injection prevention (parameterized queries)
- [x] Rate limiting (100 req/min)
- [x] Audit logging (all PII actions)
- [x] Data encryption at rest (AES-256)
- [x] TLS 1.3 for transit
- [x] Input validation
- [x] Account lockout (5 failed attempts)
- [x] Session management
- [x] Tenant isolation

---

## 8. AI INTEGRATION DESIGN

### 8.1 Smart Chatbot
- NLP-powered school assistant
- FAQ automation
- Parent/student queries
- Integration: Dialogflow/Rasa/Custom LLM

### 8.2 Student Performance Prediction
- ML model: Random Forest / XGBoost
- Features: attendance, grades, behavior
- Output: Risk score, recommendations

### 8.3 Auto Report Generation
- Template-based PDF generation
- NLP for remarks
- Historical data analysis

### 8.4 Anomaly Detection
- Fee payment anomalies
- Attendance patterns
- Academic performance drops

### 8.5 CV Screening (HR)
- Resume parsing
- Keyword matching
- Scoring algorithm

---

## 9. SCALING STRATEGY

- Horizontal scaling with Kubernetes
- Read replicas for PostgreSQL
- Redis cluster for sessions/cache
- CDN for static assets
- Message queue for async jobs
- Separate OLTP/OLAP databases

---

## 10. DEPLOYMENT

```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: ./docker/api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://user:pass@db:5432/edusys
      - REDIS_URL=redis://cache:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - db
      - cache
      - kafka

  web:
    build: ./docker/web
    ports:
      - "3000:3000"
    environment:
      - API_URL=http://api:8080

  db:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=edusys

  cache:
    image: redis:7-alpine

  kafka:
    image: confluentinc/cp-kafka:7.5.0

  minio:
    image: minio/minio
    command: server /data

volumes:
  postgres_data:
```

---

## 11. ROADMAP

### MVP (3 months)
- Authentication & RBAC
- Student management
- Basic attendance
- Fee structure
- Reports

### V1.0 (6 months)
- Full academic system
- Exam & grading
- LMS basic
- HR module
- Analytics

### V1.5 (9 months)
- Transport tracking
- Library system
- WhatsApp integration
- AI chatbot

### V2.0 (12 months)
- Predictive analytics
- Auto timetable
- Advanced AI features
- Mobile apps
