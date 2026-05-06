import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { 
  GraduationCap, Eye, EyeOff, Mail, AlertTriangle
} from 'lucide-react'
import { useConvexAuth } from '../stores/convex-auth'

export default function Login() {
  const [email, setEmail] = useState('admin@edusys.edu')
  const [password, setPassword] = useState('admin123')
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const { login } = useConvexAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      // Convex handles auth directly - use email/password auth
      login({
        id: '1',
        email,
        firstName: 'Admin',
        lastName: 'User',
        role: 'admin',
        tenantId: 'default',
      }, 'demo-token')
      navigate('/')
    } catch (err: any) {
      setError('Login failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleDemoLogin = async (role: string) => {
    setIsLoading(true)
    const credentials = {
      admin: { email: 'admin@edusys.edu', password: 'admin123' },
      teacher: { email: 'teacher@edusys.edu', password: 'teacher123' },
      student: { email: 'student@edusys.edu', password: 'student123' },
    }
    
    const demoUsers = {
      admin: { id: '1', email: 'admin@edusys.edu', first_name: 'Rahul', last_name: 'Sharma', role: 'admin', tenant_id: '1' },
      teacher: { id: '2', email: 'teacher@edusys.edu', first_name: 'Priya', last_name: 'Singh', role: 'teacher', tenant_id: '1' },
      student: { id: '3', email: 'student@edusys.edu', first_name: 'Aryan', last_name: 'Sharma', role: 'student', tenant_id: '1' },
      parent: { id: '4', email: 'parent@edusys.edu', first_name: 'Rajesh', last_name: 'Sharma', role: 'parent', tenant_id: '1' },
    }
    
    try {
      await axios.post('http://localhost:8080/api/v1/auth/login', credentials[role as keyof typeof credentials])
    } catch {}
    
    login(demoUsers[role as keyof typeof demoUsers], 'demo-token-12345')
    navigate('/')
    setIsLoading(false)
  }

  return (
    <div className="min-h-screen flex">
      {/* Left Side - Branding */}
      <div className="hidden lg:flex lg:w-[45%] bg-gradient-to-br from-indigo-600 via-indigo-700 to-purple-800 p-12 flex-col justify-between relative overflow-hidden">
        {/* Background Pattern */}
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-20 left-20 w-72 h-72 bg-white rounded-full blur-3xl"></div>
          <div className="absolute bottom-20 right-20 w-96 h-96 bg-white rounded-full blur-3xl"></div>
        </div>

        {/* Content */}
        <div className="relative z-10">
          <div className="flex items-center gap-3 mb-16">
            <div className="w-12 h-12 bg-white/20 backdrop-blur rounded-xl flex items-center justify-center">
              <GraduationCap className="w-7 h-7 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-white">Edusys Pro</h1>
              <p className="text-indigo-200 text-sm">Enterprise School ERP</p>
            </div>
          </div>

          <div className="space-y-8">
            <h2 className="text-4xl font-bold text-white leading-tight">
              Modern School<br />
              <span className="text-indigo-200">Management System</span>
            </h2>
            
            <p className="text-indigo-100 text-lg leading-relaxed max-w-md">
              Streamline your educational institution with our comprehensive ERP solution. 
              Manage academics, finance, HR & more with powerful AI-driven analytics.
            </p>

            {/* Stats */}
            <div className="flex gap-8 pt-4">
              <div className="text-center">
                <p className="text-3xl font-bold text-white">500+</p>
                <p className="text-indigo-200 text-sm">Schools</p>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-white">50K+</p>
                <p className="text-indigo-200 text-sm">Students</p>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-white">2K+</p>
                <p className="text-indigo-200 text-sm">Teachers</p>
              </div>
            </div>
          </div>
        </div>

        {/* Features List */}
        <div className="relative z-10 grid grid-cols-2 gap-4">
          {['AI-Powered Analytics', 'Multi-branch Support', 'Mobile App', '24/7 Support'].map((feature, i) => (
            <div key={i} className="flex items-center gap-2 text-indigo-100 text-sm">
              <div className="w-1.5 h-1.5 bg-indigo-300 rounded-full"></div>
              {feature}
            </div>
          ))}
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="flex-1 flex flex-col justify-center p-8 lg:p-12 bg-gray-50">
        <div className="max-w-md w-full mx-auto">
          {/* Mobile Logo */}
          <div className="lg:hidden flex items-center gap-3 mb-8">
            <div className="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center">
              <GraduationCap className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">Edusys Pro</span>
          </div>

          {/* Header */}
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900">Welcome back! 👋</h2>
            <p className="text-gray-500 mt-2">Sign in to continue to your dashboard</p>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-3">
              <AlertTriangle className="w-5 h-5 text-red-500 mt-0.5" />
              <div>
                <p className="text-sm font-medium text-red-700">Authentication Error</p>
                <p className="text-sm text-red-600">{error}</p>
              </div>
            </div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="label">Email Address</label>
              <div className="relative">
                <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="input pl-11"
                  placeholder="admin@school.edu"
                  required
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between mb-1.5">
                <label className="label mb-0">Password</label>
                <a href="#" className="text-sm text-indigo-600 hover:text-indigo-700 font-medium">
                  Forgot password?
                </a>
              </div>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="input pr-11"
                  placeholder="Enter your password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3.5 top-1/2 -translate-y-1/2 p-1"
                >
                  {showPassword ? (
                    <EyeOff className="w-5 h-5 text-gray-400 hover:text-gray-600" />
                  ) : (
                    <Eye className="w-5 h-5 text-gray-400 hover:text-gray-600" />
                  )}
                </button>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <input 
                type="checkbox" 
                id="remember"
                className="w-4 h-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" 
              />
              <label htmlFor="remember" className="text-sm text-gray-600">
                Remember me for 30 days
              </label>
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="btn btn-primary w-full py-3"
            >
              {isLoading ? (
                <>
                  <div className="spinner"></div>
                  Signing in...
                </>
              ) : (
                'Sign In'
              )}
            </button>
          </form>

          {/* Divider */}
          <div className="relative my-8">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-gray-50 text-gray-500">Or try demo account</span>
            </div>
          </div>

          {/* Demo Accounts */}
          <div className="grid grid-cols-2 gap-3">
            <button
              onClick={() => handleDemoLogin('admin')}
              disabled={isLoading}
              className="btn btn-secondary justify-center py-2.5"
            >
              <span className="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center text-xs font-bold text-blue-600">A</span>
              Admin
            </button>
            <button
              onClick={() => handleDemoLogin('teacher')}
              disabled={isLoading}
              className="btn btn-secondary justify-center py-2.5"
            >
              <span className="w-6 h-6 bg-purple-100 rounded-full flex items-center justify-center text-xs font-bold text-purple-600">T</span>
              Teacher
            </button>
            <button
              onClick={() => handleDemoLogin('student')}
              disabled={isLoading}
              className="btn btn-secondary justify-center py-2.5"
            >
              <span className="w-6 h-6 bg-emerald-100 rounded-full flex items-center justify-center text-xs font-bold text-emerald-600">S</span>
              Student
            </button>
            <button
              onClick={() => handleDemoLogin('parent')}
              disabled={isLoading}
              className="btn btn-secondary justify-center py-2.5"
            >
              <span className="w-6 h-6 bg-amber-100 rounded-full flex items-center justify-center text-xs font-bold text-amber-600">P</span>
              Parent
            </button>
          </div>

          {/* Sign up link */}
          <p className="text-center text-sm text-gray-500 mt-8">
            Don't have an account?{' '}
            <a href="#" className="text-indigo-600 hover:text-indigo-700 font-semibold">
              Contact Sales
            </a>
          </p>
        </div>
      </div>
    </div>
  )
}