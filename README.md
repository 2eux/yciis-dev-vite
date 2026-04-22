# Edusys Pro - Enterprise School ERP System

<p align="center">
  <img src="https://img.shields.io/badge/Version-1.0.0-blue?style=for-the-badge&logo=version" alt="Version" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License" />
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react" alt="React" />
</p>

A comprehensive, enterprise-grade School ERP (Enterprise Resource Planning) system built with modern technologies. Edusys Pro handles all school operations including academics, administration, finance, HR, and more.

---

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [System Architecture](#system-architecture)
- [DFD - Data Flow Diagrams](#dfd---data-flow-diagrams)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Installation Guide](#installation-guide)
- [Project Structure](#project-structure)
- [Security](#security)
- [AI Integration](#ai-integration)
- [Roadmap](#roadmap)
- [License](#license)

---

## ✨ Features

### Core Modules
- **Authentication & RBAC** - JWT-based auth with role-based access control + 2FA (Google Authenticator)
- **Multi-tenant** - SaaS-ready multi-school support
- **Audit Logging** - Full traceability of all actions

### Academic System
- **Student Management** - Profile, family, documents, health records
- **Admission CRM** - Lead → Prospect → Student pipeline
- **Timetable Engine** - Auto-generation of schedules
- **Subjects & Curriculum** - Subject mapping and teacher assignment
- **Attendance** - QR/biometric integration, real-time parent notifications

### Examination & Grading
- **Exam Builder** - Create exams with question banks
- **Mark Entry** - Grade student submissions
- **Report Cards** - PDF export with analytics

### Finance & Accounting
- **Fee Structure Engine** - Flexible fee categories
- **Billing & Invoicing** - Automated invoicing
- **Online Payments** - Midtrans/Xendit integration (Indonesia)
- **Financial Reports** - Comprehensive reporting

### HR & Payroll
- **Staff Management** - Employee records
- **Leave Management** - Approval workflows
- **Payroll Automation** - Salary processing
- **Recruitment** - AI-powered CV screening

### LMS (Learning Management)
- **Course Builder** - Video/content upload
- **Assignments & Quizzes** - Interactive learning
- **Progress Tracking** - Student analytics

### Additional Modules
- **Library** - Book tracking system
- **Transport** - Route planning, GPS tracking
- **Communication** - WhatsApp, Email, Push notifications
- **Analytics** - Real-time dashboards, KPIs, AI predictions

---

## 🛠 Tech Stack

### Backend
| Technology | Purpose |
|-----------|---------|
| **Go 1.21+** | Primary language (Fiber framework) |
| **PostgreSQL 15** | Primary database |
| **Redis 7** | Cache & session storage |
| **JWT** | Authentication |
| **bcrypt** | Password hashing |

### Frontend
| Technology | Purpose |
|-----------|---------|
| **React 18** | UI framework |
| **TypeScript** | Type safety |
| **Tailwind CSS** | Styling (shadcn/ui) |
| **Recharts** | Data visualization |
| **Zustand** | State management |
| **React Query** | Data fetching |

### DevOps & Tools
| Technology | Purpose |
|-----------|---------|
| **Docker** | Containerization |
| **Vite** | Build tool |
| **ESLint** | Code linting |

---

## 🏗 System Architecture

```
┌────────────────────────────────────────���────────────────────────────┐
│                     CLIENT LAYER                           │
├─────────────────────────────────────────────────────────────┤
│  React + TypeScript + Tailwind CSS + shadcn/ui              │
│  (SPA with Vite)                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     API GATEWAY                           │
├─────────────────────────────────────────────────────────────┤
│  Go Fiber + JWT Auth + Rate Limiting + CORS               │
│  Port: 8080                                            │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┬────────────────────┐
         ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Core API    │    │  Analytics  │    │   Queue     │
│  (REST)      │    │   (OLAP)    │    │  (Async)    │
└──────────────┘    └──────────────┘    └──────────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   DATA LAYER                             │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL        Redis           MinIO/S3              │
│  (OLTP)           (Cache)         (Files)               │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 DFD - Data Flow Diagrams

### Level 0: Context Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                                                        │
│   ┌──────────────┐                                     │
│   │   EXTERNAL  │                                     │
│   │   ENTITIES │                                     │
│   └─────┬──────┘                                     │
│         │                                            │
│    ┌────▼───────────��───┐                              │
│    │                 │                              │
│    │  EDUSYS PRO      │ ◄─── System                │
│    │                 │                              │
│    └───────┬─────────┘                              │
│            │                                       │
│    ┌──────▼──────┐                              │
│    │  EXTERNAL   │                              │
│    │  ENTITIES  │                              │
│    └───────────┘                              │
└──────────────────────────────────────────────────────┘
```

External Entities:
- **Students** - View grades, attendance, submit assignments
- **Parents** - View reports, make payments, communicate
- **Teachers** - Mark attendance, upload grades, manage courses
- **Admin** - Full system access, reports, settings
- **HR/Finance** - Payroll, fee management

### Level 1: Main Processes

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐      │
│  │Student │      │Academic │      │Finance │      │   HR   │      │
│  │Module  │      │Module   │      │Module  │      │Module │      │
│  └──┬────┘      └──┬────┘      └──┬────┘      └──┬────┘      │
│     │             │             │            │            │          │
│     └─────────────┴─────────────┴────────────┘            │
│                           │                              │          │
│                    ┌────────▼────────┐    ┌──────▼─────┐   │
│                    │              │    │           │   │
│                    │  PROCESS    │    │  PROCESS  │   │
│                    │   1.0       │    │   2.0     │   │
│                    │  Core Edu   │    │  Security  │   │
│                    │  System   │    │  System   │   │
│                    └──────┬─────┘    └──────┬─────┘   │
│                           │               │          │
│                           └───────┬───────┘          │
│                                   │              │
│                    ┌──────────────▼──────────────┐│
│                    │                             ││
│                    │     DATA STORE           ││
│                    │  ┌──────────┬──────────┐ ││
│                    │  │Database  │   Cache  │ ││
│                    │  │(Postgres)│ (Redis) │  ││
│                    │  └──────────┴──────────┘ ││
│                    └─────────────────────────────┘│
└──────────────────────────────────────────────────────���─���────┘
```

### Level 2: Core Educational Process

```
┌────────────────────────────────────────────────────────────────────────┐
│                 STUDENT MANAGEMENT PROCESS             │
│                                                        │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐      │
│   │ Admission│ ──► │ Enroll   │ ──► │ Class   │      │
│   │   Data   │     │ Student │     │ Assign │      │
│   └──────────┘     └──────────┘     └──────────┘      │
│                                                 │    │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐  │    │
│   │  Exam   │ ◄── │  Grade   │ ◄── │Attendance│  │    │
│   │ Results │     │  Entry   │     │  Mark   │   │    │
│   └──────────┘     └──────────┘     └──────────┘  │    │
│                    │                           │        │
│                    └───────────┬───────────┘        │
│                                ▼                 │
│   ┌─────────────────────────────────────────┐    │
│   │         STUDENT DATABASE                │    │
│   │  - Personal Info   - Academic History │    │
│   │  - Family Data   - Fee Records       │    │
│   │  - Documents    - Attendance      │    │
│   └─────────────────────────────────┘    │
└────────────────────────────────────────────────┘
```

### Data Flows

| Flow ID | From Process | To Process | Data Description |
|--------|-------------|-----------|----------------|
| D1 | Student Portal | Admission | Registration Form |
| D2 | Admission | Student Database | New Student Record |
| D3 | Student Database | Timetable | Student Class Info |
| D4 | Teacher | Attendance | Daily Attendance |
| D5 | Exam System | Report Card | Grade Data |
| D6 | Finance | Payment Gateway | Fee Payment |
| D7 | Analytics | Admin Dashboard | Aggregated Stats |

---

## 🗄 Database Schema

### Core Tables

```sql
-- Tenants (Multi-school)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Users (All roles)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    email CITEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'teacher', 'student', 'parent', 'finance', 'hr')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    two_secret_enabled BOOLEAN DEFAULT false,
    two_secret_secret VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Students
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    student_id VARCHAR(50) UNIQUE NOT NULL,
    gender VARCHAR(10),
    date_of_birth DATE,
    emergency_contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Academic Years
CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT false,
    PRIMARY KEY(tenant_id, name)
);

-- Sections (Classes)
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    academic_year_id UUID REFERENCES academic_years(id),
    name VARCHAR(50) NOT NULL,
    grade_level INT NOT NULL,
    room VARCHAR(50),
    capacity INT DEFAULT 35
);

-- Student Attendance
CREATE TABLE student_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    section_id UUID REFERENCES sections(id),
    date DATE NOT NULL,
    status VARCHAR(20) CHECK (status IN ('present', 'absent', 'late', 'excused')),
    UNIQUE(student_id, date)
);

-- Fee Structures
CREATE TABLE fee_structures (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    due_date DATE
);

-- Student Fees
CREATE TABLE student_fees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    fee_structure_id UUID REFERENCES fee_structures(id),
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'paid', 'overdue'))
);

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    module VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 📡 API Endpoints

### Authentication
```
POST   /api/v1/auth/login              - Login
POST   /api/v1/auth/logout             - Logout
POST   /api/v1/auth/refresh         - Refresh token
POST   /api/v1/auth/login-2fa        - Login with 2FA
GET    /api/v1/2fa/status           - Get 2FA status
GET    /api/v1/2fa/setup            - Setup 2FA (get QR)
POST   /api/v1/2fa/enable           - Enable 2FA
POST   /api/v1/2fa/disable          - Disable 2FA
GET    /api/v1/auth/me             - Current user
```

### Students
```
GET    /api/v1/students             - List students
POST   /api/v1/students            - Create student
GET    /api/v1/students/:id        - Get student
PUT    /api/v1/students/:id       - Update student
GET    /api/v1/students/:id/attendance - Student attendance
GET    /api/v1/students/:id/fees   - Student fees
```

### Academic
```
GET    /api/v1/academic-years      - List years
GET    /api/v1/sections          - List classes
POST   /api/v1/sections           - Create class
GET    /api/v1/subjects           - List subjects
GET    /api/v1/timetables         - Timetable
GET    /api/v1/timetables/generate - Auto-generate
```

### Attendance
```
GET    /api/v1/attendance/students - Student attendance
POST   /api/v1/attendance/students - Mark attendance
GET    /api/v1/attendance/reports - Attendance reports
POST   /api/v1/attendance/qr/scan - QR scan attendance
```

### Finance
```
GET    /api/v1/fees/structures   - Fee structures
POST   /api/v1/fees/assign       - Assign fees
GET    /api/v1/fees/students     - Student fees
POST   /api/v1/fees/students/:id/pay - Process payment
GET    /api/v1/fees/reports      - Financial reports
```

### HR
```
GET    /api/v1/hr/staff           - List staff
POST   /api/v1/hr/staff          - Create staff
GET    /api/v1/hr/leave         - Leave requests
POST   /api/v1/hr/payroll       - Run payroll
```

### Analytics
```
GET    /api/v1/analytics/dashboard - Dashboard data
GET    /api/v1/analytics/kpi     - KPI metrics
POST   /api/v1/ai/chat          - AI Chatbot
```

---

## 🚀 Installation Guide

### Prerequisites

| Software | Version | Purpose |
|----------|---------|---------|
| **Go** | 1.21+ | Backend runtime |
| **Node.js** | 18+ | Frontend runtime |
| **PostgreSQL** | 15+ | Database |
| **Redis** | 7+ | Cache (optional) |
| **Docker** | Latest | Containerization |

### Option 1: Docker Setup (Recommended)

```bash
# Clone the repository
git clone https://github.com/edusyspro/edusys.git
cd edusys

# Start all services
docker-compose up -d

# Access the application
# API: http://localhost:8080
# Web: http://localhost:3000
```

### Option 2: Local Development

#### Backend Setup

```bash
# Navigate to project root
cd edusys

# Copy environment file
cp .env.example .env

# Edit .env with your settings
# Required: DATABASE_URL, JWT_SECRET

# Install Go dependencies
go mod download

# Run migrations (if using local PostgreSQL)
# Create database: createdb edusys

# Start the server
go run ./cmd/server
```

#### Frontend Setup

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start development server
npm run dev

# Access: http://localhost:3000
```

### Environment Variables

```env
# Server
SERVER_PORT=8080
JWT_SECRET=your-super-secret-key
ENVIRONMENT=development
DEBUG=true

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/edusys?sslmode=disable

# Redis (optional)
REDIS_URL=redis://localhost:6379

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

# Payments (Indonesia)
MIDTRANS_KEY=your-midtrans-key
XENDIT_KEY=your-xendit-key
```

### Building for Production

```bash
# Backend
go build -o edusys-server ./cmd/server

# Frontend
cd web
npm run build
```

---

## 📁 Project Structure

```
edusys/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├���─ config/             # Configuration
│   ├── database/          # Database connection
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Auth, CORS, etc.
│   ├── models/            # Data models
│   ├── routes/            # Route definitions
│   ├── totp/              # 2FA implementation
│   └── utils/              # Utilities
├── migrations/             # SQL migrations
├── web/                   # Frontend (React)
│   ├── src/
│   │   ├── components/   # UI components
│   │   ├── pages/        # Page routes
│   │   ├── stores/       # Zustand stores
│   │   ├── lib/          # Utilities
│   │   └── styles/       # CSS
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## 🔐 Security

- ✅ JWT with short-lived access tokens (15 min)
- ✅ Refresh token rotation (30 days)
- ✅ Password hashing with bcrypt (cost 12)
- ✅ Account lockout after 5 failed attempts
- ✅ **Google Authenticator 2FA** support
- ✅ Role-based access control (RBAC)
- ✅ Rate limiting (100 req/min)
- ✅ SQL injection prevention
- ✅ XSS protection
- ✅ TLS 1.3 for all connections
- ✅ Full audit logging
- ✅ Tenant data isolation

---

## 🤖 AI Integration

### Features
- **Smart Chatbot** - NLP-powered school assistant
- **Performance Prediction** - ML-based student risk scoring
- **Anomaly Detection** - Fee/attendance anomalies
- **CV Screening** - Resume parsing & ranking
- **Auto Timetable** - AI scheduling

### API
```
POST /api/v1/ai/chat              - Chat with AI
GET  /api/v1/ai/performance     - Student performance
GET  /api/v1/ai/anomalies       - Anomaly detection
```

---

## 🗓 Roadmap

### v2.0 (Q2 2024)
- [ ] Mobile apps (iOS/Android)
- [ ] Parent portal
- [ ] SMS integration
- [ ] Advanced AI analytics

### v2.1 (Q3 2024)
- [ ] E-learning with video
- [ ] Live classes
- [ ] Certificate builder
- [ ] Document management

### v2.2 (Q4 2024)
- [ ] AI chatbot improvements
- [ ] Predictive analytics
- [ ] Auto timetable engine
- [ ] Blockchain for credentials

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io) - Express-inspired Go web framework
- [shadcn/ui](https://ui.shadcn.com) - Beautiful component library
- [Tailwind CSS](https://tailwindcss.com) - Utility-first CSS framework
- [Recharts](https://recharts.org) - Composable charting library
- [Lucide](https://lucide.dev) - Beautiful icons

---

<p align="center">
  <strong>Edusys Pro</strong> - Modern School Management System<br/>
  Built with ❤️ for educational institutions
</p>"# yciis-dev-vite"  git init git add README.md git commit -m "first commit" git branch -M main git remote add origin https://github.com/2eux/yciis-dev-vite.git git push -u origin main
"# yciis-dev-vite"  git init git add README.md git commit -m "first co
"# yciis-dev-vite"  git init git add README.md git commit -m "second commentttt"
