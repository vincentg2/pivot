export interface ClubRef {
  id: string | null
  name: string
  shortName: string
  tla: string
  crestUrl: string | null
}
export interface MatchCompetition {
  id: string
  code: string
  name: string
}
export interface FootballMatch {
  id: string
  competition: MatchCompetition
  utcDate: string
  status: string
  stage: string
  matchday: number | null
  home: ClubRef
  away: ClubRef
  homeScore: number | null
  awayScore: number | null
  favorite: boolean
}
export interface StandingRow {
  position: number
  club: ClubRef
  played: number
  won: number
  drawn: number
  lost: number
  goalsFor: number
  goalsAgainst: number
  goalDifference: number
  points: number
}
export interface Standing {
  competition: MatchCompetition
  season: { startDate: string; endDate: string; current: boolean }
  rows: StandingRow[]
}

export function localDate(value: Date): string {
  const year = value.getFullYear()
  const month = String(value.getMonth() + 1).padStart(2, '0')
  const day = String(value.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
export function addDays(value: Date, days: number): Date {
  const result = new Date(value)
  result.setDate(result.getDate() + days)
  return result
}
