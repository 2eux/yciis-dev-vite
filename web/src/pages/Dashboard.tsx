import { useQuery } from '@tanstack/react-query'
import { 
  Users, GraduationCap, Briefcase, CalendarClock, DollarSign, 
  TrendingUp, TrendingDown, AlertTriangle, MoreHorizontal,
  Calendar, Award, FileText, Wallet
} from 'lucide-react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell } from 'recharts'

const mockData = {
  overview: {
    total_students: 1250,
    total_teachers: 85,
    total_staff: 42,
    total_sections: 35,
  },
  attendance: {
    present: 1180,
    absent: 45,
    late: 25,
    percentage: 94.4,
  },
  fees: {
    collected: 12500000,
    pending: 2500000,
    overdue: 500000,
  },
}

const weeklyAttendance = [
  { day: 'Mon', present: 1180, absent: 45 },
  { day: 'Tue', present: 1200, absent: 30 },
  { day: 'Wed', present: 1195, absent: 35 },
  { day: 'Thu', present: 1210, absent: 25 },
  { day: 'Fri', present: 1185, absent: 40 },
]

const feePieData = [
  { name: 'Collected', value: 12500000, color: '#10b981' },
  { name: 'Pending', value: 2500000, color: '#f59e0b' },
  { name: 'Overdue', value: 500000, color: '#ef4444' },
]

const recentActivities = [
  { id: 1, title: 'New admission', description: 'Aryan Sharma admitted to Class 9-A', time: '5 min ago', color: 'bg-primary' },
  { id: 2, title: 'Fee payment', description: '₹45,000 from Rahul Verma', time: '15 min ago', color: 'bg-emerald-500' },
  { id: 3, title: 'Attendance', description: 'Class 8-A: 98% present', time: '30 min ago', color: 'bg-amber-500' },
  { id: 4, title: 'Results published', description: 'Science Midterm Grade 10', time: '1 hour ago', color: 'bg-violet-500' },
  { id: 5, title: 'New staff', description: 'Ms. Priya Singh as Math teacher', time: '2 hours ago', color: 'bg-primary' },
]

const upcomingEvents = [
  { id: 1, title: 'Parent-Teacher Meeting', date: '25 Mar', type: 'event' },
  { id: 2, title: 'Unit Test - Science', date: '28 Mar', type: 'exam' },
  { id: 3, title: 'Annual Day Practice', date: '01 Apr', type: 'event' },
  { id: 4, title: 'Fee Due Date', date: '05 Apr', type: 'fee' },
]

const alerts = [
  { id: 1, type: 'warning', message: '15 students have pending library books' },
  { id: 2, type: 'danger', message: '5 fee payments overdue > 30 days' },
  { id: 3, type: 'info', message: 'Mid-term exam marks pending' },
]

const subjectPerformance = [
  { name: 'Mathematics', avg: 85, color: '#4f46e5' },
  { name: 'Science', avg: 78, color: '#10b981' },
  { name: 'English', avg: 82, color: '#f59e0b' },
  { name: 'Social', avg: 75, color: '#ef4444' },
  { name: 'Hindi', avg: 88, color: '#8b5cf6' },
]

interface StatCardProps {
  title: string
  value: string | number
  change?: string
  isPositive?: boolean
  icon: any
  color: string
}

function StatCard({ title, value, change, isPositive, icon: Icon, color }: StatCardProps) {
  return (
    <div className="card">
      <div className="flex items-start justify-between">
        <div className="flex h-12 w-12 items-center justify-center rounded-xl" style={{ backgroundColor: `${color}15` }}>
          <Icon className="h-6 w-6" style={{ color }} />
        </div>
        {change && (
          <span className={`flex items-center text-xs font-medium ${isPositive ? 'text-emerald-600' : 'text-red-600'}`}>
            {isPositive ? <TrendingUp className="mr-1 h-3 w-3" /> : <TrendingDown className="mr-1 h-3 w-3" />}
            {change}
          </span>
        )}
      </div>
      <div className="mt-4">
        <h3 className="stat-value">{value}</h3>
        <p className="stat-label">{title}</p>
      </div>
    </div>
  )
}

export default function Dashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: async () => {
      await new Promise(resolve => setTimeout(resolve, 300))
      return mockData
    },
  })

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(amount)
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="card p-5">
              <div className="skeleton h-10 w-10 rounded-xl mb-4"></div>
              <div className="skeleton h-8 w-24 mb-2"></div>
              <div className="skeleton h-4 w-16"></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="page-header">
        <div>
          <h1 className="page-title">Dashboard</h1>
          <p className="page-description">Welcome back! Here's what's happening at your school today.</p>
        </div>
        <div className="flex items-center gap-3">
          <select className="flex h-10 w-auto items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm">
            <option>Academic Year 2023-24</option>
            <option>Academic Year 2022-23</option>
          </select>
          <button className="btn btn-primary">
            <Calendar className="mr-2 h-4 w-4" />
            View Calendar
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard title="Total Students" value={data?.overview.total_students.toLocaleString()} change="+12" isPositive={true} icon={Users} color="#4f46e5" />
        <StatCard title="Total Teachers" value={data?.overview.total_teachers} change="+3" isPositive={true} icon={GraduationCap} color="#10b981" />
        <StatCard title="Total Staff" value={data?.overview.total_staff} icon={Briefcase} color="#f59e0b" />
        <StatCard title="Today's Attendance" value={`${data?.attendance.percentage}%`} change="+2.4%" isPositive={true} icon={CalendarClock} color="#8b5cf6" />
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Attendance Chart */}
        <div className="card lg:col-span-2">
          <div className="card-header">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="card-title text-lg">Weekly Attendance</h3>
                <p className="text-sm text-muted-foreground">Student attendance overview</p>
              </div>
              <button className="btn btn-ghost btn-sm">
                <MoreHorizontal className="h-4 w-4" />
              </button>
            </div>
          </div>
          <div className="card-content">
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={weeklyAttendance}>
                  <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" vertical={false} />
                  <XAxis dataKey="day" axisLine={false} tickLine={false} tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }} />
                  <YAxis axisLine={false} tickLine={false} tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }} />
                  <Tooltip 
                    contentStyle={{ borderRadius: '8px', border: '1px solid hsl(var(--border))', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }}
                    cursor={{ fill: 'hsl(var(--accent))' }}
                  />
                  <Bar dataKey="present" name="Present" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
                  <Bar dataKey="absent" name="Absent" fill="#fca5a5" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        {/* Fee Distribution */}
        <div className="card">
          <div className="card-header">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="card-title text-lg">Fee Collection</h3>
                <p className="text-sm text-muted-foreground">Payment status overview</p>
              </div>
              <button className="btn btn-ghost btn-sm">
                <MoreHorizontal className="h-4 w-4" />
              </button>
            </div>
          </div>
          <div className="card-content">
            <div className="h-48">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={feePieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={50}
                    outerRadius={70}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {feePieData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => formatCurrency(value as number)} />
                </PieChart>
              </ResponsiveContainer>
            </div>
            <div className="space-y-3 mt-4">
              {feePieData.map((item, i) => (
                <div key={i} className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className="h-3 w-3 rounded-full" style={{ backgroundColor: item.color }}></div>
                    <span className="text-sm text-muted-foreground">{item.name}</span>
                  </div>
                  <span className="font-semibold">{formatCurrency(item.value)}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Bottom Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Activities */}
        <div className="card lg:col-span-2">
          <div className="card-header">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="card-title text-lg">Recent Activities</h3>
                <p className="text-sm text-muted-foreground">Latest updates from your school</p>
              </div>
              <button className="text-sm text-primary font-medium hover:text-primary/80">View All</button>
            </div>
          </div>
          <div className="card-content p-0">
            <div className="divide-y">
              {recentActivities.map((activity) => (
                <div key={activity.id} className="flex items-center gap-4 px-6 py-4 hover:bg-accent transition-colors cursor-pointer">
                  <div className={`flex h-10 w-10 items-center justify-center rounded-xl ${activity.color}`}>
                    <Users className="h-5 w-5 text-white" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium">{activity.title}</p>
                    <p className="text-sm text-muted-foreground truncate">{activity.description}</p>
                  </div>
                  <span className="text-xs text-muted-foreground">{activity.time}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Right Column */}
        <div className="space-y-6">
          {/* Upcoming Events */}
          <div className="card">
            <div className="card-header">
              <div className="flex items-center justify-between">
                <h3 className="card-title text-lg">Upcoming Events</h3>
                <Calendar className="h-5 w-5 text-muted-foreground" />
              </div>
            </div>
            <div className="card-content space-y-3">
              {upcomingEvents.map((event) => (
                <div key={event.id} className="flex items-center justify-between rounded-lg border p-3 hover:bg-accent transition-colors cursor-pointer">
                  <div>
                    <p className="font-medium text-sm">{event.title}</p>
                    <p className="text-xs text-muted-foreground">{event.date}</p>
                  </div>
                  <span className={`badge ${event.type === 'exam' ? 'badge-secondary' : event.type === 'fee' ? 'badge-destructive' : 'badge-primary'}`}>
                    {event.type}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Alerts */}
          <div className="card">
            <div className="card-header">
              <div className="flex items-center justify-between">
                <h3 className="card-title text-lg">Alerts</h3>
                <span className="badge badge-destructive">{alerts.length}</span>
              </div>
            </div>
            <div className="card-content space-y-3">
              {alerts.map((alert) => (
                <div 
                  key={alert.id} 
                  className={`flex items-center gap-3 rounded-lg p-3 ${
                    alert.type === 'danger' ? 'bg-destructive/10' : 
                    alert.type === 'warning' ? 'bg-amber-100 dark:bg-amber-950' : 'bg-blue-100 dark:bg-blue-950'
                  }`}
                >
                  <AlertTriangle className={`h-5 w-5 ${
                    alert.type === 'danger' ? 'text-destructive' : 
                    alert.type === 'warning' ? 'text-amber-600' : 'text-blue-600'
                  }`} />
                  <p className="text-sm">{alert.message}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Subject Performance */}
      <div className="card">
        <div className="card-header">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="card-title text-lg">Subject Performance</h3>
              <p className="text-sm text-muted-foreground">Average marks by subject</p>
            </div>
            <button className="btn btn-secondary btn-sm">
              <Award className="mr-2 h-4 w-4" />
              View Report Cards
            </button>
          </div>
        </div>
        <div className="card-content">
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            {subjectPerformance.map((subject, i) => (
              <div key={i} className="text-center rounded-lg border p-4">
                <div className="relative mx-auto mb-3 h-16 w-16">
                  <svg className="h-16 w-16 transform -rotate-90">
                    <circle cx="32" cy="32" r="28" stroke="hsl(var(--border))" strokeWidth="4" fill="none" />
                    <circle 
                      cx="32" cy="32" r="28" 
                      stroke={subject.color} 
                      strokeWidth="4" 
                      fill="none"
                      strokeDasharray={`${subject.avg * 1.76} 176`}
                      strokeLinecap="round"
                    />
                  </svg>
                  <span className="absolute inset-0 flex items-center justify-center text-sm font-bold">
                    {subject.avg}%
                  </span>
                </div>
                <p className="font-medium">{subject.name}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}