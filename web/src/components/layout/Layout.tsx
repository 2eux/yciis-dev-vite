import { useState } from 'react'
import { useNavigate, useLocation, NavLink } from 'react-router-dom'
import { 
  GraduationCap, Users, BookOpen, CalendarClock, FileText, 
  Wallet, Briefcase, Truck, Library, Settings, LogOut, Bell, Search,
  Menu, Plus, ChevronDown, Home, BarChart3, X
} from 'lucide-react'
import { useConvexAuth } from '../stores/convex-auth'

const navItems = [
  { path: '/', icon: Home, label: 'Dashboard' },
  { path: '/students', icon: Users, label: 'Students' },
  { path: '/academic', icon: BookOpen, label: 'Academic' },
  { path: '/attendance', icon: CalendarClock, label: 'Attendance' },
  { path: '/exams', icon: FileText, label: 'Exams' },
  { path: '/fees', icon: Wallet, label: 'Finance' },
  { path: '/hr', icon: Briefcase, label: 'HR' },
  { path: '/lms', icon: BarChart3, label: 'LMS' },
  { path: '/library', icon: Library, label: 'Library' },
  { path: '/transport', icon: Truck, label: 'Transport' },
]

export default function Layout() {
  const { user, logout } = useConvexAuth()
  const navigate = useNavigate()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Sticky Header */}
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center px-4 md:px-6">
          {/* Logo */}
          <div className="flex items-center gap-4">
            <button 
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="lg:hidden p-2 rounded-md hover:bg-accent"
            >
              {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </button>
            <div className="flex items-center gap-3 cursor-pointer" onClick={() => navigate('/')}>
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary">
                <GraduationCap className="h-5 w-5 text-primary-foreground" />
              </div>
              <div className="hidden sm:block">
                <h1 className="text-lg font-bold leading-none">Edusys Pro</h1>
                <p className="text-xs text-muted-foreground">Enterprise ERP</p>
              </div>
            </div>
          </div>

          {/* Horizontal Navigation */}
          <nav className="hidden lg:flex items-center gap-1 mx-6">
            {navItems.map((item) => (
              <NavLink 
                key={item.path}
                to={item.path}
                className={({ isActive }) => 
                  `inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 h-10 px-4 py-2 ${
                    isActive 
                      ? 'bg-accent text-accent-foreground' 
                      : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                  }`
                }
              >
                {item.label}
              </NavLink>
            ))}
          </nav>

          {/* Right Section */}
          <div className="flex flex-1 items-center justify-end gap-2">
            {/* Search */}
            <div className="hidden md:flex items-center">
              <div className="relative w-64">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <input 
                  type="text" 
                  placeholder="Search..." 
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 pl-9 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>
            </div>

            {/* Quick Add */}
            <button className="hidden md:inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
              <Plus className="h-4 w-4 mr-2" />
              Quick Add
            </button>

            {/* Notifications */}
            <button className="relative inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-10 w-10">
              <Bell className="h-5 w-5" />
              <span className="absolute -top-0.5 -right-0.5 flex h-2.5 w-2.5">
                <span className="relative flex h-full w-full">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-red-500"></span>
                </span>
              </span>
            </button>

            {/* User Menu */}
            <div className="relative">
              <button 
                onClick={() => setUserMenuOpen(!userMenuOpen)}
                className="flex items-center gap-2 rounded-md hover:bg-accent p-2"
              >
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary">
                  <span className="text-sm font-medium text-primary-foreground">
                    {user?.first_name?.[0]}{user?.last_name?.[0]}
                  </span>
                </div>
                <ChevronDown className="h-4 w-4 text-muted-foreground hidden sm:block" />
              </button>
              
              {userMenuOpen && (
                <div className="absolute right-0 mt-2 w-56 rounded-md border bg-popover p-1 shadow-md animate-in z-50">
                  <div className="px-2 py-1.5 border-b">
                    <p className="text-sm font-medium">{user?.first_name} {user?.last_name}</p>
                    <p className="text-xs text-muted-foreground">{user?.email}</p>
                  </div>
                  <div className="p-1">
                    <button 
                      onClick={() => { navigate('/settings'); setUserMenuOpen(false) }}
                      className="relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
                    >
                      <Settings className="mr-2 h-4 w-4" />
                      Settings
                    </button>
                    <button 
                      onClick={handleLogout}
                      className="relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground text-red-500"
                    >
                      <LogOut className="mr-2 h-4 w-4" />
                      Logout
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <nav className="lg:hidden border-t px-4 py-4 bg-background">
            <div className="grid grid-cols-2 gap-2">
              {navItems.map((item) => (
                <NavLink 
                  key={item.path}
                  to={item.path}
                  onClick={() => setMobileMenuOpen(false)}
                  className={({ isActive }) => 
                    `flex items-center gap-3 rounded-md px-4 py-3 ${
                      isActive 
                        ? 'bg-accent text-accent-foreground' 
                        : 'text-foreground'
                    }`
                  }
                >
                  <item.icon className="h-5 w-5" />
                  <span className="font-medium">{item.label}</span>
                </NavLink>
              ))}
            </div>
            <div className="mt-3 pt-3 border-t">
              <button 
                onClick={() => { navigate('/settings'); setMobileMenuOpen(false) }}
                className="flex items-center gap-3 w-full rounded-md px-4 py-3 text-foreground"
              >
                <Settings className="h-5 w-5" />
                <span className="font-medium">Settings</span>
              </button>
              <button 
                onClick={handleLogout}
                className="flex items-center gap-3 w-full rounded-md px-4 py-3 text-red-500"
              >
                <LogOut className="h-5 w-5" />
                <span className="font-medium">Logout</span>
              </button>
            </div>
          </nav>
        )}
      </header>

      {/* Main Content */}
      <main className="p-4 md:p-6">
        <div className="mx-auto max-w-7xl">
          <Outlet />
        </div>
      </main>

      {userMenuOpen && (
        <div className="fixed inset-0 z-40" onClick={() => setUserMenuOpen(false)} />
      )}
    </div>
  )
}