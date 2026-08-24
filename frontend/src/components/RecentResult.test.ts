import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import RecentResult from './RecentResult.vue'
import { i18n } from '@/i18n'
import type { FootballMatch } from '@/lib/football'

const result: FootballMatch = {
  id: 'result-1',
  utcDate: '2026-08-23T19:00:00Z',
  status: 'FINISHED',
  stage: null,
  matchday: 1,
  homeScore: 2,
  awayScore: 1,
  favorite: true,
  competition: { code: 'FL1', name: 'Ligue 1' },
  home: { id: 'club-1', name: 'Paris', shortName: 'PSG', tla: 'PSG', crestUrl: null },
  away: { id: 'club-2', name: 'Rennes', shortName: 'Rennes', tla: 'REN', crestUrl: null },
  goals: [
    { minute: 12, injuryTime: null, type: 'REGULAR', scorerName: 'A. Striker' },
    { minute: 90, injuryTime: 3, type: 'PENALTY', scorerName: 'B. Forward' },
  ],
}

describe('RecentResult', () => {
  it('shows scorer minutes and goal type when supplied by the provider', () => {
    const wrapper = mount(RecentResult, {
      props: { match: result },
      global: { plugins: [i18n] },
    })
    expect(wrapper.get('.recent-scorers').text()).toBe('A. Striker 12′ · B. Forward 90+3′ (pen.)')
  })
})
