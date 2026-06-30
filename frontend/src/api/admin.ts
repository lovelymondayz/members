import { api } from './client'

export interface AdminUser {
  user_id: string
  name: string
  email: string
  role_id: number
  role: string
  is_active: boolean
  store_id: string
  store_name: string
  store_address: string
  store_phone: string
  store_color: string
  created_at: string
}

export interface Store {
  store_id: string
  admin_id: string
  name: string
  logo_url?: string
  address?: string
  phone?: string
  card_color_hex: string
  created_at: string
  updated_at: string
}

// ─── Admin (Super Admin only) ───────────────────────────

export async function getAdmins(): Promise<AdminUser[]> {
  const res = await api.get('/admin/admins')
  return res.data
}

export async function createAdmin(data: {
  name: string
  email: string
  password?: string
  store_name?: string
  address?: string
  phone?: string
  card_color_hex?: string
}) {
  const res = await api.post('/admin/admins', data)
  return res.data
}

export async function updateAdmin(id: string, data: { name?: string; email?: string }) {
  const res = await api.put(`/admin/admins/${id}`, data)
  return res.data
}

export async function deleteAdmin(id: string) {
  const res = await api.delete(`/admin/admins/${id}`)
  return res.data
}

// ─── Stores (Super Admin only) ──────────────────────────

export async function getStores(): Promise<Store[]> {
  const res = await api.get('/stores')
  return res.data
}

export async function updateStore(id: string, data: {
  name?: string
  address?: string
  phone?: string
  card_color_hex?: string
  logo_url?: string
}) {
  const res = await api.put(`/stores/${id}`, data)
  return res.data
}

export async function deleteStore(id: string) {
  const res = await api.delete(`/stores/${id}`)
  return res.data
}
