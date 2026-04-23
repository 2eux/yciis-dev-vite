import { supabase } from './supabase'
import type { Database, Tenant, Student, Section, AcademicYear, Exam, StudentFee, Staff, LMSCourse, Notification, UserProfile, DashboardStats } from '../types/database'

// =============================================
// AUTH QUERIES
// =============================================

export const auth = {
  // Sign up with email/password
  signUp: async (email: string, password: string, metadata: Record<string, unknown> = {}) => {
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
      options: { data: metadata },
    })
    return { data, error }
  },

  // Sign in with email/password
  signIn: async (email: string, password: string) => {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    return { data, error }
  },

  // Sign out
  signOut: async () => {
    const { error } = await supabase.auth.signOut()
    return { error }
  },

  // Get current session
  getSession: async () => {
    const { data: { session }, error } = await supabase.auth.getSession()
    return { session, error }
  },

  // Get current user
  getUser: async () => {
    const { data: { user }, error } = await supabase.auth.getUser()
    return { user, error }
  },

  // Reset password
  resetPassword: async (email: string) => {
    const { data, error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: `${window.location.origin}/reset-password`,
    })
    return { data, error }
  },

  // Update password
  updatePassword: async (newPassword: string) => {
    const { data, error } = await supabase.auth.updateUser({ password: newPassword })
    return { data, error }
  },

  // Listen to auth changes
  onAuthStateChange: (callback: (event: string, session: unknown) => void) => {
    return supabase.auth.onAuthStateChange(callback)
  },
}

// =============================================
// USER PROFILE QUERIES
// =============================================

export const userProfile = {
  // Get current user profile
  getCurrent: async () => {
    const { data: { user } } = await supabase.auth.getUser()
    if (!user) return { data: null, error: null }

    const { data, error } = await supabase
      .from('user_profiles')
      .select('*')
      .eq('user_id', user.id)
      .single()

    return { data: data as UserProfile, error }
  },

  // Get user profile by ID
  getById: async (id: string) => {
    const { data, error } = await supabase
      .from('user_profiles')
      .select('*')
      .eq('id', id)
      .single()

    return { data: data as UserProfile, error }
  },

  // Get user profile by user_id
  getByUserId: async (userId: string) => {
    const { data, error } = await supabase
      .from('user_profiles')
      .select('*')
      .eq('user_id', userId)
      .single()

    return { data: data as UserProfile, error }
  },

  // Update profile
  update: async (id: string, updates: Partial<UserProfile>) => {
    const { data, error } = await supabase
      .from('user_profiles')
      .update(updates)
      .eq('id', id)
      .select()
      .single()

    return { data: data as UserProfile, error }
  },
}

// =============================================
// STUDENT QUERIES
// =============================================

export const students = {
  // List students with filters
  list: async (tenantId: string, options?: {
    search?: string
    sectionId?: string
    status?: string
    limit?: number
    page?: number
  }) => {
    let query = supabase
      .from('students')
      .select('*', { count: 'exact' })
      .eq('tenant_id', tenantId)

    if (options?.search) {
      query = query.or(`first_name.ilike.%${options.search}%,last_name.ilike.%${options.search}%,student_id.ilike.%${options.search}%`)
    }
    if (options?.sectionId) {
      query = query.eq('section_id', options.sectionId)
    }
    if (options?.status) {
      query = query.eq('status', options.status)
    }
    if (options?.limit) {
      query = query.limit(options.limit)
    }

    const { data, error, count } = await query

    return { 
      data: data as Student[], 
      error, 
      count,
    }
  },

  // Get student by ID
  getById: async (id: string) => {
    const { data, error } = await supabase
      .from('students')
      .select('*')
      .eq('id', id)
      .single()

    return { data: data as Student, error }
  },

  // Create student
  create: async (student: Partial<Student>) => {
    const { data, error } = await supabase
      .from('students')
      .insert(student)
      .select()
      .single()

    return { data: data as Student, error }
  },

  // Update student
  update: async (id: string, updates: Partial<Student>) => {
    const { data, error } = await supabase
      .from('students')
      .update(updates)
      .eq('id', id)
      .select()
      .single()

    return { data: data as Student, error }
  },

  // Delete student
  delete: async (id: string) => {
    const { error } = await supabase
      .from('students')
      .delete()
      .eq('id', id)

    return { error }
  },
}

// =============================================
// ACADEMIC YEAR QUERIES
// =============================================

export const academicYears = {
  // List academic years
  list: async (tenantId: string) => {
    const { data, error } = await supabase
      .from('academic_years')
      .select('*')
      .eq('tenant_id', tenantId)
      .order('start_date', { ascending: false })

    return { data: data as AcademicYear[], error }
  },

  // Get current academic year
  getCurrent: async (tenantId: string) => {
    const { data, error } = await supabase
      .from('academic_years')
      .select('*')
      .eq('tenant_id', tenantId)
      .eq('is_current', true)
      .single()

    return { data: data as AcademicYear, error }
  },

  // Create academic year
  create: async (year: Partial<AcademicYear>) => {
    const { data, error } = await supabase
      .from('academic_years')
      .insert(year)
      .select()
      .single()

    return { data: data as AcademicYear, error }
  },

  // Update academic year
  update: async (id: string, updates: Partial<AcademicYear>) => {
    const { data, error } = await supabase
      .from('academic_years')
      .update(updates)
      .eq('id', id)
      .select()
      .single()

    return { data: data as AcademicYear, error }
  },
}

// =============================================
// SECTION QUERIES
// =============================================

export const sections = {
  // List sections
  list: async (tenantId: string, academicYearId?: string) => {
    let query = supabase
      .from('sections')
      .select('*')
      .eq('tenant_id', tenantId)
      .order('grade_level')
      .order('name')

    if (academicYearId) {
      query = query.eq('academic_year_id', academicYearId)
    }

    const { data, error } = await query

    return { data: data as Section[], error }
  },

  // Get section with students
  getWithStudents: async (sectionId: string) => {
    const { data: section, error: sectionError } = await supabase
      .from('sections')
      .select('*')
      .eq('id', sectionId)
      .single()

    if (sectionError) return { section: null, students: [], error: sectionError }

    const { data: students, error: studentsError } = await supabase
      .from('section_students')
      .select(`
        student_id,
        roll_number,
        status,
        students (
          id,
          student_id,
          first_name,
          last_name,
          gender,
          date_of_birth
        )
      `)
      .eq('section_id', sectionId)
      .eq('status', 'active')

    return { section, students, error: studentsError }
  },

  // Create section
  create: async (section: Partial<Section>) => {
    const { data, error } = await supabase
      .from('sections')
      .insert(section)
      .select()
      .single()

    return { data: data as Section, error }
  },

  // Update section
  update: async (id: string, updates: Partial<Section>) => {
    const { data, error } = await supabase
      .from('sections')
      .update(updates)
      .eq('id', id)
      .select()
      .single()

    return { data: data as Section, error }
  },
}

// =============================================
// ATTENDANCE QUERIES
// =============================================

export const attendance = {
  // Get student attendance
  getStudentAttendance: async (studentId: string, academicYearId: string) => {
    const { data, error } = await supabase
      .from('student_attendance')
      .select('*')
      .eq('student_id', studentId)
      .eq('academic_year_id', academicYearId)
      .order('date', { ascending: false })

    return { data, error }
  },

  // Get section attendance for a date
  getSectionAttendance: async (sectionId: string, date: string) => {
    const { data, error } = await supabase
      .from('student_attendance')
      .select(`
        *,
        student:students (
          id,
          student_id,
          first_name,
          last_name
        )
      `)
      .eq('section_id', sectionId)
      .eq('date', date)

    return { data, error }
  },

  // Mark attendance
  mark: async (records: Array<{
    student_id: string
    section_id: string
    academic_year_id: string
    date: string
    status: string
    time_in?: string
    remarks?: string
  }>) => {
    const { data, error } = await supabase
      .from('student_attendance')
      .upsert(records, { onConflict: 'student_id,date' })
      .select()

    return { data, error }
  },
}

// =============================================
// EXAM QUERIES
// =============================================

export const exams = {
  // List exams
  list: async (tenantId: string, options?: {
    sectionId?: string
    subjectId?: string
    academicYearId?: string
  }) => {
    let query = supabase
      .from('exams')
      .select('*')
      .eq('tenant_id', tenantId)

    if (options?.sectionId) query = query.eq('section_id', options.sectionId)
    if (options?.subjectId) query = query.eq('subject_id', options.subjectId)
    if (options?.academicYearId) query = query.eq('academic_year_id', options.academicYearId)

    const { data, error } = await query.order('start_date', { ascending: false })

    return { data: data as Exam[], error }
  },

  // Get exam with questions
  getWithQuestions: async (examId: string) => {
    const { data: exam, error: examError } = await supabase
      .from('exams')
      .select('*')
      .eq('id', examId)
      .single()

    if (examError) return { exam: null, questions: [], error: examError }

    const { data: questions, error: questionsError } = await supabase
      .from('exam_questions')
      .select('*')
      .eq('exam_id', examId)
      .order('sort_order')

    return { exam, questions, error: questionsError }
  },

  // Create exam
  create: async (exam: Partial<Exam>) => {
    const { data, error } = await supabase
      .from('exams')
      .insert(exam)
      .select()
      .single()

    return { data: data as Exam, error }
  },

  // Submit marks
  submitMarks: async (marks: Array<{
    student_id: string
    exam_id: string
    question_id: string
    marks_obtained: number
    graded_by: string
  }>) => {
    const { data, error } = await supabase
      .from('student_exam_marks')
      .upsert(marks)
      .select()

    return { data, error }
  },
}

// =============================================
// FEE QUERIES
// =============================================

export const fees = {
  // Get student fees
  getStudentFees: async (studentId: string) => {
    const { data, error } = await supabase
      .from('student_fees')
      .select(`
        *,
        fee_structure:_fee_structures (
          name,
          amount,
          due_date
        )
      `)
      .eq('student_id', studentId)
      .order('due_date')

    return { data: data as StudentFee[], error }
  },

  // Get fee summary
  getFeeSummary: async (tenantId: string, academicYearId: string) => {
    const { data, error } = await supabase
      .from('student_fees')
      .select('status, final_amount, paid_amount')
      .eq('tenant_id', tenantId)
      .eq('academic_year_id', academicYearId)

    if (error) return { collected: 0, pending: 0, overdue: 0, error }

    const summary = {
      collected: data?.filter(f => f.status === 'paid').reduce((sum, f) => sum + (f.paid_amount || 0), 0) || 0,
      pending: data?.filter(f => f.status === 'pending').reduce((sum, f) => sum + (f.final_amount - f.paid_amount), 0) || 0,
      overdue: data?.filter(f => f.status === 'overdue').reduce((sum, f) => sum + (f.final_amount - f.paid_amount), 0) || 0,
    }

    return { ...summary, error }
  },

  // Process payment
  processPayment: async (feeId: string, amount: number, method: string, gateway: string) => {
    // Update fee status
    const { data: fee, error: feeError } = await supabase
      .from('student_fees')
      .update({
        status: 'paid',
        paid_amount: amount,
        paid_at: new Date().toISOString(),
        payment_method: method,
      })
      .eq('id', feeId)
      .select()
      .single()

    if (feeError) return { error: feeError }

    // Record transaction
    const { data: transaction, error: txError } = await supabase
      .from('payment_transactions')
      .insert({
        student_fee_id: feeId,
        amount,
        payment_method: method,
        gateway,
        status: 'success',
        paid_at: new Date().toISOString(),
      })
      .select()
      .single()

    return { fee, transaction, error: txError }
  },
}

// =============================================
// STAFF QUERIES
// =============================================

export const staff = {
  // List staff
  list: async (tenantId: string, options?: {
    department?: string
    status?: string
  }) => {
    let query = supabase
      .from('staff')
      .select('*')
      .eq('tenant_id', tenantId)

    if (options?.department) query = query.eq('department', options.department)
    if (options?.status) query = query.eq('status', options.status)

    const { data, error } = await query.order('name')

    return { data: data as Staff[], error }
  },

  // Get staff by user_id
  getByUserId: async (userId: string) => {
    const { data, error } = await supabase
      .from('staff')
      .select('*')
      .eq('user_id', userId)
      .single()

    return { data: data as Staff, error }
  },

  // Create staff
  create: async (member: Partial<Staff>) => {
    const { data, error } = await supabase
      .from('staff')
      .insert(member)
      .select()
      .single()

    return { data: data as Staff, error }
  },

  // Update staff
  update: async (id: string, updates: Partial<Staff>) => {
    const { data, error } = await supabase
      .from('staff')
      .update(updates)
      .eq('id', id)
      .select()
      .single()

    return { data: data as Staff, error }
  },
}

// =============================================
// LMS QUERIES
// =============================================

export const lms = {
  // List courses
  listCourses: async (tenantId: string, options?: {
    teacherId?: string
    published?: boolean
  }) => {
    let query = supabase
      .from('lms_courses')
      .select('*')
      .eq('tenant_id', tenantId)

    if (options?.teacherId) query = query.eq('teacher_id', options.teacherId)
    if (options?.published !== undefined) query = query.eq('is_published', options.published)

    const { data, error } = await query.order('created_at', { ascending: false })

    return { data: data as LMSCourse[], error }
  },

  // Get course with content
  getCourseWithContent: async (courseId: string) => {
    const { data: course, error: courseError } = await supabase
      .from('lms_courses')
      .select('*')
      .eq('id', courseId)
      .single()

    if (courseError) return { course: null, sections: [], error: courseError }

    const { data: sections, error: sectionsError } = await supabase
      .from('lms_sections')
      .select('*, content:lms_content(*)')
      .eq('course_id', courseId)
      .order('sort_order')

    return { course, sections, error: sectionsError }
  },

  // Enroll student
  enroll: async (courseId: string, studentId: string, tenantId: string) => {
    const { data, error } = await supabase
      .from('lms_enrollments')
      .insert({
        tenant_id: tenantId,
        course_id: courseId,
        student_id: studentId,
      })
      .select()
      .single()

    return { data, error }
  },

  // Update progress
  updateProgress: async (enrollmentId: string, progress: number) => {
    const { data, error } = await supabase
      .from('lms_enrollments')
      .update({
        progress_percentage: progress,
        completed_at: progress >= 100 ? new Date().toISOString() : null,
      })
      .eq('id', enrollmentId)
      .select()
      .single()

    return { data, error }
  },
}

// =============================================
// NOTIFICATIONS QUERIES
// =============================================

export const notifications = {
  // Get user notifications
  list: async (userId: string, options?: { unreadOnly?: boolean; limit?: number }) => {
    let query = supabase
      .from('notifications')
      .select('*')
      .eq('user_id', userId)

    if (options?.unreadOnly) query = query.eq('is_read', false)
    if (options?.limit) query = query.limit(options.limit)

    const { data, error } = await query.order('created_at', { ascending: false })

    return { data: data as Notification[], error }
  },

  // Mark as read
  markAsRead: async (id: string) => {
    const { error } = await supabase
      .from('notifications')
      .update({ is_read: true, read_at: new Date().toISOString() })
      .eq('id', id)

    return { error }
  },

  // Mark all as read
  markAllAsRead: async (userId: string) => {
    const { error } = await supabase
      .from('notifications')
      .update({ is_read: true, read_at: new Date().toISOString() })
      .eq('user_id', userId)
      .eq('is_read', false)

    return { error }
  },

  // Create notification
  create: async (notification: Partial<Notification>) => {
    const { data, error } = await supabase
      .from('notifications')
      .insert(notification)
      .select()
      .single()

    return { data: data as Notification, error }
  },
}

// =============================================
// DASHBOARD QUERIES
// =============================================

export const dashboard = {
  // Get dashboard stats
  getStats: async (tenantId: string) => {
    const { data, error } = await supabase
      .from('dashboard_stats')
      .select('*')
      .eq('tenant_id', tenantId)
      .single()

    return { data: data as DashboardStats, error }
  },

  // Get today's attendance
  getTodayAttendance: async (tenantId: string, date: string) => {
    const { data, error } = await supabase
      .from('student_attendance')
      .select('status')
      .eq('tenant_id', tenantId)
      .eq('date', date)

    if (error) return { present: 0, absent: 0, late: 0, percentage: 0, error }

    const present = data?.filter(a => a.status === 'present').length || 0
    const absent = data?.filter(a => a.status === 'absent').length || 0
    const late = data?.filter(a => a.status === 'late').length || 0
    const total = present + absent + late
    const percentage = total > 0 ? Math.round(((present + late) / total) * 100) : 0

    return { present, absent, late, percentage, error }
  },

  // Get recent activities
  getRecentActivities: async (tenantId: string, limit: number = 10) => {
    const { data, error } = await supabase
      .from('audit_logs')
      .select('*, user:user_profiles(first_name, last_name)')
      .eq('tenant_id', tenantId)
      .order('created_at', { ascending: false })
      .limit(limit)

    return { data, error }
  },
}

// =============================================
// REAL-TIME SUBSCRIPTIONS
// =============================================

export const realtime = {
  // Subscribe to notifications
  subscribeToNotifications: (userId: string, callback: (payload: unknown) => void) => {
    return supabase
      .channel('notifications')
      .on('postgres_changes', {
        event: 'INSERT',
        schema: 'public',
        table: 'notifications',
        filter: `user_id=eq.${userId}`,
      }, callback)
      .subscribe()
  },

  // Subscribe to attendance changes
  subscribeToAttendance: (sectionId: string, date: string, callback: (payload: unknown) => void) => {
    return supabase
      .channel('attendance')
      .on('postgres_changes', {
        event: '*',
        schema: 'public',
        table: 'student_attendance',
        filter: `section_id=eq.${sectionId}`,
      }, callback)
      .subscribe()
  },

  // Subscribe to messages
  subscribeToMessages: (userId: string, callback: (payload: unknown) => void) => {
    return supabase
      .channel('messages')
      .on('postgres_changes', {
        event: 'INSERT',
        schema: 'public',
        table: 'messages',
        filter: `receiver_id=eq.${userId}`,
      }, callback)
      .subscribe()
  },
}