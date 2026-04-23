export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json }
  | Json[]

export interface Database {
  public: {
    Tables: {
      tenants: {
        Row: {
          id: string
          name: string
          code: string
          timezone: string
          currency_code: string
          logo_url: string | null
          address: string | null
          phone: string | null
          email: string | null
          is_active: boolean
          subscription_plan: string
          subscription_expires_at: string | null
          settings: Json
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          name: string
          code: string
          timezone?: string
          currency_code?: string
          logo_url?: string | null
          address?: string | null
          phone?: string | null
          email?: string | null
          is_active?: boolean
          subscription_plan?: string
          subscription_expires_at?: string | null
          settings?: Json
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          name?: string
          code?: string
          timezone?: string
          currency_code?: string
          logo_url?: string | null
          address?: string | null
          phone?: string | null
          email?: string | null
          is_active?: boolean
          subscription_plan?: string
          subscription_expires_at?: string | null
          settings?: Json
          created_at?: string
          updated_at?: string
        }
      }
      user_profiles: {
        Row: {
          id: string
          user_id: string
          tenant_id: string | null
          role: string
          first_name: string
          last_name: string | null
          phone: string | null
          avatar_url: string | null
          is_active: boolean
          two_secret_enabled: boolean
          two_secret_secret: string | null
          last_login_at: string | null
          last_login_ip: string | null
          failed_login_attempts: number
          locked_until: string | null
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          user_id: string
          tenant_id?: string | null
          role: string
          first_name: string
          last_name?: string | null
          phone?: string | null
          avatar_url?: string | null
          is_active?: boolean
          two_secret_enabled?: boolean
          two_secret_secret?: string | null
          last_login_at?: string | null
          last_login_ip?: string | null
          failed_login_attempts?: number
          locked_until?: string | null
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          user_id?: string
          tenant_id?: string | null
          role?: string
          first_name?: string
          last_name?: string | null
          phone?: string | null
          avatar_url?: string | null
          is_active?: boolean
          two_secret_enabled?: boolean
          two_secret_secret?: string | null
          last_login_at?: string | null
          last_login_ip?: string | null
          failed_login_attempts?: number
          locked_until?: string | null
          created_at?: string
          updated_at?: string
        }
      }
      students: {
        Row: {
          id: string
          tenant_id: string | null
          user_id: string | null
          student_id: string
          gender: string | null
          date_of_birth: string | null
          place_of_birth: string | null
          nationality: string
          religion: string | null
          blood_type: string | null
          address: string | null
          city: string | null
          province: string | null
          postal_code: string | null
          emergency_contact_name: string | null
          emergency_contact_phone: string | null
          emergency_contact_relation: string | null
          notes: string | null
          documents: Json
          health_info: Json
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          user_id?: string | null
          student_id: string
          gender?: string | null
          date_of_birth?: string | null
          place_of_birth?: string | null
          nationality?: string
          religion?: string | null
          blood_type?: string | null
          address?: string | null
          city?: string | null
          province?: string | null
          postal_code?: string | null
          emergency_contact_name?: string | null
          emergency_contact_phone?: string | null
          emergency_contact_relation?: string | null
          notes?: string | null
          documents?: Json
          health_info?: Json
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          user_id?: string | null
          student_id?: string
          gender?: string | null
          date_of_birth?: string | null
          place_of_birth?: string | null
          nationality?: string
          religion?: string | null
          blood_type?: string | null
          address?: string | null
          city?: string | null
          province?: string | null
          postal_code?: string | null
          emergency_contact_name?: string | null
          emergency_contact_phone?: string | null
          emergency_contact_relation?: string | null
          notes?: string | null
          documents?: Json
          health_info?: Json
          created_at?: string
          updated_at?: string
        }
      }
      academic_years: {
        Row: {
          id: string
          tenant_id: string | null
          name: string
          start_date: string
          end_date: string
          is_current: boolean
          is_active: boolean
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          name: string
          start_date: string
          end_date: string
          is_current?: boolean
          is_active?: boolean
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          name?: string
          start_date?: string
          end_date?: string
          is_current?: boolean
          is_active?: boolean
          created_at?: string
          updated_at?: string
        }
      }
      sections: {
        Row: {
          id: string
          tenant_id: string | null
          academic_year_id: string | null
          name: string
          grade_level: number | null
          room: string | null
          capacity: number
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          academic_year_id?: string | null
          name: string
          grade_level?: number | null
          room?: string | null
          capacity?: number
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          academic_year_id?: string | null
          name?: string
          grade_level?: number | null
          room?: string | null
          capacity?: number
          created_at?: string
          updated_at?: string
        }
      }
      student_attendance: {
        Row: {
          id: string
          tenant_id: string | null
          student_id: string
          section_id: string | null
          academic_year_id: string | null
          date: string
          status: string
          time_in: string | null
          time_out: string | null
          remarks: string | null
          marked_by: string | null
          device_id: string | null
          created_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          student_id: string
          section_id?: string | null
          academic_year_id?: string | null
          date: string
          status: string
          time_in?: string | null
          time_out?: string | null
          remarks?: string | null
          marked_by?: string | null
          device_id?: string | null
          created_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          student_id?: string
          section_id?: string | null
          academic_year_id?: string | null
          date?: string
          status?: string
          time_in?: string | null
          time_out?: string | null
          remarks?: string | null
          marked_by?: string | null
          device_id?: string | null
          created_at?: string
        }
      }
      exams: {
        Row: {
          id: string
          tenant_id: string | null
          academic_year_id: string | null
          name: string
          type: string
          section_id: string | null
          subject_id: string | null
          start_date: string
          end_date: string
          duration_minutes: number | null
          total_marks: number | null
          passing_marks: number | null
          instructions: string | null
          is_published: boolean
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          academic_year_id?: string | null
          name: string
          type: string
          section_id?: string | null
          subject_id?: string | null
          start_date: string
          end_date: string
          duration_minutes?: number | null
          total_marks?: number | null
          passing_marks?: number | null
          instructions?: string | null
          is_published?: boolean
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          academic_year_id?: string | null
          name?: string
          type?: string
          section_id?: string | null
          subject_id?: string | null
          start_date?: string
          end_date?: string
          duration_minutes?: number | null
          total_marks?: number | null
          passing_marks?: number | null
          instructions?: string | null
          is_published?: boolean
          created_at?: string
          updated_at?: string
        }
      }
      student_fees: {
        Row: {
          id: string
          tenant_id: string | null
          student_id: string
          fee_structure_id: string | null
          academic_year_id: string | null
          amount: number
          discount_amount: number
          final_amount: number
          status: string
          due_date: string | null
          paid_amount: number
          paid_at: string | null
          payment_method: string | null
          transaction_id: string | null
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          student_id: string
          fee_structure_id?: string | null
          academic_year_id?: string | null
          amount: number
          discount_amount?: number
          final_amount: number
          status?: string
          due_date?: string | null
          paid_amount?: number
          paid_at?: string | null
          payment_method?: string | null
          transaction_id?: string | null
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          student_id?: string
          fee_structure_id?: string | null
          academic_year_id?: string | null
          amount?: number
          discount_amount?: number
          final_amount?: number
          status?: string
          due_date?: string | null
          paid_amount?: number
          paid_at?: string | null
          payment_method?: string | null
          transaction_id?: string | null
          created_at?: string
          updated_at?: string
        }
      }
      staff: {
        Row: {
          id: string
          tenant_id: string | null
          user_id: string | null
          employee_id: string
          department: string | null
          position: string | null
          join_date: string | null
          employment_type: string
          status: string
          salary: number | null
          bank_account: string | null
          bank_name: string | null
          emergency_contact_name: string | null
          emergency_contact_phone: string | null
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          user_id?: string | null
          employee_id: string
          department?: string | null
          position?: string | null
          join_date?: string | null
          employment_type?: string
          status?: string
          salary?: number | null
          bank_account?: string | null
          bank_name?: string | null
          emergency_contact_name?: string | null
          emergency_contact_phone?: string | null
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          user_id?: string | null
          employee_id?: string
          department?: string | null
          position?: string | null
          join_date?: string | null
          employment_type?: string
          status?: string
          salary?: number | null
          bank_account?: string | null
          bank_name?: string | null
          emergency_contact_name?: string | null
          emergency_contact_phone?: string | null
          created_at?: string
          updated_at?: string
        }
      }
      lms_courses: {
        Row: {
          id: string
          tenant_id: string | null
          teacher_id: string | null
          academic_year_id: string | null
          title: string
          description: string | null
          thumbnail_url: string | null
          is_published: boolean
          is_free: boolean
          price: number
          language: string
          difficulty: string
          duration_hours: number | null
          created_at: string
          updated_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          teacher_id?: string | null
          academic_year_id?: string | null
          title: string
          description?: string | null
          thumbnail_url?: string | null
          is_published?: boolean
          is_free?: boolean
          price?: number
          language?: string
          difficulty?: string
          duration_hours?: number | null
          created_at?: string
          updated_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          teacher_id?: string | null
          academic_year_id?: string | null
          title?: string
          description?: string | null
          thumbnail_url?: string | null
          is_published?: boolean
          is_free?: boolean
          price?: number
          language?: string
          difficulty?: string
          duration_hours?: number | null
          created_at?: string
          updated_at?: string
        }
      }
      notifications: {
        Row: {
          id: string
          tenant_id: string | null
          user_id: string
          title: string
          body: string
          notification_type: string
          data: Json
          is_read: boolean
          read_at: string | null
          created_at: string
        }
        Insert: {
          id?: string
          tenant_id?: string | null
          user_id: string
          title: string
          body: string
          notification_type: string
          data?: Json
          is_read?: boolean
          read_at?: string | null
          created_at?: string
        }
        Update: {
          id?: string
          tenant_id?: string | null
          user_id?: string
          title?: string
          body?: string
          notification_type?: string
          data?: Json
          is_read?: boolean
          read_at?: string | null
          created_at?: string
        }
      }
    }
    Views: {
      dashboard_stats: {
        Row: {
          tenant_id: string
          tenant_name: string
          total_students: number
          total_staff: number
          total_sections: number
          total_teachers: number
        }
      }
    }
    Functions: {
      create_user_with_profile: {
        Args: {
          p_email: string
          p_password: string
          p_role: string
          p_first_name: string
          p_last_name: string
          p_tenant_id: string | null
        }
        Returns: string
      }
      get_current_academic_year: {
        Args: {
          p_tenant_id: string
        }
        Returns: {
          id: string
          name: string
          is_current: boolean
        }[]
      }
      get_student_attendance_summary: {
        Args: {
          p_student_id: string
          p_academic_year_id: string
        }
        Returns: {
          present: number
          absent: number
          late: number
          total_days: number
          percentage: number
        }[]
      }
    }
  }
}

export type Tenant = Database['public']['Tables']['tenants']['Row']
export type UserProfile = Database['public']['Tables']['user_profiles']['Row']
export type Student = Database['public']['Tables']['students']['Row']
export type AcademicYear = Database['public']['Tables']['academic_years']['Row']
export type Section = Database['public']['Tables']['sections']['Row']
export type StudentAttendance = Database['public']['Tables']['student_attendance']['Row']
export type Exam = Database['public']['Tables']['exams']['Row']
export type StudentFee = Database['public']['Tables']['student_fees']['Row']
export type Staff = Database['public']['Tables']['staff']['Row']
export type LMSCourse = Database['public']['Tables']['lms_courses']['Row']
export type Notification = Database['public']['Tables']['notifications']['Row']
export type DashboardStats = Database['public']['Views']['dashboard_stats']['Row']