// =============================================
// EDUSYS PRO - Convex Students API
// =============================================

import { v } from "convex/values";
import { query, mutation } from "./_generated/server";
import { requireAuth, getCurrentTenant } from "./helpers";

// =============================================
// QUERIES
// =============================================

export const list = query({
  args: {
    search: v.optional(v.string()),
    sectionId: v.optional(v.id("sections")),
    status: v.optional(v.string()),
    limit: v.optional(v.float64()),
    page: v.optional(v.float64()),
  },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) throw new Error("Not authenticated");

    const limit = args.limit || 20;
    const page = args.page || 0;
    const skip = page * limit;

    let query = ctx.db.query("students")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId));

    if (args.status) {
      query = ctx.db.query("students")
        .withIndex("by_status", (q) => q.eq("status", args.status));
    }

    const students = await query.order("desc").collect();

    if (args.search) {
      const search = args.search.toLowerCase();
      return students
        .filter(s =>
          (s.firstName?.toLowerCase().includes(search)) ||
          (s.lastName?.toLowerCase().includes(search)) ||
          (s.studentId?.toLowerCase().includes(search))
        )
        .slice(skip, skip + limit);
    }

    if (args.sectionId) {
      const filtered = students.filter(s => s.sectionId === args.sectionId);
      return filtered.slice(skip, skip + limit);
    }

    return students.slice(skip, skip + limit);
  },
});

export const getById = query({
  args: { id: v.id("students") },
  handler: async (ctx, args) => {
    const student = await ctx.db.get(args.id);
    if (!student) throw new Error("Student not found");
    return student;
  },
});

export const getByStudentId = query({
  args: { studentId: v.string() },
  handler: async (ctx, args) => {
    const student = await ctx.db
      .query("students")
      .withIndex("by_studentId", (q) => q.eq("studentId", args.studentId))
      .first();
    
    if (!student) throw new Error("Student not found");
    return student;
  },
});

export const getAttendance = query({
  args: {
    studentId: v.id("students"),
    academicYearId: v.id("academicYears"),
  },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("attendance")
      .withIndex("by_student", (q) => q.eq("studentId", args.studentId))
      .filter((q) => q.eq(q.field("academicYearId"), args.academicYearId))
      .order("desc")
      .collect();
  },
});

export const getFees = query({
  args: { studentId: v.id("students") },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("studentFees")
      .withIndex("by_student", (q) => q.eq("studentId", args.studentId))
      .order("desc")
      .collect();
  },
});

export const getFamily = query({
  args: { studentId: v.id("students") },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("studentParents")
      .withIndex("by_student", (q) => q.eq("studentId", args.studentId))
      .collect();
  },
});

// =============================================
// MUTATIONS
// =============================================

export const create = mutation({
  args: {
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
    sectionId: v.optional(v.id("sections")),
  },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) throw new Error("Not authenticated");

    // Generate student ID
    const year = new Date().getFullYear();
    const count = (await ctx.db.query("students")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .collect()).length;
    
    const studentId = `STU${year}${String(count + 1).padStart(4, "0")}`;

    const student = await ctx.db.insert("students", {
      tenantId,
      studentId,
      firstName: args.firstName,
      lastName: args.lastName,
      gender: args.gender,
      dateOfBirth: args.dateOfBirth,
      placeOfBirth: args.placeOfBirth,
      nationality: args.nationality,
      religion: args.religion,
      bloodType: args.bloodType,
      address: args.address,
      city: args.city,
      province: args.province,
      postalCode: args.postalCode,
      emergencyContactName: args.emergencyContactName,
      emergencyContactPhone: args.emergencyContactPhone,
      emergencyContactRelation: args.emergencyContactRelation,
      sectionId: args.sectionId,
      status: "active",
      rollNumber: String(count + 1).padStart(3, "0"),
    });

    // Add audit log
    await ctx.db.insert("auditLogs", {
      tenantId,
      userId: identity.subject,
      action: "CREATE",
      module: "students",
      entityType: "students",
      entityId: student,
      newValues: { ...args, studentId },
      ipAddress: "",
      userAgent: "",
    });

    return student;
  },
});

export const update = mutation({
  args: {
    id: v.id("students"),
    firstName: v.optional(v.string()),
    lastName: v.optional(v.string()),
    gender: v.optional(v.string()),
    dateOfBirth: v.optional(v.string()),
    address: v.optional(v.string()),
    city: v.optional(v.string()),
    province: v.optional(v.string()),
    postalCode: v.optional(v.string()),
    emergencyContactName: v.optional(v.string()),
    emergencyContactPhone: v.optional(v.string()),
    emergencyContactRelation: v.optional(v.string()),
    notes: v.optional(v.string()),
    status: v.optional(v.string()),
    sectionId: v.optional(v.id("sections")),
  },
  handler: async (ctx, args) => {
    const { id, ...updates } = args;
    await ctx.db.patch(id, updates);
    return await ctx.db.get(id);
  },
});

export const remove = mutation({
  args: { id: v.id("students") },
  handler: async (ctx, args) => {
    await ctx.db.patch(args.id, { status: "inactive" });
    return { success: true };
  },
});

export const addFamilyMember = mutation({
  args: {
    studentId: v.id("students"),
    relation: v.string(),
    firstName: v.string(),
    lastName: v.optional(v.string()),
    email: v.optional(v.string()),
    phone: v.optional(v.string()),
    occupation: v.optional(v.string()),
    company: v.optional(v.string()),
    incomeBracket: v.optional(v.string()),
    isPrimary: v.boolean(),
  },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    return await ctx.db.insert("studentParents", {
      tenantId,
      studentId: args.studentId,
      relation: args.relation,
      firstName: args.firstName,
      lastName: args.lastName,
      email: args.email,
      phone: args.phone,
      occupation: args.occupation,
      company: args.company,
      incomeBracket: args.incomeBracket,
      isPrimary: args.isPrimary,
    });
  },
});