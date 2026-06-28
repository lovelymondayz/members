import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../api/client'

interface Invoice {
  invoice_id: string
  invoice_number: string
  amount: number
  status: 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled'
  description?: string
  due_date?: string
  member?: { member_code: string; user_name?: string }
}

interface Member {
  member_id: string
  member_code: string
  user?: { name: string }
}

const statusColors: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-700',
  sent: 'bg-blue-100 text-blue-700',
  paid: 'bg-green-100 text-green-700',
  overdue: 'bg-red-100 text-red-700',
  cancelled: 'bg-gray-100 text-gray-500'
}

export default function InvoicesPage() {
  const [showCreate, setShowCreate] = useState(false)
  const [memberID, setMemberID] = useState('')
  const [amount, setAmount] = useState('')
  const [dueDate, setDueDate] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [payingID, setPayingID] = useState<string | null>(null)
  const [payAmount, setPayAmount] = useState('')
  const qc = useQueryClient()

  const { data: invoices, isLoading } = useQuery({
    queryKey: ['invoices'],
    queryFn: async () => {
      const res = await api.get('/invoices')
      return res.data as Invoice[]
    },
  })

  const { data: members } = useQuery({
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
      await api.post('/invoices', {
        member_id: memberID,
        amount: parseFloat(amount),
        due_date: dueDate,
        description,
      })
      setMemberID('')
      setAmount('')
      setDueDate('')
      setDescription('')
      setShowCreate(false)
      qc.invalidateQueries({ queryKey: ['invoices'] })
    } catch {
      alert('Failed to create invoice')
    } finally {
      setSubmitting(false)
    }
  }

  const handlePay = async (invoiceID: string) => {
    if (!payAmount) return
    try {
      await api.post(`/invoices/${invoiceID}/pay`, { amount: parseFloat(payAmount) })
      setPayingID(null)
      setPayAmount('')
      qc.invalidateQueries({ queryKey: ['invoices'] })
    } catch {
      alert('Failed to record payment')
    }
  }

  if (isLoading) return <div className="text-gray-500">Loading invoices...</div>

  return (
    <div>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Invoices</h1>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="bg-brand-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-brand-700"
        >
          + Create Invoice
        </button>
      </div>

      {showCreate && (
        <div className="mt-6 bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">New Invoice</h2>
          <form onSubmit={handleCreate} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Member</label>
              <select
                value={memberID}
                onChange={(e) => setMemberID(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                required
              >
                <option value="">Select member...</option>
                {members?.map((m) => (
                  <option key={m.member_id} value={m.member_id}>
                    {m.user?.name || m.member_code}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Amount (Rp)</label>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="500000"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
              <input
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <input
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Monthly membership fee"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
              />
            </div>
            <div className="md:col-span-2 flex gap-3">
              <button
                type="submit"
                disabled={submitting}
                className="bg-brand-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-brand-700 disabled:opacity-50"
              >
                {submitting ? 'Creating...' : 'Create Invoice'}
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
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Invoice #</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Member</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Amount</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="hidden md:table-cell text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Due Date</th>
              <th className="text-left px-6 py-3 text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {invoices?.map((inv) => (
              <tr key={inv.invoice_id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-mono text-gray-600">{inv.invoice_number}</td>
                <td className="px-6 py-4 text-sm font-medium text-gray-900">
                  {inv.member?.user_name || inv.member?.member_code || '—'}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900">Rp {inv.amount.toLocaleString('id-ID')}</td>
                <td className="px-6 py-4">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${statusColors[inv.status]}`}>
                    {inv.status}
                  </span>
                </td>
                <td className="hidden md:table-cell px-6 py-4 text-sm text-gray-500">
                  {inv.due_date ? new Date(inv.due_date).toLocaleDateString('id-ID') : '—'}
                </td>
                <td className="px-6 py-4">
                  {inv.status !== 'paid' && inv.status !== 'cancelled' && (
                    <>
                      {payingID === inv.invoice_id ? (
                        <div className="flex items-center gap-2">
                          <input
                            type="number"
                            value={payAmount}
                            onChange={(e) => setPayAmount(e.target.value)}
                            placeholder="Amount"
                            className="w-24 px-2 py-1 border border-gray-300 rounded text-xs"
                          />
                          <button
                            onClick={() => handlePay(inv.invoice_id)}
                            className="text-green-600 hover:text-green-700 text-xs font-medium"
                          >
                            Confirm
                          </button>
                          <button
                            onClick={() => { setPayingID(null); setPayAmount('') }}
                            className="text-gray-400 hover:text-gray-600 text-xs"
                          >
                            ✕
                          </button>
                        </div>
                      ) : (
                        <button
                          onClick={() => { setPayingID(inv.invoice_id); setPayAmount(String(inv.amount)) }}
                          className="text-green-600 hover:text-green-700 text-sm font-medium"
                        >
                          Record Payment
                        </button>
                      )}
                    </>
                  )}
                </td>
              </tr>
            ))}
            {(!invoices || invoices.length === 0) && (
              <tr>
                <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                  No invoices yet. Create your first invoice to get started.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
