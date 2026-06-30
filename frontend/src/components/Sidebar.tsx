import { NavLink } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: '📊', roles: ['super_admin', 'admin', 'member'] },
  { to: '/admins', label: 'Admins', icon: '🛡️', roles: ['super_admin'] },
  { to: '/stores', label: 'Stores', icon: '🏪', roles: ['super_admin'] },
  { to: '/members', label: 'Members', icon: '👥', roles: ['super_admin', 'admin'] },
  { to: '/invoices', label: 'Invoices', icon: '📄', roles: ['super_admin', 'admin'] },
  { to: '/my-invoices', label: 'My Invoices', icon: '📄', roles: ['member'] },
  { to: '/members/me/card', label: 'My Card', icon: '💳', roles: ['member'] },
]

export default function Sidebar() {
  const { user } = useAuthStore()

  return (
    <aside className="w-64 bg-white border-r border-gray-200 min-h-screen p-4 hidden lg:block">
      <div className="mb-8">
        <h1 className="text-xl font-bold text-brand-700">Members</h1>
        <p className="text-sm text-gray-500">Membership Platform</p>
      </div>

      <nav className="space-y-1">
        {navItems
          .filter((item) => user && item.roles.includes(user.role))
          .map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-brand-50 text-brand-700'
                    : 'text-gray-600 hover:bg-gray-50'
                }`
              }
            >
              <span>{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
      </nav>

      <div className="mt-auto pt-8 border-t border-gray-200 mt-8">
        <p className="text-sm font-medium text-gray-900">{user?.name}</p>
        <p className="text-xs text-gray-500 capitalize">{user?.role}</p>
        <button
          onClick={() => { localStorage.removeItem('members-auth'); window.location.href = '/login' }}
          className="mt-3 text-xs text-red-500 hover:text-red-700"
        >
          Logout
        </button>
      </div>
    </aside>
  )
}
