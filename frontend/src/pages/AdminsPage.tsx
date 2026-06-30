import { useState, useEffect } from 'react'
import { getAdmins, createAdmin, updateAdmin, deleteAdmin, AdminUser } from '../api/admin'

export default function AdminsPage() {
  const [admins, setAdmins] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<AdminUser | null>(null)
  const [form, setForm] = useState({
    name: '', email: '', password: '',
    store_name: '', address: '', phone: '', card_color_hex: '#1E40AF',
  })

  useEffect(() => { fetchAdmins() }, [])

  const fetchAdmins = async () => {
    try {
      const data = await getAdmins()
      setAdmins(data)
    } catch { alert('Failed to load admins') }
    finally { setLoading(false) }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createAdmin(form)
      setShowForm(false)
      setForm({ name: '', email: '', password: '', store_name: '', address: '', phone: '', card_color_hex: '#1E40AF' })
      fetchAdmins()
    } catch { alert('Failed to create admin') }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editing) return
    try {
      await updateAdmin(editing.user_id, { name: form.name, email: form.email })
      setEditing(null)
      setForm({ name: '', email: '', password: '', store_name: '', address: '', phone: '', card_color_hex: '#1E40AF' })
      fetchAdmins()
    } catch { alert('Failed to update admin') }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this admin and their store?')) return
    try { await deleteAdmin(id); fetchAdmins() }
    catch { alert('Failed to delete admin') }
  }

  const openEdit = (a: AdminUser) => {
    setEditing(a)
    setForm({ name: a.name, email: a.email, password: '', store_name: a.store_name, address: a.store_address, phone: a.store_phone, card_color_hex: a.store_color })
    setShowForm(false)
  }

  if (loading) return <div className="p-6 text-gray-500">Loading...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">Store Admins</h1>
        <button
          onClick={() => { setEditing(null); setShowForm(!showForm); setForm({ name: '', email: '', password: '', store_name: '', address: '', phone: '', card_color_hex: '#1E40AF' }) }}
          className="px-4 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 font-medium text-sm"
        >
          {showForm ? 'Cancel' : '+ Add Admin'}
        </button>
      </div>

      {/* Create / Edit Form */}
      {(showForm || editing) && (
        <form onSubmit={editing ? handleUpdate : handleCreate} className="bg-white rounded-xl shadow-sm border p-6 space-y-4">
          <h3 className="font-semibold">{editing ? 'Edit Admin' : 'Create New Admin'}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Admin Name *</label>
              <input type="text" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" required />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email *</label>
              <input type="email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" required />
            </div>
            {!editing && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Password <span className="text-gray-400">(optional, random if empty)</span></label>
                <input type="text" value={form.password} onChange={e => setForm({ ...form, password: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" placeholder="Random if empty" />
              </div>
            )}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Store Name</label>
              <input type="text" value={form.store_name} onChange={e => setForm({ ...form, store_name: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" placeholder="Auto-generated" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
              <input type="text" value={form.address} onChange={e => setForm({ ...form, address: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
              <input type="text" value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Card Color</label>
              <input type="color" value={form.card_color_hex} onChange={e => setForm({ ...form, card_color_hex: e.target.value })}
                className="w-full h-10 px-3 py-2 border rounded-lg cursor-pointer" />
            </div>
          </div>
          <div className="flex gap-2">
            <button type="submit" className="px-6 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 font-medium text-sm">
              {editing ? 'Update' : 'Create'}
            </button>
            <button type="button" onClick={() => { setEditing(null); setShowForm(false) }}
              className="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 text-sm">
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Admins Table */}
      <div className="bg-white rounded-xl shadow-sm border overflow-hidden">
        {admins.length === 0 ? (
          <div className="p-8 text-center text-gray-400">No admins yet. Create one above.</div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Admin</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Store</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Address</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Phone</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Color</th>
                <th className="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {admins.map(a => (
                <tr key={a.user_id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <div className="font-medium text-gray-900">{a.name}</div>
                    <div className="text-sm text-gray-500">{a.email}</div>
                  </td>
                  <td className="px-4 py-3 text-gray-700">{a.store_name || '—'}</td>
                  <td className="px-4 py-3 text-gray-500 text-sm">{a.store_address || '—'}</td>
                  <td className="px-4 py-3 text-gray-500 text-sm">{a.store_phone || '—'}</td>
                  <td className="px-4 py-3">
                    <span className="inline-block w-5 h-5 rounded-full border" style={{ backgroundColor: a.store_color || '#1E40AF' }} />
                    <span className="ml-1 text-xs text-gray-500">{a.store_color || ''}</span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button onClick={() => openEdit(a)} className="text-brand-600 hover:text-brand-800 text-sm font-medium mr-3">Edit</button>
                    <button onClick={() => handleDelete(a.user_id)} className="text-red-600 hover:text-red-800 text-sm font-medium">Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
