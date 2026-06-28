import { Link } from 'react-router-dom'

export default function NotFoundPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center px-4">
      <div className="text-center">
        {/* Ghost emoji floating animation */}
        <div className="text-8xl mb-6 animate-bounce">👻</div>

        {/* 404 text */}
        <h1 className="text-7xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-pink-500 via-purple-500 to-cyan-500 mb-4">
          404
        </h1>

        {/* Message */}
        <p className="text-xl text-slate-300 mb-2">
          Halaman tidak ditemukan
        </p>
        <p className="text-slate-500 mb-8 max-w-md mx-auto">
          Seperti member yang hilang... halaman yang kamu cari mungkin sudah dipindahkan atau tidak pernah ada.
        </p>

        {/* Action buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            to="/dashboard"
            className="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold rounded-xl hover:from-purple-700 hover:to-pink-700 transition-all duration-200 shadow-lg shadow-purple-500/25 hover:shadow-purple-500/40"
          >
            🏠 Kembali ke Dashboard
          </Link>
          <Link
            to="/login"
            className="px-8 py-3 bg-slate-800 text-slate-300 font-semibold rounded-xl hover:bg-slate-700 transition-all duration-200 border border-slate-700"
          >
            🔑 Login
          </Link>
        </div>

        {/* Decorative dots */}
        <div className="mt-12 flex justify-center gap-2">
          <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse" />
          <div className="w-2 h-2 bg-pink-500 rounded-full animate-pulse" style={{ animationDelay: '0.2s' }} />
          <div className="w-2 h-2 bg-cyan-500 rounded-full animate-pulse" style={{ animationDelay: '0.4s' }} />
        </div>
      </div>
    </div>
  )
}
