import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import AvatarMonogram from './AvatarMonogram.vue'

describe('AvatarMonogram', () => {
  it('derives stable initials without loading a remote image', () => {
    const wrapper = mount(AvatarMonogram, {
      props: { name: 'Camille Martin', seed: 'stable-seed' },
    })
    expect(wrapper.text()).toBe('CM')
    expect(wrapper.attributes('style')).toContain('--avatar-hue')
  })
})
