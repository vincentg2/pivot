import { defineStore } from 'pinia'
import { api } from '@/lib/api'

export type Theme = 'system' | 'light' | 'dark'
export interface User {
  id: string
  email: string
  nickname: string
  avatarSeed: string
  theme: Theme
  role: 'user' | 'admin'
  createdAt: string
}

function applyTheme(theme: Theme) {
  const dark =
    theme === 'dark' || (theme === 'system' && matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', dark)
  document.documentElement.dataset.theme = theme
}

export const useSessionStore = defineStore('session', {
  state: () => ({ user: null as User | null, ready: false }),
  actions: {
    async restore() {
      try {
        this.user = (await api<{ user: User }>('/auth/me')).user
        applyTheme(this.user.theme)
      } catch {
        this.user = null
        applyTheme('system')
      } finally {
        this.ready = true
      }
    },
    async login(email: string, password: string) {
      this.user = (
        await api<{ user: User }>('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ email, password }),
        })
      ).user
      applyTheme(this.user.theme)
    },
    async register(payload: {
      email: string
      password: string
      nickname: string
      invitationCode: string
    }) {
      await api('/auth/register', { method: 'POST', body: JSON.stringify(payload) })
      await this.login(payload.email, payload.password)
    },
    async logout() {
      await api('/auth/logout', { method: 'POST' })
      this.user = null
      applyTheme('system')
    },
    async updateProfile(nickname: string, theme: Theme) {
      this.user = (
        await api<{ user: User }>('/profile', {
          method: 'PATCH',
          body: JSON.stringify({ nickname, theme }),
        })
      ).user
      applyTheme(theme)
    },
    async deleteAccount() {
      await api('/profile', { method: 'DELETE' })
      this.user = null
      applyTheme('system')
    },
  },
})
