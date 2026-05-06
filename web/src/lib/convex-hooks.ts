// =============================================
// EDUSYS PRO - Convex Hooks
// =============================================

import { useQuery, useMutation } from "convex/react";
import { api } from "../../../convex/_generated/api";

// =============================================
// STUDENT HOOKS
// =============================================

export function useStudents(options?: {
  search?: string;
  sectionId?: string;
  status?: string;
  limit?: number;
  page?: number;
}) {
  return useQuery(api.students.list, options ?? {});
}

export function useStudent(id: string | undefined) {
  return useQuery(api.students.getById, id ? { id: id as any } : "skip");
}

export function useCreateStudent() {
  return useMutation(api.students.create);
}

export function useUpdateStudent() {
  return useMutation(api.students.update);
}

export function useDeleteStudent() {
  return useMutation(api.students.remove);
}

// =============================================
// ATTENDANCE HOOKS
// =============================================

export function useSectionAttendance(sectionId?: string, date?: string) {
  return useQuery(
    api.attendance.listBySection,
    sectionId && date ? { sectionId: sectionId as any, date } : "skip"
  );
}

export function useStudentAttendance(studentId?: string, academicYearId?: string) {
  return useQuery(
    api.attendance.listByStudent,
    studentId && academicYearId ? { studentId: studentId as any, academicYearId: academicYearId as any } : "skip"
  );
}

export function useMarkAttendance() {
  return useMutation(api.attendance.mark);
}

export function useMarkBulkAttendance() {
  return useMutation(api.attendance.markBulk);
}

export function useAttendanceSummary(studentId?: string, academicYearId?: string) {
  return useQuery(
    api.attendance.getSummary,
    studentId && academicYearId ? { studentId: studentId as any, academicYearId: academicYearId as any } : "skip"
  );
}

// =============================================
// STUDENT FEES HOOKS
// =============================================

export function useStudentFees(studentId?: string) {
  return useQuery(
    api.students.getFees,
    studentId ? { studentId: studentId as any } : "skip"
  );
}

// =============================================
// DASHBOARD HOOKS
// =============================================

export function useDashboardStats() {
  return useQuery(api.helpers.getDashboardStats, {});
}

export function useTodayAttendance(date: string) {
  return useQuery(api.helpers.getTodayAttendance, { date });
}

export function useFeeSummary() {
  return useQuery(api.helpers.getFeeSummary, {});
}

export function useRecentActivities(limit?: number) {
  return useQuery(api.helpers.getRecentActivities, { limit });
}

// =============================================
// ACADEMIC HOOKS
// =============================================

export function useSections(academicYearId?: string) {
  // Simple list for now - can be expanded
  return useQuery(api.students.list, {});
}

// =============================================
// FAMILY HOOKS
// =============================================

export function useStudentFamily(studentId?: string) {
  return useQuery(
    api.students.getFamily,
    studentId ? { studentId: studentId as any } : "skip"
  );
}

export function useAddFamilyMember() {
  return useMutation(api.students.addFamilyMember);
}