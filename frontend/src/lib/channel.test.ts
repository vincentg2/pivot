import { describe, expect, it } from 'vitest'
import { channelIconUrl } from './channel'

describe('channelIconUrl', () => {
  it('maps channel variants to a shared remote brand mark', () => {
    expect(channelIconUrl('Canal+ Foot')).toContain('/CanalPlus.svg')
    expect(channelIconUrl('DAZN 1')).toContain('/dazn.svg')
    expect(channelIconUrl('Disney+')).toContain('/Disney%2B_logo.svg')
  })

  it('leaves unknown channels on the text fallback', () => {
    expect(channelIconUrl('Ligue 1+ add. 4')).toBeNull()
  })
})
