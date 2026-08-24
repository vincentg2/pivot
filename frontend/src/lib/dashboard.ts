import { addDays, type FootballMatch } from './football'

export function favoriteMatchesWithin(
  matches: FootballMatch[],
  now: Date,
  days: number,
): FootballMatch[] {
  const start = new Date(now)
  start.setHours(0, 0, 0, 0)
  const end = addDays(start, days)
  end.setHours(23, 59, 59, 999)
  return matches
    .filter((match) => {
      const kickoff = new Date(match.utcDate)
      return match.favorite && kickoff >= start && kickoff <= end
    })
    .sort((left, right) => Date.parse(left.utcDate) - Date.parse(right.utcDate))
}

export function matchesForClub(matches: FootballMatch[], clubId: string | null): FootballMatch[] {
  if (!clubId) return matches
  return matches.filter((match) => match.home.id === clubId || match.away.id === clubId)
}

export function latestFavoriteResults(matches: FootballMatch[], limit = 3): FootballMatch[] {
  return matches
    .filter((match) => match.favorite && match.status === 'FINISHED')
    .sort((left, right) => Date.parse(right.utcDate) - Date.parse(left.utcDate))
    .slice(0, limit)
}
