import { Outlet, Navigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import Sidebar from '../components/Sidebar'

export default function MainLayout() {
  const { isAuthenticated } = useAuthStore()

  if (!isAuthenticated) return <Navigate to="/login" replace />

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 p-6 lg:p-8">
        <Outlet />
      </main>
    </div>
  )
}
