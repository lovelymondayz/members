import { useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { api } from '../api/client'

interface CardData {
  member_id: string
  member_code: string
  name: string
  tier: string
  store_name: string
  store_logo?: string
  card_color: string
  joined_at: string
  qr_data: string
}

export default function MemberCardPage() {
  const { id } = useParams()

  const { data: card, isLoading } = useQuery({
    queryKey: ['member-card', id],
    queryFn: async () => {
      const res = await api.get(`/members/${id}/card`)
      return res.data as CardData
    },
    enabled: !!id,
  })

  if (isLoading) return <div className="text-gray-500">Loading card...</div>
  if (!card) return <div className="text-gray-500">Member not found</div>

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Membership Card</h1>

      <div
        className="rounded-2xl p-6 text-white shadow-xl"
        style={{ backgroundColor: card.card_color }}
      >
        <div className="flex items-center justify-between mb-8">
          <div>
            <p className="text-white/70 text-xs uppercase tracking-wide">{card.store_name}</p>
            <h2 className="text-xl font-bold mt-1">{card.name}</h2>
          </div>
          {card.store_logo && (
            <img src={card.store_logo} alt="Logo" className="w-12 h-12 rounded-full bg-white/20" />
          )}
        </div>

        <div className="flex items-end justify-between">
          <div>
            <p className="text-white/70 text-xs">Member Code</p>
            <p className="text-lg font-mono font-bold">{card.member_code}</p>
          </div>
          <span className="bg-white/20 px-3 py-1 rounded-full text-sm font-medium capitalize">
            {card.tier}
          </span>
        </div>

        <div className="mt-6 pt-4 border-t border-white/20">
          <p className="text-white/70 text-xs">Member since {card.joined_at}</p>
        </div>
      </div>

      {/* QR Code placeholder */}
      <div className="mt-6 bg-white rounded-xl border border-gray-200 p-6 text-center">
        <p className="text-sm text-gray-500 mb-4">QR Code</p>
        <div className="w-32 h-32 mx-auto bg-gray-100 rounded-lg flex items-center justify-center">
          <span className="text-xs text-gray-400">QR: {card.qr_data}</span>
        </div>
        <p className="text-xs text-gray-400 mt-2">Scan for quick verification</p>
      </div>
    </div>
  )
}
