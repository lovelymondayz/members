import { useState, useEffect } from 'react'
import { getStores, updateStore, deleteStore, Store } from '../api/admin'

export default function StoresPage() {
  const [stores, setStores] = useState<Store[]>([])
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState<Store | null>(null)
  const [form, setForm] = useState({ name: '', address: '', phone: '', card_color_hex: '#1E40AF' })

  useEffect(() => { fetchStores() }, [])

  const fetchStores = async () => {
    try {
      const data = await getStores()
      setStores(data)
    } catch { alert('Failed to load stores') }
    finally { setLoading(false) }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editing) return
    try {
      await updateStore(editing.store_id, form)
      setEditing(null)
      fetchStores()
    } catch { alert('Failed to update store') }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this store and all its data (members, invoices)?')) return
    try { await deleteStore(id); fetchStores() }
    catch { alert('Failed to delete store') }
  }

  const openEdit = (s: Store) => {
    setEditing(s)
    setForm({ name: s.name, address: s.address || '', phone: s.phone || '', card_color_hex: s.card_color_hex })
  }

  if (loading) return <div className="p-6 text-gray-500">Loading...</div>

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">Stores</h1>

      {/* Edit Form */}
      {editing && (
        <form onSubmit={handleUpdate} className="bg-white rounded-xl shadow-sm border p-6 space-y-4">
          <h3 className="font-semibold">Edit: {editing.name}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Store Name *</label>
              <input type="text" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-brand-500" required />
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
            <button type="submit" className="px-6 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 font-medium text-sm">Save</button>
            <button type="button" onClick={() => setEditing(null)}
            className="px-2 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 text-sm">Cancel</button>
          </div>
        </form>
      )}

      {/* Stores Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {stores.map(s => (
          <div key={s.store_id} className="bg-white rounded-xl shadow-sm border p-5">
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg flex items-center justify-center text-white font-bold text-sm"
                  style={{ backgroundColor: s.card_color_hex }}>
                  {s.name.charAt(0).toUpperCase()}
                </div>
                <div>
                  <div className="font-semibold text-gray-900">{s.name}</div>
                  <div className="text-xs text-gray-500 font-mono">{s.store_id.slice(0, 8)}...</div>
                </div>
              </div>
            </div>
            {s.address && <p className="text-sm text-gray-600 mb-1">📍 {s.address}</p>}
            {s.phone && <p className="text-sm text-gray-600 mb-1">📞 {s.phone}</p>}
            <div className="flex items-center gap-2 mt-3">
              <span className="inline-block w-4 h-4 rounded-full border" style={{ backgroundColor: s.card_color_hex }} />
              <span className="text-xs text-gray-500">{s.card_color_hex}</span>
            </div>
            <div className="flex gap-2 mt-4 pt-3 border-t">
              <button onClick={() => openEdit(s)} className="text-brand-600 hover:text-brand-800 text-sm font-medium">Edit</button>
              <button onClick={() => handleDelete(s.store_id)} className="text-red-600 hover:text-red-800 text-sm font-medium">Delete</button>
            </div>
          </div>
        ))}
      </div>

      {stores.length === 0 && (
        <div className="text-center text-gray-400 py-8">No stores yet.</div>
      )}
    </div>
  )
}
