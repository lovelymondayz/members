import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'
import { useAuthStore } from '../store/authStore'

export default function DashboardPage() {
  const { user } = useAuthStore()

  const { data: stats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      try {
        const res = await api.get('/admin/dashboard')
        return res.data as { total_stores: number; total_members: number; total_revenue: number }
      } catch {
        return null
      }
    },
  })

  const statCards = user?.role === 'super_admin'
    ? [
        { title: 'Total Stores', value: String(stats?.total_stores ?? 0), icon: '🏪' },
        { title: 'Total Members', value: String(stats?.total_members ?? 0), icon: '👥' },
        { title: 'Total Revenue', value: `Rp ${(stats?.total_revenue ?? 0).toLocaleString('id-ID')}`, icon: '💰' },
      ]
    : [
        { title: 'Total Members', value: String(stats?.total_members ?? 0), icon: '👥' },
        { title: 'Total Invoices', value: '0', icon: '📄' },
        { title: 'Revenue', value: `Rp ${(stats?.total_revenue ?? 0).toLocaleString('id-ID')}`, icon: '💰' },
      ]

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p className="text-gray-500 mt-1">Welcome back, {user?.name}</p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        {statCards.map((card) => (
          <StatCard key={card.title} title={card.title} value={card.value} icon={card.icon} />
        ))}
      </div>
    </div>
  )
}

function StatCard({ title, value, icon }: { title: string; value: string; icon: string }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-6">
      <div className="flex items-center justify-between">
        <span className="text-2xl">{icon}</span>
        <span className="text-2xl font-bold text-gray-900">{value}</span>
      </div>
      <p className="text-sm text-gray-500 mt-2">{title}</p>
    </div>
  )
}
