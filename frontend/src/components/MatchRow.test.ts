import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import MatchRow from './MatchRow.vue'
import type { FootballMatch } from '@/lib/football'

const match: FootballMatch = {
  id: 'match-1',
  utcDate: '2026-08-28T18:45:00Z',
  status: 'TIMED',
  stage: null,
  matchday: 1,
  homeScore: null,
  awayScore: null,
  favorite: true,
  competition: { code: 'FL1', name: 'Ligue 1' },
  home: { id: 'club-1', name: 'Lille', shortName: 'Lille', tla: 'LIL', crestUrl: null },
  away: { id: 'club-2', name: 'Paris Saint-Germain', shortName: 'PSG', tla: 'PSG', crestUrl: null },
}

describe('MatchRow TV status', () => {
  it('shows the linked TV channel', () => {
    const wrapper = mount(MatchRow, { props: { match, channels: ['Canal+ Foot'] } })
    expect(wrapper.get('[aria-label="TV channels"]').text()).toContain('Canal+ Foot')
  })

  it('makes a pending TV listing explicit on the dashboard', () => {
    const wrapper = mount(MatchRow, { props: { match, showChannelStatus: true } })
    expect(wrapper.get('[aria-label="TV channels"]').text()).toContain(
      'TV channel not announced yet',
    )
  })
})
