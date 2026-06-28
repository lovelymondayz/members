import axios from 'axios'

export const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Attach token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('members-auth')
    ? JSON.parse(localStorage.getItem('members-auth') || '{}').state?.token
    : null

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})
