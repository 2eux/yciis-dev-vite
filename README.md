# Edusys Pro - Enterprise School ERP System

<p align="center">
  <img src="https://img.shields.io/badge/Version-1.0.0-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Convex-8B5CF6?style=for-the-badge&logo=convex" />
  <img src="https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react" />
</p>

Enterprise-grade School ERP powered by **Convex** (serverless backend) and **React**. Real-time, multi-tenant, and AI-ready.

---

## 📋 Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [System Architecture](#system-architecture)
- [Data Flow Diagrams](#data-flow-diagrams)
- [Database Schema](#database-schema)
- [Convex API](#convex-api)
- [Installation](#installation)
- [Project Structure](#project-structure)
- [Deployment](#deployment)
- [Security](#security)
- [Roadmap](#roadmap)

---

## ✨ Features

| Module | Features |
|--------|----------|
| **Auth** | JWT RBAC, Google OAuth, Google Authenticator 2FA, audit logs |
| **Students** | Master data, family/guardians, documents, health notes, class mapping |
| **Academic** | Subjects, sections, curriculum, timetable engine, teacher assignments |
| **Attendance** | Daily student marks, QR/biometric, real-time parent push |
| **Admission CRM** | Lead → Prospect → Student pipeline, online forms, approvals |
| **Exams** | Builder, question bank, grading, report cards (PDF) |
| **LMS** | Course builder, video/content, assignments, quizzes, progress |
| **Finance** | Fee structures, invoicing, online payments, reports |
| **HR** | Staff profiles, leave management, payroll automation |
| **Library** | Book catalog, issue/return, stock tracking |
| **Transport** | Routes, vehicles, GPS, student boarding logs |
| **Analytics** | Dashboard, KPIs, AI predictions, anomaly detection |

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Convex (serverless: DB, Auth, Functions, Real-time, File Storage) |
| **Frontend** | React 18 + TypeScript + Vite |
| **UI** | Tailwind CSS + shadcn/ui |
| **State** | Zustand |
| **Charts** | Recharts |
| **Icons** | Lucide React |
| **Deploy** | Docker + Coolify |

---

## 🏗 System Architecture

```
React SPA (Vite) ──► Convex React Client ──► Convex Cloud
                        │
                        ├── Queries (read data, real-time subscriptions)
                        ├── Mutations (write data, ACID transactions)
                        ├── Auth (JWT, Google OAuth, magic link)
                        ├── File Storage (images, documents)
                        └── Scheduled Functions (cron jobs)

Optional: Go API Gateway (port 8080) for rate limiting / custom logic
```

**Key:** Convex replaces PostgreSQL, Redis, S3, REST server, and WebSocket server — all in one platform.

---

## 📊 Data Flow Diagrams

### Level 0 - Context Diagram

```
  ┌──────────┐     ┌──────────────┐     ┌──────────┐
  │ Students │────►│              │────►│ Reports  │
  └──────────┘     │   EDUSYS     │     └──────────┘
  ┌──────────┐     │     PRO      │     ┌──────────┐
  │ Teachers │────►│  (Convex)    │────►│ Analytics│
  └──────────┘     │              │     └──────────┘
  ┌──────────┐     │              │     ┌──────────┐
  │ Parents  │────►│              │────►│ Notices  │
  └──────────┘     └──────────────┘     └──────────┘
```

### Level 1 - Process Flow

```
Admission ──► Enrollment ──► Class Assignment ──► Timetable
                                                    │
Attendance ◄────────────────────────────────────────┘
    │
    ▼
Grading ──► Report Cards ──► Parent Portal
                                │
Fee Management ◄───────────────┘
    │
    ▼
Payment Gateway ──► Financial Reports
```

### Data Flows

| ID | From | To | Description |
|----|------|-----|-------------|
| D1 | Student Portal | Admission | Registration form |
| D2 | Admission | Students | New student record |
| D3 | Teacher | Attendance | Daily marks |
| D4 | Exam System | Report Card | Grade data |
| D5 | Finance | Payment Gateway | Fee transaction |
| D6 | Analytics | Dashboard | Aggregated stats |
| D7 | Scheduler | Notifications | Push alerts |

---

## 🗄 Database Schema

Convex uses a **document-based schema** defined in `convex/schema.ts`. Full schema: 25+ tables.

### Core Tables

```typescript
// Tenants (multi-school isolation)
tenants: defineTable({
  name: v.string(),
  code: v.string(),
  timezone: v.string(),
  isActive: v.boolean(),
}).index("by_code", ["code"]),

// User profiles linked to Convex Auth
userProfiles: defineTable({
  userId: v.string(),         // Convex Auth user ID
  tenantId: v.id("tenants"),
  role: v.union(v.literal("admin"), v.literal("teacher"), ...),
  firstName: v.string(),
  twoFAEnabled: v.boolean(),
}).index("by_userId", ["userId"]),

// Students
students: defineTable({
  tenantId: v.id("tenants"),
  studentId: v.string(),       // Auto-generated: STU20250001
  firstName: v.string(),
  lastName: v.string(),
  sectionId: v.id("sections"),
  status: v.string(),
}).index("by_studentId", ["studentId"])
 .index("by_section", ["sectionId"]),

// Attendance (realtime-enabled)
attendance: defineTable({
  studentId: v.id("students"),
  sectionId: v.id("sections"),
  date: v.string(),
  status: v.union(v.literal("present"), v.literal("absent"), ...),
}).index("by_student_date", ["studentId", "date"])
 .index("by_section_date", ["sectionId", "date"]),

// Fees with payment tracking
studentFees: defineTable({
  studentId: v.id("students"),
  amount: v.float64(),
  status: v.union(v.literal("pending"), v.literal("paid"), ...),
  paidAmount: v.float64(),
}).index("by_student", ["studentId"])
 .index("by_status", ["status"]),
```

Full schema: `convex/schema.ts` — 25 tables with indexes and validators.

---

## 📡 Convex API

Convex replaces REST endpoints with **type-safe queries and mutations** that are called directly from the frontend with zero API glue code.

### Students

```typescript
// Query: List students
const students = useQuery(api.students.list, { status: "active" })

// Mutation: Create student
const createStudent = useMutation(api.students.create)
await createStudent({ firstName: "John", lastName: "Doe", sectionId: "..." })

// Query: Get student by ID
const student = useQuery(api.students.getById, { id: studentId })

// Mutation: Update student
const updateStudent = useMutation(api.students.update)
await updateStudent({ id: studentId, address: "New address" })

// Mutation: Remove (soft delete → status: inactive)
const removeStudent = useMutation(api.students.remove)
await removeStudent({ id: studentId })
```

### Attendance

```typescript
// Query: Get section attendance for a date
const attendance = useQuery(api.attendance.listBySection, { sectionId, date })

// Mutation: Mark single or bulk attendance
const markAttendance = useMutation(api.attendance.mark)
await markAttendance({ records: [{ studentId, sectionId, date, status: "present" }] })

// Query: Attendance summary with percentages
const summary = useQuery(api.attendance.getSummary, { studentId, academicYearId })
// Returns: { present: 180, absent: 10, late: 5, total: 195, percentage: 95 }
```

### Dashboard

```typescript
// Query: School overview stats
const stats = useQuery(api.helpers.getDashboardStats, {})
// Returns: { totalStudents, totalStaff, totalSections, totalTeachers }

// Query: Today's attendance
const today = useQuery(api.helpers.getTodayAttendance, { date: "2025-04-23" })
// Returns: { present, absent, late, total, percentage }

// Query: Fee collection summary
const fees = useQuery(api.helpers.getFeeSummary, {})
// Returns: { collected, pending, overdue }

// Query: Recent activity log
const activities = useQuery(api.helpers.getRecentActivities, { limit: 10 })
```

### Real-time Subscriptions

All Convex queries are **automatically real-time** — the UI updates instantly when data changes, with zero extra code.

---

## 🚀 Installation

### Prerequisites

- **Node.js** 18+
- **Convex account** (free tier at [convex.dev](https://convex.dev))

### Setup

```bash
# 1. Clone and install
git clone https://github.com/edusyspro/edusys.git
cd edusys/web
npm install

# 2. Initialize Convex (one-time)
cd ..                          # back to project root
npx convex dev                 # creates convex.json + _generated/ types

# 3. Start both Convex + Vite
# Convex dev server: auto-started by step 2
# Vite dev server:
cd web
npm run dev
```

Access: **http://localhost:5173**

### Environment Variables

```env
# .env (auto-created by `npx convex dev`)
VITE_CONVEX_URL=https://your-project.convex.cloud

# Optional: Google OAuth
VITE_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com

# Optional: Payments (Indonesia)
MIDTRANS_CLIENT_KEY=
MIDTRANS_SERVER_KEY=
XENDIT_API_KEY=
```

---

## 📁 Project Structure

```
edusys/
├── convex/                    # Convex backend
│   ├── schema.ts              # 25+ tables with indexes
│   ├── students.ts            # Student queries & mutations
│   ├── attendance.ts          # Attendance queries & mutations
│   ├── helpers.ts             # Auth, dashboard, utilities
│   └── _generated/            # Auto-generated (npx convex dev)
│       └── api.ts             # Type-safe API
│
├── web/                       # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   └── layout/        # Layout (sticky header nav)
│   │   ├── pages/
│   │   │   ├── Login.tsx      # Auth page
│   │   │   ├── Dashboard.tsx  # Admin dashboard
│   │   │   ├── Students.tsx   # Student management
│   │   │   ├── Academic.tsx   # Academic module
│   │   │   ├── Attendance.tsx # Attendance tracker
│   │   │   ├── Exams.tsx      # Exam builder
│   │   │   ├── Fees.tsx       # Finance
│   │   │   ├── HR.tsx         # Staff management
│   │   │   ├── LMS.tsx        # Learning management
│   │   │   ├── Library.tsx    # Book tracking
│   │   │   ├── Transport.tsx  # Routes & vehicles
│   │   │   └── Settings.tsx   # System settings
│   │   ├── stores/
│   │   │   └── convex-auth.ts # Auth state (Zustand)
│   │   ├── lib/
│   │   │   ├── convex-hooks.ts # Typed hooks
│   │   │   └── utils.ts       # Utility functions
│   │   ├── App.tsx             # Root (ConvexProvider + Router)
│   │   ├── main.tsx            # Entry point
│   │   └── index.css           # Tailwind + shadcn/ui design system
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
│
├── coolify/
│   ├── compose.yaml           # Coolify deployment
│   └── Dockerfile             # Go API (optional)
│
├── docker-compose.yml         # Local Docker
├── .env.example
├── REVIEW.md
├── SPEC.md
├── SECURITY.md
└── README.md
```

---

## 🐳 Deployment

### Docker (Local)

```bash
docker-compose up -d
# Web: http://localhost:3000
```

### Convex Cloud (Production)

```bash
npx convex deploy              # Deploy backend to production
cd web && npm run build        # Build frontend
```

### Coolify

1. Deploy Convex backend: `npx convex deploy`
2. Set `VITE_CONVEX_URL` to your production Convex URL
3. Deploy `coolify/compose.yaml` on Coolify

---

## 🔐 Security

| Feature | Implementation |
|---------|---------------|
| **Authentication** | Convex Auth with JWT + Google OAuth |
| **2FA** | Google Authenticator (TOTP) — `convex/totp.ts` |
| **RBAC** | Row-level via Convex auth checks in queries/mutations |
| **Multi-tenant** | tenantId filtering in every query |
| **Audit Logs** | All mutations write to `auditLogs` table |
| **Rate Limiting** | Handled by Convex platform |
| **HTTPS** | Enforced by Convex Cloud |
| **Input Validation** | Convex schema validators (`v.string()`, `v.union()`) |

---

## 🤖 AI Features (Roadmap)

- Smart chatbot for parent/student queries
- Student performance prediction (ML)
- Auto report generation
- Fee/attendance anomaly detection
- CV screening for HR recruitment

---

## 🗓 Roadmap

| Version | Features |
|---------|----------|
| **v1.0** (Current) | Auth, Students, Attendance, Dashboard |
| **v1.5** | Exams, Fees, HR, Library |
| **v2.0** | LMS, Transport, Analytics |
| **v2.5** | Mobile apps, AI chatbot, WhatsApp |
| **v3.0** | Predictive analytics, auto timetable |

---

## 📄 License

MIT — see [LICENSE](LICENSE)

---

<p align="center">
  <strong>Edusys Pro</strong> — Modern School ERP<br/>
  Built with React + Convex
</p>"# yciis-convex" 
"# yciis-convex" 
"# yciis-convex" 
