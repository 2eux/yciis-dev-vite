// =============================================
// EDUSYS PRO - Convex Attendance API
// =============================================

import { v } from "convex/values";
import { query, mutation } from "./_generated/server";
import { requireAuth, getCurrentTenant } from "./helpers";

// =============================================
// QUERIES
// =============================================

export const listBySection = query({
  args: {
    sectionId: v.id("sections"),
    date: v.string(),
  },
  handler: async (ctx, args) => {
    return await ctx.db
      .query("attendance")
      .withIndex("by_section_date", (q) =>
        q.eq("sectionId", args.sectionId).eq("date", args.date)
      )
      .collect();
  },
});

export const listByStudent = query({
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

export const getSummary = query({
  args: {
    studentId: v.id("students"),
    academicYearId: v.id("academicYears"),
  },
  handler: async (ctx, args) => {
    const records = await ctx.db
      .query("attendance")
      .withIndex("by_student", (q) => q.eq("studentId", args.studentId))
      .filter((q) => q.eq(q.field("academicYearId"), args.academicYearId))
      .collect();

    const total = records.length;
    const present = records.filter(r => r.status === "present").length;
    const absent = records.filter(r => r.status === "absent").length;
    const late = records.filter(r => r.status === "late").length;
    const excused = records.filter(r => r.status === "excused").length;
    const percentage = total > 0 ? Math.round(((present + late) / total) * 100) : 0;

    return { present, absent, late, excused, total, percentage };
  },
});

// =============================================
// MUTATIONS
// =============================================

export const mark = mutation({
  args: {
    records: v.array(v.object({
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
      deviceId: v.optional(v.string()),
    })),
  },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    const identity = await requireAuth(ctx);
    const results = [];

    for (const record of args.records) {
      // Upsert: check for existing record, update or insert
      const existing = await ctx.db
        .query("attendance")
        .withIndex("by_student_date", (q) =>
          q.eq("studentId", record.studentId).eq("date", record.date)
        )
        .first();

      if (existing) {
        await ctx.db.patch(existing._id, {
          status: record.status,
          timeIn: record.timeIn,
          timeOut: record.timeOut,
          remarks: record.remarks,
          markedBy: identity._id,
        });
        results.push(existing._id);
      } else {
        const id = await ctx.db.insert("attendance", {
          tenantId,
          studentId: record.studentId,
          sectionId: record.sectionId,
          academicYearId: record.academicYearId,
          date: record.date,
          status: record.status,
          timeIn: record.timeIn,
          timeOut: record.timeOut,
          remarks: record.remarks,
          markedBy: identity._id,
          deviceId: record.deviceId,
        });
        results.push(id);
      }
    }

    return results;
  },
});

export const markBulk = mutation({
  args: {
    sectionId: v.id("sections"),
    academicYearId: v.id("academicYears"),
    date: v.string(),
    attendance: v.array(v.object({
      studentId: v.id("students"),
      status: v.union(
        v.literal("present"),
        v.literal("absent"),
        v.literal("late"),
        v.literal("excused")
      ),
    })),
  },
  handler: async (ctx, args) => {
    const tenantId = await getCurrentTenant(ctx);
    const identity = await requireAuth(ctx);

    for (const entry of args.attendance) {
      await ctx.db.insert("attendance", {
        tenantId,
        studentId: entry.studentId,
        sectionId: args.sectionId,
        academicYearId: args.academicYearId,
        date: args.date,
        status: entry.status,
        markedBy: identity._id,
      });
    }

    return { success: true, count: args.attendance.length };
  },
});