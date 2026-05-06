// =============================================
// EDUSYS PRO - Convex Helpers & Middleware
// =============================================

import { v } from "convex/values";
import { query, mutation } from "./_generated/server";

export async function requireAuth(ctx: any) {
  const identity = await ctx.auth.getUserIdentity();
  if (!identity) {
    throw new Error("Authentication required");
  }
  return identity;
}

export async function getCurrentUser(ctx: any) {
  const identity = await ctx.auth.getUserIdentity();
  if (!identity) throw new Error("Not authenticated");

  const profile = await ctx.db
    .query("userProfiles")
    .withIndex("by_userId", (q: any) => q.eq("userId", identity.subject))
    .first();

  return { identity, profile };
}

export async function getCurrentTenant(ctx: any) {
  const { profile } = await getCurrentUser(ctx);
  if (!profile?.tenantId) {
    throw new Error("No tenant assigned");
  }
  return profile.tenantId;
}

export async function requireRole(ctx: any, allowedRoles: string[]) {
  const { profile } = await getCurrentUser(ctx);
  if (!profile || !allowedRoles.includes(profile.role)) {
    throw new Error("Insufficient permissions");
  }
  return profile;
}

export function sanitizePhone(phone: string): string {
  return phone.replace(/[^0-9+]/g, "");
}

export function generateStudentId(year: number, count: number): string {
  return `STU${year}${String(count + 1).padStart(4, "0")}`;
}

export function formatCurrency(amount: number, currency: string = "IDR"): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amount);
}

// =============================================
// DASHBOARD QUERIES
// =============================================

export const getDashboardStats = query({
  args: {},
  handler: async (ctx) => {
    const tenantId = await getCurrentTenant(ctx);

    const students = await ctx.db
      .query("students")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .collect();

    const staff = await ctx.db
      .query("staff")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .collect();

    const sections = await ctx.db
      .query("sections")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .collect();

    const activeStudents = students.filter(s => s.status === "active").length;
    const activeStaff = staff.filter(s => s.status === "active").length;

    return {
      totalStudents: activeStudents,
      totalStaff: activeStaff,
      totalSections: sections.length,
      totalTeachers: staff.filter(s => s.department === "Teaching").length,
    };
  },
});

export const getTodayAttendance = query({
  args: { date: v.string() },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);

    const records = await ctx.db
      .query("attendance")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .filter((q) => q.eq(q.field("date"), args.date))
      .collect();

    const total = records.length;
    const present = records.filter(r => r.status === "present").length;
    const absent = records.filter(r => r.status === "absent").length;
    const late = records.filter(r => r.status === "late").length;
    const percentage = total > 0 ? Math.round(((present + late) / total) * 100) : 0;

    return { present, absent, late, total, percentage };
  },
});

export const getFeeSummary = query({
  args: {},
  handler: async (ctx) => {
    const tenantId = await getCurrentTenant(ctx);

    const fees = await ctx.db
      .query("studentFees")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .collect();

    const collected = fees
      .filter(f => f.status === "paid")
      .reduce((sum, f) => sum + f.paidAmount, 0);

    const pending = fees
      .filter(f => f.status === "pending")
      .reduce((sum, f) => sum + (f.finalAmount - f.paidAmount), 0);

    const overdue = fees
      .filter(f => f.status === "overdue")
      .reduce((sum, f) => sum + (f.finalAmount - f.paidAmount), 0);

    return { collected, pending, overdue };
  },
});

export const getRecentActivities = query({
  args: { limit: v.optional(v.float64()) },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    const limit = args.limit || 10;

    return await ctx.db
      .query("auditLogs")
      .withIndex("by_tenant", (q) => q.eq("tenantId", tenantId))
      .order("desc")
      .take(limit);
  },
});