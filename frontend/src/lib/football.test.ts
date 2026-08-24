import { describe, expect, it } from 'vitest'
import { addDays, localDate } from './football'

describe('football date helpers', () => {
  it('keeps date navigation in local calendar time', () => {
    const start = new Date(2026, 11, 31, 12)
    expect(localDate(addDays(start, 1))).toBe('2027-01-01')
  })
})
