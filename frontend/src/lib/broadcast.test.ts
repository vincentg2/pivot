import { describe, expect, it } from 'vitest'
import { broadcastChannelsByMatch, type BroadcastListing } from './broadcast'

function listing(overrides: Partial<BroadcastListing>): BroadcastListing {
  return {
    id: 'listing-1',
    matchId: 'match-1',
    startsAt: '2026-08-24T19:00:00Z',
    homeName: 'Lille',
    awayName: 'PSG',
    label: 'Lille - PSG',
    competitionName: 'Ligue 1',
    kind: 'live',
    channels: ['Canal+ Foot'],
    sourceUrl: null,
    external: false,
    corrected: false,
    ...overrides,
  }
}

describe('broadcastChannelsByMatch', () => {
  it('groups and deduplicates channels for linked matches', () => {
    const channels = broadcastChannelsByMatch([
      listing({}),
      listing({ id: 'listing-2', channels: ['Canal+ Foot', 'DAZN'] }),
      listing({ id: 'external', matchId: null, channels: ['beIN Sports 1'] }),
    ])

    expect(channels.get('match-1')).toEqual(['Canal+ Foot', 'DAZN'])
    expect(channels.size).toBe(1)
  })
})
