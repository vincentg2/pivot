import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import MatchRow from './MatchRow.vue'
import type { FootballMatch } from '@/lib/football'
import { i18n } from '@/i18n'

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
  const mountRow = (props: InstanceType<typeof MatchRow>['$props']) =>
    mount(MatchRow, { props, global: { plugins: [i18n] } })

  it('shows the linked TV channel', () => {
    const wrapper = mountRow({ match, channels: ['Canal+ Foot'] })
    expect(wrapper.get('[aria-label="TV channels"]').text()).toContain('Canal+ Foot')
  })

  it('makes a pending TV listing explicit on the dashboard', () => {
    const wrapper = mountRow({ match, showChannelStatus: true })
    expect(wrapper.get('[aria-label="TV channels"]').text()).toContain(
      'Chaîne TV pas encore annoncée',
    )
  })

  it('shows a compact date when requested by the dashboard', () => {
    const wrapper = mountRow({ match, showDate: true })
    expect(wrapper.get('time').text()).toContain('28 août')
  })
})
