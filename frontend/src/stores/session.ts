import { defineStore } from 'pinia'
import { api } from '@/lib/api'
import { applyLocale, type Locale } from '@/i18n'

export type Theme = 'system' | 'light' | 'dark'
export interface User {
  id: string
  email: string
  nickname: string
  avatarSeed: string
  theme: Theme
  locale: Locale
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
        applyLocale(this.user.locale ?? 'fr')
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
      applyLocale(this.user.locale ?? 'fr')
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
    async updateProfile(nickname: string, theme: Theme, locale: Locale) {
      this.user = (
        await api<{ user: User }>('/profile', {
          method: 'PATCH',
          body: JSON.stringify({ nickname, theme, locale }),
        })
      ).user
      applyTheme(theme)
      applyLocale(locale)
    },
    async deleteAccount() {
      await api('/profile', { method: 'DELETE' })
      this.user = null
      applyTheme('system')
    },
  },
})
