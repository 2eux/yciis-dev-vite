// =============================================
// EDUSYS PRO - Convex Schema
// =============================================
// Convex uses a document-based schema with
// tables defined by TypeScript types and
// validated with the `v` validator.
// =============================================

import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  // =============================================
  // TENANTS (Schools/Organizations)
  // =============================================
  tenants: defineTable({
    name: v.string(),
    code: v.string(),
    timezone: v.string(),
    currencyCode: v.string(),
    logoUrl: v.optional(v.string()),
    address: v.optional(v.string()),
    phone: v.optional(v.string()),
    email: v.optional(v.string()),
    isActive: v.boolean(),
    subscriptionPlan: v.string(),
    subscriptionExpiresAt: v.optional(v.number()),
    settings: v.optional(v.any()),
  })
    .index("by_code", ["code"])
    .index("by_active", ["isActive"]),

  // =============================================
  // USER PROFILES
  // =============================================
  userProfiles: defineTable({
    userId: v.string(),
    tenantId: v.optional(v.id("tenants")),
    role: v.union(
      v.literal("super_admin"),
      v.literal("admin"),
      v.literal("teacher"),
      v.literal("student"),
      v.literal("parent"),
      v.literal("finance"),
      v.literal("hr")
    ),
    firstName: v.string(),
    lastName: v.optional(v.string()),
    phone: v.optional(v.string()),
    avatarUrl: v.optional(v.string()),
    isActive: v.boolean(),
    twoFAEnabled: v.boolean(),
    twoFASecret: v.optional(v.string()),
    lastLoginAt: v.optional(v.number()),
    lastLoginIp: v.optional(v.string()),
  })
    .index("by_userId", ["userId"])
    .index("by_tenant", ["tenantId"])
    .index("by_role", ["role"]),

  // =============================================
  // ACADEMIC YEARS
  // =============================================
  academicYears: defineTable({
    tenantId: v.id("tenants"),
    name: v.string(),
    startDate: v.string(),
    endDate: v.string(),
    isCurrent: v.boolean(),
    isActive: v.boolean(),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_current", ["tenantId", "isCurrent"]),

  // =============================================
  // SECTIONS (Classes)
  // =============================================
  sections: defineTable({
    tenantId: v.id("tenants"),
    academicYearId: v.id("academicYears"),
    name: v.string(),
    gradeLevel: v.float64(),
    room: v.optional(v.string()),
    capacity: v.float64(),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_year", ["academicYearId"]),

  // =============================================
  // STUDENTS
  // =============================================
  students: defineTable({
    tenantId: v.id("tenants"),
    userId: v.optional(v.string()),
    studentId: v.string(),
    firstName: v.string(),
    lastName: v.string(),
    gender: v.optional(v.string()),
    dateOfBirth: v.optional(v.string()),
    placeOfBirth: v.optional(v.string()),
    nationality: v.optional(v.string()),
    religion: v.optional(v.string()),
    bloodType: v.optional(v.string()),
    address: v.optional(v.string()),
    city: v.optional(v.string()),
    province: v.optional(v.string()),
    postalCode: v.optional(v.string()),
    emergencyContactName: v.optional(v.string()),
    emergencyContactPhone: v.optional(v.string()),
    emergencyContactRelation: v.optional(v.string()),
    notes: v.optional(v.string()),
    documents: v.optional(v.any()),
    healthInfo: v.optional(v.any()),
    status: v.string(),
    sectionId: v.optional(v.id("sections")),
    rollNumber: v.optional(v.string()),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_studentId", ["studentId"])
    .index("by_section", ["sectionId"])
    .index("by_user", ["userId"])
    .index("by_status", ["status"]),

  // =============================================
  // STUDENT PARENTS
  // =============================================
  studentParents: defineTable({
    tenantId: v.id("tenants"),
    studentId: v.id("students"),
    userId: v.optional(v.string()),
    relation: v.string(),
    firstName: v.string(),
    lastName: v.optional(v.string()),
    email: v.optional(v.string()),
    phone: v.optional(v.string()),
    occupation: v.optional(v.string()),
    company: v.optional(v.string()),
    incomeBracket: v.optional(v.string()),
    isPrimary: v.boolean(),
  })
    .index("by_student", ["studentId"])
    .index("by_user", ["userId"]),

  // =============================================
  // SUBJECTS
  // =============================================
  subjects: defineTable({
    tenantId: v.id("tenants"),
    name: v.string(),
    code: v.string(),
    description: v.optional(v.string()),
    isOptional: v.boolean(),
    isActive: v.boolean(),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_code", ["tenantId", "code"]),

  // =============================================
  // SUBJECT ASSIGNMENTS
  // =============================================
  sectionSubjects: defineTable({
    tenantId: v.id("tenants"),
    sectionId: v.id("sections"),
    subjectId: v.id("subjects"),
    academicYearId: v.id("academicYears"),
    teacherId: v.optional(v.id("userProfiles")),
  })
    .index("by_section", ["sectionId"])
    .index("by_subject", ["subjectId"])
    .index("by_teacher", ["teacherId"]),

  // =============================================
  // TIMETABLES
  // =============================================
  timetables: defineTable({
    tenantId: v.id("tenants"),
    sectionId: v.id("sections"),
    sectionSubjectId: v.id("sectionSubjects"),
    teacherId: v.optional(v.id("userProfiles")),
    dayOfWeek: v.float64(),
    periodStart: v.float64(),
    periodEnd: v.float64(),
    room: v.optional(v.string()),
    academicYearId: v.id("academicYears"),
  })
    .index("by_section", ["sectionId"])
    .index("by_teacher", ["teacherId"]),

  // =============================================
  // ATTENDANCE
  // =============================================
  attendance: defineTable({
    tenantId: v.id("tenants"),
    studentId: v.id("students"),
    sectionId: v.id("sections"),
    academicYearId: v.id("academicYears"),
    date: v.string(),
    status: v.union(
      v.literal("present"),
      v.literal("absent"),
      v.literal("late"),
      v.literal("excused")
    ),
    timeIn: v.optional(v.string()),
    timeOut: v.optional(v.string()),
    remarks: v.optional(v.string()),
    markedBy: v.optional(v.id("userProfiles")),
    deviceId: v.optional(v.string()),
  })
    .index("by_student", ["studentId"])
    .index("by_section_date", ["sectionId", "date"])
    .index("by_student_date", ["studentId", "date"]),

  // =============================================
  // EXAMS
  // =============================================
  exams: defineTable({
    tenantId: v.id("tenants"),
    academicYearId: v.id("academicYears"),
    name: v.string(),
    type: v.union(
      v.literal("quiz"),
      v.literal("midterm"),
      v.literal("final"),
      v.literal("semester"),
      v.literal("annual")
    ),
    sectionId: v.id("sections"),
    subjectId: v.id("subjects"),
    startDate: v.string(),
    endDate: v.string(),
    durationMinutes: v.optional(v.float64()),
    totalMarks: v.optional(v.float64()),
    passingMarks: v.optional(v.float64()),
    instructions: v.optional(v.string()),
    isPublished: v.boolean(),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_section", ["sectionId"])
    .index("by_subject", ["subjectId"]),

  // =============================================
  // EXAM MARKS
  // =============================================
  examMarks: defineTable({
    tenantId: v.id("tenants"),
    studentId: v.id("students"),
    examId: v.id("exams"),
    subjectId: v.id("subjects"),
    marksObtained: v.float64(),
    totalMarks: v.float64(),
    grade: v.optional(v.string()),
    remarks: v.optional(v.string()),
    gradedBy: v.optional(v.id("userProfiles")),
    gradedAt: v.optional(v.number()),
  })
    .index("by_exam", ["examId"])
    .index("by_student_exam", ["studentId", "examId"]),

  // =============================================
  // FEE STRUCTURES
  // =============================================
  feeStructures: defineTable({
    tenantId: v.id("tenants"),
    academicYearId: v.id("academicYears"),
    name: v.string(),
    description: v.optional(v.string()),
    amount: v.float64(),
    dueDate: v.optional(v.string()),
    isRecurring: v.boolean(),
    frequency: v.optional(v.string()),
    isActive: v.boolean(),
  })
    .index("by_tenant", ["tenantId"]),

  // =============================================
  // STUDENT FEES
  // =============================================
  studentFees: defineTable({
    tenantId: v.id("tenants"),
    studentId: v.id("students"),
    feeStructureId: v.id("feeStructures"),
    academicYearId: v.id("academicYears"),
    amount: v.float64(),
    discountAmount: v.float64(),
    finalAmount: v.float64(),
    status: v.union(
      v.literal("pending"),
      v.literal("partial"),
      v.literal("paid"),
      v.literal("overdue"),
      v.literal("waived")
    ),
    dueDate: v.optional(v.string()),
    paidAmount: v.float64(),
    paidAt: v.optional(v.number()),
    paymentMethod: v.optional(v.string()),
    transactionId: v.optional(v.string()),
  })
    .index("by_student", ["studentId"])
    .index("by_status", ["status"])
    .index("by_due", ["dueDate", "status"]),

  // =============================================
  // STAFF
  // =============================================
  staff: defineTable({
    tenantId: v.id("tenants"),
    userId: v.optional(v.string()),
    employeeId: v.string(),
    department: v.optional(v.string()),
    position: v.optional(v.string()),
    joinDate: v.optional(v.string()),
    employmentType: v.optional(v.string()),
    status: v.string(),
    salary: v.optional(v.float64()),
    bankAccount: v.optional(v.string()),
    bankName: v.optional(v.string()),
    emergencyContactName: v.optional(v.string()),
    emergencyContactPhone: v.optional(v.string()),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_employeeId", ["employeeId"])
    .index("by_department", ["department"]),

  // =============================================
  // LEAVE REQUESTS
  // =============================================
  leaveRequests: defineTable({
    tenantId: v.id("tenants"),
    staffId: v.id("staff"),
    leaveType: v.string(),
    startDate: v.string(),
    endDate: v.string(),
    totalDays: v.float64(),
    reason: v.optional(v.string()),
    status: v.string(),
    approvedBy: v.optional(v.id("userProfiles")),
    approvedAt: v.optional(v.number()),
    remarks: v.optional(v.string()),
  })
    .index("by_staff", ["staffId"])
    .index("by_status", ["status"]),

  // =============================================
  // LMS COURSES
  // =============================================
  lmsCourses: defineTable({
    tenantId: v.id("tenants"),
    teacherId: v.optional(v.id("userProfiles")),
    academicYearId: v.optional(v.id("academicYears")),
    title: v.string(),
    description: v.optional(v.string()),
    thumbnailUrl: v.optional(v.string()),
    isPublished: v.boolean(),
    isFree: v.boolean(),
    price: v.float64(),
    language: v.string(),
    difficulty: v.string(),
    durationHours: v.optional(v.float64()),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_teacher", ["teacherId"]),

  // =============================================
  // LMS SECTIONS
  // =============================================
  lmsSections: defineTable({
    courseId: v.id("lmsCourses"),
    title: v.string(),
    description: v.optional(v.string()),
    sortOrder: v.float64(),
  })
    .index("by_course", ["courseId"]),

  // =============================================
  // LMS CONTENT
  // =============================================
  lmsContent: defineTable({
    sectionId: v.id("lmsSections"),
    title: v.string(),
    contentType: v.union(
      v.literal("video"),
      v.literal("document"),
      v.literal("quiz"),
      v.literal("assignment"),
      v.literal("link")
    ),
    contentUrl: v.optional(v.string()),
    contentText: v.optional(v.string()),
    durationMinutes: v.optional(v.float64()),
    sortOrder: v.float64(),
    isPreview: v.boolean(),
  })
    .index("by_section", ["sectionId"]),

  // =============================================
  // LMS ENROLLMENTS
  // =============================================
  lmsEnrollments: defineTable({
    tenantId: v.id("tenants"),
    courseId: v.id("lmsCourses"),
    studentId: v.id("students"),
    progressPercentage: v.float64(),
    completedAt: v.optional(v.number()),
    enrolledAt: v.number(),
  })
    .index("by_student", ["studentId"])
    .index("by_course", ["courseId"])
    .index("by_student_course", ["studentId", "courseId"]),

  // =============================================
  // BOOKS
  // =============================================
  books: defineTable({
    tenantId: v.id("tenants"),
    isbn: v.optional(v.string()),
    title: v.string(),
    author: v.optional(v.string()),
    publisher: v.optional(v.string()),
    publicationYear: v.optional(v.float64()),
    category: v.optional(v.string()),
    location: v.optional(v.string()),
    totalCopies: v.float64(),
    availableCopies: v.float64(),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_isbn", ["isbn"]),

  // =============================================
  // BOOK ISSUES
  // =============================================
  bookIssues: defineTable({
    tenantId: v.id("tenants"),
    bookId: v.id("books"),
    studentId: v.optional(v.id("students")),
    staffId: v.optional(v.id("staff")),
    issueDate: v.string(),
    dueDate: v.string(),
    returnDate: v.optional(v.string()),
    status: v.string(),
  })
    .index("by_book", ["bookId"])
    .index("by_student", ["studentId"]),

  // =============================================
  // ADMISSION LEADS
  // =============================================
  admissionLeads: defineTable({
    tenantId: v.id("tenants"),
    firstName: v.string(),
    lastName: v.optional(v.string()),
    email: v.optional(v.string()),
    phone: v.string(),
    gender: v.optional(v.string()),
    dateOfBirth: v.optional(v.string()),
    gradeApplied: v.optional(v.float64()),
    source: v.optional(v.string()),
    status: v.string(),
    assignedTo: v.optional(v.id("userProfiles")),
    notes: v.optional(v.string()),
    documents: v.optional(v.any()),
    followUpAt: v.optional(v.number()),
    convertedAt: v.optional(v.number()),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_status", ["status"]),

  // =============================================
  // NOTIFICATIONS
  // =============================================
  notifications: defineTable({
    tenantId: v.id("tenants"),
    userId: v.string(),
    title: v.string(),
    body: v.string(),
    type: v.string(),
    data: v.optional(v.any()),
    isRead: v.boolean(),
    readAt: v.optional(v.number()),
  })
    .index("by_user", ["userId", "isRead"])
    .index("by_tenant", ["tenantId"]),

  // =============================================
  // AUDIT LOGS
  // =============================================
  auditLogs: defineTable({
    tenantId: v.optional(v.id("tenants")),
    userId: v.optional(v.string()),
    action: v.string(),
    module: v.string(),
    entityType: v.optional(v.string()),
    entityId: v.optional(v.id("students")),
    oldValues: v.optional(v.any()),
    newValues: v.optional(v.any()),
    ipAddress: v.optional(v.string()),
    userAgent: v.optional(v.string()),
  })
    .index("by_tenant", ["tenantId"])
    .index("by_user", ["userId"])
    .index("by_action", ["action"]),

  // =============================================
  // TRANSPORT ROUTES
  // =============================================
  transportRoutes: defineTable({
    tenantId: v.id("tenants"),
    name: v.string(),
    startPoint: v.string(),
    endPoint: v.string(),
    waypoints: v.optional(v.any()),
    distanceKm: v.optional(v.float64()),
    estimatedTimeMinutes: v.optional(v.float64()),
  })
    .index("by_tenant", ["tenantId"]),

  // =============================================
  // VEHICLES
  // =============================================
  vehicles: defineTable({
    tenantId: v.id("tenants"),
    vehicleNumber: v.string(),
    vehicleType: v.string(),
    model: v.optional(v.string()),
    capacity: v.optional(v.float64()),
    driverName: v.optional(v.string()),
    driverPhone: v.optional(v.string()),
    insuranceExpiry: v.optional(v.string()),
    fitnessExpiry: v.optional(v.string()),
    isActive: v.boolean(),
  })
    .index("by_tenant", ["tenantId"]),

  // =============================================
  // BOARDING LOGS
  // =============================================
  boardingLogs: defineTable({
    tenantId: v.id("tenants"),
    studentId: v.id("students"),
    routeId: v.id("transportRoutes"),
    date: v.string(),
    pickupTime: v.optional(v.string()),
    pickupLocation: v.optional(v.string()),
    dropoffTime: v.optional(v.string()),
    dropoffLocation: v.optional(v.string()),
    status: v.string(),
  })
    .index("by_student", ["studentId"])
    .index("by_route", ["routeId"]),
});