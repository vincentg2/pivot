import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { useFavoritesStore, type Club } from './favorites'

function club(id: string): Club {
  return {
    id,
    name: `Club ${id}`,
    shortName: `Club ${id}`,
    tla: id,
    crestUrl: null,
    websiteUrl: null,
    venue: '',
  }
}

describe('favorites store', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('enforces the five-club limit before calling the API', async () => {
    const store = useFavoritesStore()
    store.clubs = ['A', 'B', 'C', 'D', 'E'].map(club)

    await expect(store.toggle(club('F'))).rejects.toThrow('up to five clubs')
    expect(store.clubIds).toEqual(['A', 'B', 'C', 'D', 'E'])
  })
})
