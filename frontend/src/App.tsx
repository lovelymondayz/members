import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import MainLayout from './layouts/MainLayout'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import MembersPage from './pages/MembersPage'
import MemberCardPage from './pages/MemberCardPage'
import InvoicesPage from './pages/InvoicesPage'
import NotFoundPage from './pages/NotFoundPage'

function ProtectedRoute({ children, allowedRoles }: { children: React.ReactNode; allowedRoles?: string[] }) {
  const { user, isAuthenticated } = useAuthStore()

  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (allowedRoles && user && !allowedRoles.includes(user.role)) {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="members" element={
          <ProtectedRoute allowedRoles={['super_admin', 'admin']}>
            <MembersPage />
          </ProtectedRoute>
        } />
        <Route path="members/:id/card" element={<MemberCardPage />} />
        <Route path="invoices" element={
          <ProtectedRoute allowedRoles={['super_admin', 'admin']}>
            <InvoicesPage />
          </ProtectedRoute>
        } />
        <Route path="my-invoices" element={
          <ProtectedRoute allowedRoles={['member']}>
            <InvoicesPage />
          </ProtectedRoute>
        } />
      </Route>
      {/* 404 Catch-All */}
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}
