import { describe, expect, it } from 'vitest'
import { favoriteMatchesWithin, latestFavoriteResults, matchesForClub } from './dashboard'
import type { FootballMatch } from './football'

function match(id: string, day: number, favorite = true, status = 'TIMED'): FootballMatch {
  return {
    id,
    utcDate: `2026-08-${String(day).padStart(2, '0')}T18:00:00Z`,
    status,
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

  it('returns the latest finished matches for favorite clubs', () => {
    const matches = [
      match('older', 10, true, 'FINISHED'),
      match('upcoming', 28),
      match('latest', 20, true, 'FINISHED'),
      match('not-favorite', 22, false, 'FINISHED'),
    ]

    expect(latestFavoriteResults(matches, 2).map((item) => item.id)).toEqual(['latest', 'older'])
  })

  it('filters the collection for one favorite club', () => {
    const homeMatch = match('home-match', 25)
    const awayMatch = match('away-match', 26)
    const unrelated = match('unrelated', 27)
    homeMatch.home.id = 'favorite'
    awayMatch.away.id = 'favorite'

    expect(
      matchesForClub([homeMatch, unrelated, awayMatch], 'favorite').map((item) => item.id),
    ).toEqual(['home-match', 'away-match'])
    expect(matchesForClub([homeMatch, unrelated], null)).toHaveLength(2)
  })
})
