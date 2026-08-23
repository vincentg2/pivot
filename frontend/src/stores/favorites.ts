import { defineStore } from 'pinia'

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({ clubIds: [] as string[] }),
  actions: {
    replace(ids: string[]) {
      this.clubIds = ids.slice(0, 5)
    },
  },
})
