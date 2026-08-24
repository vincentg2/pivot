import { describe, expect, it } from 'vitest'
import { favoriteMatchesWithin, pageItems } from './dashboard'
import type { FootballMatch } from './football'

function match(id: string, day: number, favorite = true): FootballMatch {
  return {
    id,
    utcDate: `2026-08-${String(day).padStart(2, '0')}T18:00:00Z`,
    status: 'TIMED',
    stage: '',
    matchday: 1,
    homeScore: null,
    awayScore: null,
    favorite,
    competition: { id: 'competition', code: 'FL1', name: 'Ligue 1' },
    home: { id: 'home', name: 'Home', shortName: 'Home', tla: 'HOM', crestUrl: null },
    away: { id: 'away', name: 'Away', shortName: 'Away', tla: 'AWY', crestUrl: null },
  }
}

describe('dashboard match window', () => {
  it('filters favorite matches by horizon and sorts them', () => {
    const now = new Date('2026-08-24T10:00:00Z')
    const matches = [
      match('later', 31),
      match('outside', 20),
      match('soon', 25),
      match('other', 26, false),
    ]

    expect(favoriteMatchesWithin(matches, now, 7).map((item) => item.id)).toEqual(['soon', 'later'])
  })

  it('paginates a match collection', () => {
    expect(pageItems([1, 2, 3, 4, 5, 6], 1, 5)).toEqual([6])
  })
})
