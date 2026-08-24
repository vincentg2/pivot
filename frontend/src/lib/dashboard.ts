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

export function pageItems<T>(items: T[], page: number, pageSize: number): T[] {
  const start = Math.max(0, page) * pageSize
  return items.slice(start, start + pageSize)
}
