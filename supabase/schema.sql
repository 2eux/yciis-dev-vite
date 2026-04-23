-- =============================================
-- EDUSYS PRO - Supabase Database Schema
-- =============================================
-- Run this in Supabase Dashboard > SQL Editor
-- Or use: supabase db push
-- =============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================
-- AUTH & TENANTS
-- =============================================

-- Custom public profile view combining auth.users with tenant info
CREATE OR REPLACE VIEW auth.user_profiles AS
SELECT 
  au.id,
  au.email,
  au.created_at,
  au.email_confirmed_at,
  au.last_sign_in_at,
  au.encrypted_password,
  au.recovery_token,
  au.confirmation_token,
  au.email_change_token_new,
  au.reauthenticate_token,
  au.is_super_admin,
  au.email_change_token_current,
  au.email_change_token_new,
  au.email_change_token_current,
  au.phone_change_token,
  t.id AS tenant_id,
  t.name AS tenant_name,
  t.code AS tenant_code
FROM auth.users au
LEFT JOIN public.tenants t ON t.id = au.raw_user_meta_data->>'tenant_id';

-- =============================================
-- TENANTS (Schools/Organizations)
-- =============================================

CREATE TABLE public.tenants (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  timezone TEXT DEFAULT 'Asia/Jakarta',
  currency_code CHAR(3) DEFAULT 'IDR',
  logo_url TEXT,
  address TEXT,
  phone TEXT,
  email TEXT,
  is_active BOOLEAN DEFAULT true,
  subscription_plan TEXT DEFAULT 'basic',
  subscription_expires_at TIMESTAMPTZ,
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE public.tenants ENABLE ROW LEVEL SECURITY;

-- RLS Policies for tenants
CREATE POLICY "Tenants are viewable by authenticated users" 
  ON public.tenants FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role', 'anon'));

CREATE POLICY "Tenants are insertable by service role" 
  ON public.tenants FOR INSERT 
  WITH CHECK (auth.role() = 'service_role');

CREATE POLICY "Tenants are updatable by service role" 
  ON public.tenants FOR UPDATE 
  USING (auth.role() = 'service_role');

-- =============================================
-- USERS (Custom profile data)
-- =============================================

CREATE TABLE public.user_profiles (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID UNIQUE NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE SET NULL,
  role TEXT NOT NULL CHECK (role IN ('super_admin', 'admin', 'teacher', 'student', 'parent', 'finance', 'hr')),
  first_name TEXT NOT NULL,
  last_name TEXT,
  phone TEXT,
  avatar_url TEXT,
  is_active BOOLEAN DEFAULT true,
  two_secret_enabled BOOLEAN DEFAULT false,
  two_secret_secret TEXT,
  last_login_at TIMESTAMPTZ,
  last_login_ip TEXT,
  failed_login_attempts INT DEFAULT 0,
  locked_until TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.user_profiles ENABLE ROW LEVEL SECURITY;

CREATE POLICY "User profiles are viewable by authenticated users" 
  ON public.user_profiles FOR SELECT 
  USING (auth.uid() = user_id OR auth.role() IN ('service_role', 'admin'));

CREATE POLICY "Users can update own profile" 
  ON public.user_profiles FOR UPDATE 
  USING (auth.uid() = user_id);

CREATE POLICY "Service role can insert profiles" 
  ON public.user_profiles FOR INSERT 
  WITH CHECK (auth.role() = 'service_role');

-- =============================================
-- STUDENT MANAGEMENT
-- =============================================

CREATE TABLE public.students (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  student_id TEXT UNIQUE NOT NULL,
  gender TEXT CHECK (gender IN ('male', 'female', 'other')),
  date_of_birth DATE,
  place_of_birth TEXT,
  nationality TEXT DEFAULT 'Indonesia',
  religion TEXT,
  blood_type TEXT,
  address TEXT,
  city TEXT,
  province TEXT,
  postal_code TEXT,
  emergency_contact_name TEXT,
  emergency_contact_phone TEXT,
  emergency_contact_relation TEXT,
  notes TEXT,
  documents JSONB DEFAULT '[]',
  health_info JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.students ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Students viewable by authenticated" 
  ON public.students FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role') 
    OR user_id = auth.uid() 
    OR EXISTS (SELECT 1 FROM auth.users WHERE id = auth.uid() AND raw_user_meta_data->>'role' IN ('admin', 'super_admin', 'teacher')));

CREATE POLICY "Admins can manage students" 
  ON public.students FOR ALL 
  USING (auth.role() = 'service_role');

-- =============================================
-- STUDENT FAMILY (Parents/Guardians)
-- =============================================

CREATE TABLE public.student_parents (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID REFERENCES public.students(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  relation TEXT NOT NULL CHECK (relation IN ('father', 'mother', 'guardian', 'other')),
  first_name TEXT NOT NULL,
  last_name TEXT,
  email TEXT,
  phone TEXT,
  occupation TEXT,
  company TEXT,
  income_bracket TEXT,
  is_primary BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.student_parents ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Student parents viewable by authenticated" 
  ON public.student_parents FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role'));

-- =============================================
-- ACADEMIC YEARS
-- =============================================

CREATE TABLE public.academic_years (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  is_current BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(tenant_id, name)
);

ALTER TABLE public.academic_years ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Academic years viewable by authenticated" 
  ON public.academic_years FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role'));

-- =============================================
-- SECTIONS (Classes)
-- =============================================

CREATE TABLE public.sections (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  grade_level INT CHECK (grade_level BETWEEN 0 AND 12),
  room TEXT,
  capacity INT DEFAULT 35,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(tenant_id, academic_year_id, name)
);

ALTER TABLE public.sections ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Sections viewable by authenticated" 
  ON public.sections FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role'));

-- =============================================
-- SECTION STUDENTS (Enrollment)
-- =============================================

CREATE TABLE public.section_students (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID REFERENCES public.students(id) ON DELETE CASCADE,
  section_id UUID REFERENCES public.sections(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  roll_number TEXT,
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'promoted', 'transferred', 'dropped')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(student_id, academic_year_id)
);

ALTER TABLE public.section_students ENABLE ROW LEVEL SECURITY;

-- =============================================
-- SUBJECTS
-- =============================================

CREATE TABLE public.subjects (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  description TEXT,
  is_optional BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(tenant_id, code)
);

ALTER TABLE public.subjects ENABLE ROW LEVEL SECURITY;

-- =============================================
-- SUBJECT ASSIGNMENTS
-- =============================================

CREATE TABLE public.section_subjects (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  section_id UUID REFERENCES public.sections(id) ON DELETE CASCADE,
  subject_id UUID REFERENCES public.subjects(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  teacher_id UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(section_id, subject_id, academic_year_id)
);

ALTER TABLE public.section_subjects ENABLE ROW LEVEL SECURITY;

-- =============================================
-- TIMETABLES
-- =============================================

CREATE TABLE public.timetables (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  section_id UUID REFERENCES public.sections(id) ON DELETE CASCADE,
  section_subject_id UUID REFERENCES public.section_subjects(id) ON DELETE CASCADE,
  teacher_id UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  day_of_week INT CHECK (day_of_week BETWEEN 1 AND 7),
  period_start INT NOT NULL,
  period_end INT NOT NULL,
  room TEXT,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.timetables ENABLE ROW LEVEL SECURITY;

-- =============================================
-- ATTENDANCE
-- =============================================

CREATE TABLE public.student_attendance (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID NOT NULL,
  section_id UUID REFERENCES public.sections(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('present', 'absent', 'late', 'excused')),
  time_in TIME,
  time_out TIME,
  remarks TEXT,
  marked_by UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  device_id TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(student_id, date)
);

ALTER TABLE public.student_attendance ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Attendance viewable by authenticated" 
  ON public.student_attendance FOR SELECT 
  USING (auth.role() IN ('authenticated', 'service_role'));

-- =============================================
-- EXAMS
-- =============================================

CREATE TABLE public.exams (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  type TEXT CHECK (type IN ('quiz', 'midterm', 'final', 'semester', 'annual')),
  section_id UUID REFERENCES public.sections(id) ON DELETE CASCADE,
  subject_id UUID REFERENCES public.subjects(id) ON DELETE CASCADE,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  duration_minutes INT,
  total_marks DECIMAL(10,2),
  passing_marks DECIMAL(10,2),
  instructions TEXT,
  is_published BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.exams ENABLE ROW LEVEL SECURITY;

-- =============================================
-- EXAM QUESTIONS
-- =============================================

CREATE TABLE public.exam_questions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  exam_id UUID REFERENCES public.exams(id) ON DELETE CASCADE,
  question_text TEXT NOT NULL,
  question_type TEXT CHECK (question_type IN ('mcq', 'true_false', 'short', 'long')),
  option_a TEXT,
  option_b TEXT,
  option_c TEXT,
  option_d TEXT,
  correct_answer TEXT NOT NULL,
  marks DECIMAL(10,2) NOT NULL,
  sort_order INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.exam_questions ENABLE ROW LEVEL SECURITY;

-- =============================================
-- STUDENT MARKS
-- =============================================

CREATE TABLE public.student_exam_marks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID NOT NULL,
  exam_id UUID REFERENCES public.exams(id) ON DELETE CASCADE,
  question_id UUID REFERENCES public.exam_questions(id) ON DELETE CASCADE,
  answer_text TEXT,
  marks_obtained DECIMAL(10,2),
  graded_by UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  graded_at TIMESTAMPTZ,
  remarks TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(student_id, exam_id, question_id)
);

ALTER TABLE public.student_exam_marks ENABLE ROW LEVEL SECURITY;

-- =============================================
-- FEE STRUCTURES
-- =============================================

CREATE TABLE public.fee_structures (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  amount DECIMAL(15,2) NOT NULL,
  due_date DATE,
  is_recurring BOOLEAN DEFAULT false,
  frequency TEXT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.fee_structures ENABLE ROW LEVEL SECURITY;

-- =============================================
-- STUDENT FEES
-- =============================================

CREATE TABLE public.student_fees (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID REFERENCES public.students(id) ON DELETE CASCADE,
  fee_structure_id UUID REFERENCES public.fee_structures(id) ON DELETE CASCADE,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE CASCADE,
  amount DECIMAL(15,2) NOT NULL,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  final_amount DECIMAL(15,2) NOT NULL,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'partial', 'paid', 'overdue', 'waived')),
  due_date DATE,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  paid_at TIMESTAMPTZ,
  payment_method TEXT,
  transaction_id TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.student_fees ENABLE ROW LEVEL SECURITY;

-- =============================================
-- PAYMENT TRANSACTIONS
-- =============================================

CREATE TABLE public.payment_transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_fee_id UUID REFERENCES public.student_fees(id) ON DELETE SET NULL,
  amount DECIMAL(15,2) NOT NULL,
  payment_method TEXT NOT NULL,
  gateway TEXT,
  gateway_transaction_id TEXT,
  gateway_response JSONB,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'refunded')),
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.payment_transactions ENABLE ROW LEVEL SECURITY;

-- =============================================
-- HR - STAFF
-- =============================================

CREATE TABLE public.staff (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  employee_id TEXT UNIQUE NOT NULL,
  department TEXT,
  position TEXT,
  join_date DATE,
  employment_type TEXT DEFAULT 'full_time' CHECK (employment_type IN ('full_time', 'part_time', 'contract', 'intern')),
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'on_leave', 'terminated')),
  salary DECIMAL(15,2),
  bank_account TEXT,
  bank_name TEXT,
  emergency_contact_name TEXT,
  emergency_contact_phone TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.staff ENABLE ROW LEVEL SECURITY;

-- =============================================
-- HR - LEAVE REQUESTS
-- =============================================

CREATE TABLE public.leave_requests (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  staff_id UUID REFERENCES public.staff(id) ON DELETE CASCADE,
  leave_type TEXT NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  total_days INT NOT NULL,
  reason TEXT,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
  approved_by UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  approved_at TIMESTAMPTZ,
  remarks TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.leave_requests ENABLE ROW LEVEL SECURITY;

-- =============================================
-- AUDIT LOGS
-- =============================================

CREATE TABLE public.audit_logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE SET NULL,
  user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  module TEXT NOT NULL,
  entity_type TEXT,
  entity_id UUID,
  old_values JSONB,
  new_values JSONB,
  ip_address TEXT,
  user_agent TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.audit_logs ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Audit logs viewable by service role only" 
  ON public.audit_logs FOR SELECT 
  USING (auth.role() = 'service_role');

-- =============================================
-- LMS COURSES
-- =============================================

CREATE TABLE public.lms_courses (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  teacher_id UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  academic_year_id UUID REFERENCES public.academic_years(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  description TEXT,
  thumbnail_url TEXT,
  is_published BOOLEAN DEFAULT false,
  is_free BOOLEAN DEFAULT false,
  price DECIMAL(10,2) DEFAULT 0,
  language TEXT DEFAULT 'en',
  difficulty TEXT DEFAULT 'beginner',
  duration_hours DECIMAL(10,2),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.lms_courses ENABLE ROW LEVEL SECURITY;

-- =============================================
-- LMS SECTIONS
-- =============================================

CREATE TABLE public.lms_sections (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  course_id UUID REFERENCES public.lms_courses(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT,
  sort_order INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.lms_sections ENABLE ROW LEVEL SECURITY;

-- =============================================
-- LMS CONTENT
-- =============================================

CREATE TABLE public.lms_content (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  section_id UUID REFERENCES public.lms_sections(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  content_type TEXT NOT NULL CHECK (content_type IN ('video', 'document', 'quiz', 'assignment', 'link')),
  content_url TEXT,
  content_text TEXT,
  duration_minutes INT,
  sort_order INT,
  is_preview BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.lms_content ENABLE ROW LEVEL SECURITY;

-- =============================================
-- LMS ENROLLMENTS
-- =============================================

CREATE TABLE public.lms_enrollments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  course_id UUID REFERENCES public.lms_courses(id) ON DELETE CASCADE,
  student_id UUID NOT NULL,
  progress_percentage DECIMAL(5,2) DEFAULT 0,
  completed_at TIMESTAMPTZ,
  enrolled_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(course_id, student_id)
);

ALTER TABLE public.lms_enrollments ENABLE ROW LEVEL SECURITY;

-- =============================================
-- BOOKS
-- =============================================

CREATE TABLE public.books (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  isbn TEXT,
  title TEXT NOT NULL,
  author TEXT,
  publisher TEXT,
  publication_year INT,
  category TEXT,
  location TEXT,
  total_copies INT DEFAULT 1,
  available_copies INT DEFAULT 1,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.books ENABLE ROW LEVEL SECURITY;

-- =============================================
-- BOOK ISSUES
-- =============================================

CREATE TABLE public.book_issues (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  book_id UUID REFERENCES public.books(id) ON DELETE CASCADE,
  student_id UUID,
  staff_id UUID,
  issue_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  status TEXT DEFAULT 'issued' CHECK (status IN ('issued', 'returned', 'lost', 'overdue')),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.book_issues ENABLE ROW LEVEL SECURITY;

-- =============================================
-- ADMISSION LEADS
-- =============================================

CREATE TABLE public.admission_leads (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  first_name TEXT NOT NULL,
  last_name TEXT,
  email TEXT,
  phone TEXT NOT NULL,
  gender TEXT,
  date_of_birth DATE,
  grade_applied INT,
  source TEXT,
  status TEXT DEFAULT 'lead' CHECK (status IN ('lead', 'prospect', 'registered', 'rejected', 'waitlisted')),
  assigned_to UUID REFERENCES public.user_profiles(id) ON DELETE SET NULL,
  notes TEXT,
  documents JSONB DEFAULT '[]',
  follow_up_at TIMESTAMPTZ,
  converted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.admission_leads ENABLE ROW LEVEL SECURITY;

-- =============================================
-- MESSAGES
-- =============================================

CREATE TABLE public.messages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  receiver_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  subject TEXT,
  body TEXT NOT NULL,
  message_type TEXT DEFAULT 'general' CHECK (message_type IN ('general', 'announcement', 'notice')),
  priority TEXT DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
  is_read BOOLEAN DEFAULT false,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.messages ENABLE ROW LEVEL SECURITY;

-- =============================================
-- NOTIFICATIONS
-- =============================================

CREATE TABLE public.notifications (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  notification_type TEXT NOT NULL,
  data JSONB DEFAULT '{}',
  is_read BOOLEAN DEFAULT false,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.notifications ENABLE ROW LEVEL SECURITY;

-- =============================================
-- TRANSPORT ROUTES
-- =============================================

CREATE TABLE public.transport_routes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  start_point TEXT NOT NULL,
  end_point TEXT NOT NULL,
  waypoints JSONB DEFAULT '[]',
  distance_km DECIMAL(10,2),
  estimated_time_minutes INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.transport_routes ENABLE ROW LEVEL SECURITY;

-- =============================================
-- VEHICLES
-- =============================================

CREATE TABLE public.vehicles (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  vehicle_number TEXT NOT NULL,
  vehicle_type TEXT NOT NULL,
  model TEXT,
  capacity INT,
  driver_name TEXT,
  driver_phone TEXT,
  insurance_expiry DATE,
  fitness_expiry DATE,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE public.vehicles ENABLE ROW LEVEL SECURITY;

-- =============================================
-- BOARDING LOGS
-- =============================================

CREATE TABLE public.boarding_logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id UUID REFERENCES public.tenants(id) ON DELETE CASCADE,
  student_id UUID NOT NULL,
  route_assignment_id UUID REFERENCES public.transport_routes(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  pickup_time TIMESTAMP,
  pickup_location TEXT,
  dropoff_time TIMESTAMP,
  dropoff_location TEXT,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'boarded', 'dropped', 'noshow')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(student_id, route_assignment_id, date)
);

ALTER TABLE public.boarding_logs ENABLE ROW LEVEL SECURITY;

-- =============================================
-- TRIGGER FUNCTIONS
-- =============================================

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers to tables with updated_at
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_students_updated_at BEFORE UPDATE ON students
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sections_updated_at BEFORE UPDATE ON sections
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fee_structures_updated_at BEFORE UPDATE ON fee_structures
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_fees_updated_at BEFORE UPDATE ON student_fees
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lms_courses_updated_at BEFORE UPDATE ON lms_courses
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- FUNCTION: Create new user with profile
-- =============================================

CREATE OR REPLACE FUNCTION public.create_user_with_profile(
  p_email TEXT,
  p_password TEXT,
  p_role TEXT,
  p_first_name TEXT,
  p_last_name TEXT,
  p_tenant_id UUID DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
  v_user_id UUID;
BEGIN
  -- Create user in auth schema
  INSERT INTO auth.users (email, encrypted_password, email_confirmed_at, raw_user_meta_data)
  VALUES (p_email, crypt(p_password, gen_salt('bf')), NOW(), jsonb_build_object('tenant_id', p_tenant_id, 'role', p_role))
  RETURNING id INTO v_user_id;

  -- Create profile
  INSERT INTO public.user_profiles (user_id, tenant_id, role, first_name, last_name)
  VALUES (v_user_id, p_tenant_id, p_role, p_first_name, p_last_name);

  RETURN v_user_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- =============================================
-- VIEW: Dashboard Stats
-- =============================================

CREATE OR REPLACE VIEW public.dashboard_stats AS
SELECT 
  t.id AS tenant_id,
  t.name AS tenant_name,
  (SELECT COUNT(*) FROM students s WHERE s.tenant_id = t.id) AS total_students,
  (SELECT COUNT(*) FROM staff st WHERE st.tenant_id = t.id AND st.status = 'active') AS total_staff,
  (SELECT COUNT(*) FROM sections sec WHERE sec.tenant_id = t.id) AS total_sections,
  (SELECT COUNT(*) FROM user_profiles up WHERE up.tenant_id = t.id AND up.role = 'teacher') AS total_teachers
FROM tenants t
WHERE t.is_active = true;

-- =============================================
-- FUNCTION: Get current academic year
-- =============================================

CREATE OR REPLACE FUNCTION public.get_current_academic_year(p_tenant_id UUID)
RETURNS TABLE(id UUID, name TEXT, is_current BOOLEAN) AS $$
BEGIN
  RETURN QUERY
  SELECT ay.id, ay.name, ay.is_current
  FROM academic_years ay
  WHERE ay.tenant_id = p_tenant_id AND ay.is_current = true
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- FUNCTION: Student attendance summary
-- =============================================

CREATE OR REPLACE FUNCTION public.get_student_attendance_summary(
  p_student_id UUID,
  p_academic_year_id UUID
)
RETURNS TABLE(present INT, absent INT, late INT, total_days INT, percentage DECIMAL) AS $$
BEGIN
  RETURN QUERY
  SELECT 
    COUNT(*) FILTER (WHERE status = 'present') AS present,
    COUNT(*) FILTER (WHERE status = 'absent') AS absent,
    COUNT(*) FILTER (WHERE status = 'late') AS late,
    COUNT(*) AS total_days,
    COALESCE(
      ROUND(
        (COUNT(*) FILTER (WHERE status IN ('present', 'late'))::DECIMAL / NULLIF(COUNT(*), 0) * 100
      ), 0
    ) AS percentage
  FROM student_attendance sa
  WHERE sa.student_id = p_student_id 
    AND sa.academic_year_id = p_academic_year_id;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- SEED DATA: Sample Tenant
-- =============================================

INSERT INTO public.tenants (id, name, code, timezone, currency_code, email, is_active)
VALUES 
  (uuid_generate_v4(), 'Edusys Pro Academy', 'EDU001', 'Asia/Jakarta', 'IDR', 'info@edusys.edu', true)
ON CONFLICT (code) DO NOTHING;

-- =============================================
-- END OF SCHEMA
-- =============================================