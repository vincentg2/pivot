export interface BroadcastListing {
  id: string
  matchId: string | null
  startsAt: string
  homeName: string
  awayName: string
  label: string
  competitionName: string
  kind: 'live' | 'delayed' | 'replay'
  channels: string[]
  sourceUrl: string | null
  external: boolean
  corrected: boolean
  hidden?: boolean
}

export interface BroadcastAudit {
  id: string
  listingId: string
  adminId: string | null
  action: 'corrected' | 'cleared'
  createdAt: string
}

export function listingTime(value: string): string {
  return new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit' }).format(
    new Date(value),
  )
}

export function broadcastChannelsByMatch(listings: BroadcastListing[]): Map<string, string[]> {
  const channels = new Map<string, string[]>()
  for (const listing of listings) {
    if (!listing.matchId) continue
    channels.set(listing.matchId, [
      ...new Set([...(channels.get(listing.matchId) || []), ...listing.channels]),
    ])
  }
  return channels
}
