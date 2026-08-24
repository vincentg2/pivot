import { defineStore } from 'pinia'
import { api } from '@/lib/api'

export interface Club {
  id: string
  name: string
  shortName: string
  tla: string
  crestUrl: string | null
  websiteUrl: string | null
  venue: string
  competitions?: Competition[]
}

export interface Competition {
  id: string
  code: string
  name: string
  country: string
  emblemUrl: string | null
  clubCount: number
}

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({ clubs: [] as Club[], ready: false, saving: false }),
  getters: {
    clubIds: (state) => state.clubs.map((club) => club.id),
    has: (state) => (id: string) => state.clubs.some((club) => club.id === id),
  },
  actions: {
    async load() {
      this.clubs = (await api<{ favorites: Club[] }>('/favorites')).favorites
      this.ready = true
    },
    async replace(ids: string[]) {
      this.saving = true
      try {
        this.clubs = (
          await api<{ favorites: Club[] }>('/favorites', {
            method: 'PUT',
            body: JSON.stringify({ clubIds: ids }),
          })
        ).favorites
      } finally {
        this.saving = false
      }
    },
    async toggle(club: Club) {
      if (this.has(club.id)) {
        await this.replace(this.clubIds.filter((id) => id !== club.id))
        return
      }
      if (this.clubs.length >= 5) throw new Error('You can follow up to five clubs.')
      await this.replace([...this.clubIds, club.id])
    },
  },
})
