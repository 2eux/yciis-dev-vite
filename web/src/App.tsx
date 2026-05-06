import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConvexProvider, ConvexReactClient } from "convex/react"
import Layout from './components/layout/Layout'
import Dashboard from './pages/Dashboard'
import Login from './pages/Login'
import Students from './pages/Students'
import Academic from './pages/Academic'
import Attendance from './pages/Attendance'
import Exams from './pages/Exams'
import Fees from './pages/Fees'
import HR from './pages/HR'
import LMS from './pages/LMS'
import Library from './pages/Library'
import Transport from './pages/Transport'
import Settings from './pages/Settings'
import { useConvexAuth } from './stores/convex-auth'

const convexUrl = import.meta.env.VITE_CONVEX_URL
if (!convexUrl) throw new Error('VITE_CONVEX_URL is required')

const convex = new ConvexReactClient(convexUrl)

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useConvexAuth()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

export default function App() {
  return (
    <ConvexProvider client={convex}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <PrivateRoute>
                <Layout />
              </PrivateRoute>
            }
          >
            <Route index element={<Dashboard />} />
            <Route path="students" element={<Students />} />
            <Route path="academic" element={<Academic />} />
            <Route path="attendance" element={<Attendance />} />
            <Route path="exams" element={<Exams />} />
            <Route path="fees" element={<Fees />} />
            <Route path="hr" element={<HR />} />
            <Route path="lms" element={<LMS />} />
            <Route path="library" element={<Library />} />
            <Route path="transport" element={<Transport />} />
            <Route path="settings" element={<Settings />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConvexProvider>
  )
}