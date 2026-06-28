import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

interface Member {
  member_id: string
  user_id?: string
  user_name?: string
  member_code: string
  tier: string
  store_id: string
  joined_at?: string
  store_name?: string
  store_card_color?: string
}

export default function MembersPage() {
  const [showCreate, setShowCreate] = useState(false)
  const [code, setCode] = useState('')
  const [tier, setTier] = useState('standard')
  const [submitting, setSubmitting] = useState(false)
  const qc = useQueryClient()
  const { data: members, isLoading } = useQuery({
    queryKey: ['members'],
    queryFn: async () => {
      const res = await api.get('/members')
      return res.data as Member[]
    },
  })

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      await api.post('/members', { member_code: code, tier })
      setCode('')
      setTier('standard')
      setShowCreate(false)
      qc.invalidateQueries({ queryKey: ['members'] })
    } catch {
      alert('Failed to create member')
    } finally {
      setSubmitting(false)
    }
  }

  if (isLoading) return <div className="text-gray-500">Loading members...</div>

  return (
    <div>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Members</h1>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="bg-brand-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-brand-700"
        >
          + Add Member
        </button>
      </div>

      {showCreate && (
        <div className="mt-6 bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">New Member</h2>
          <form onSubmit={handleCreate} className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Member Code</label>
              <input
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="M001"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Tier</label>
              <select
                value={tier}
                onChange={(e) => setTier(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
              >
                <option value="standard">Standard</option>
                <option value="premium">Premium</option>
                <option value="vip">VIP</option>
              </select>
            </div>
            <div className="flex items-end gap-2">
              <button
                type="submit"
                disabled={submitting}
                className="bg-brand-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-brand-700 disabled:opacity-50"
              >
                {submitting ? 'Creating...' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => setShowCreate(false)}
                className="border border-gray-300 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-50"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="mt-6 bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Code</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Name</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Tier</th>
              <th className="hidden md:table-cell text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Joined</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {members?.map((member) => (
              <tr key={member.member_id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-mono text-gray-600">{member.member_code}</td>
                <td className="px-6 py-4 text-sm font-medium text-gray-900">
                 {member.user_name || <span className="text-gray-400">—</span>}
                </td>
                <td className="px-6 py-4">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-brand-50 text-brand-700 capitalize">
                    {member.tier}
                  </span>
                </td>
                <td className="hidden md:table-cell px-6 py-4 text-sm text-gray-500">
                  {member.joined_at ? new Date(member.joined_at).toLocaleDateString('id-ID') : '—'}
                </td>
                <td className="px-6 py-4">
                  <Link
                    to={`/members/${member.member_id}/card`}
                    className="text-brand-600 hover:text-brand-700 text-sm font-medium"
                  >
                    View Card
                  </Link>
                </td>
              </tr>
            ))}
            {(!members || members.length === 0) && (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-gray-500">
                  No members yet. Add your first member to get started.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
