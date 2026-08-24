import { expect, test } from '@playwright/test'

test('an administrator issues a one-time link and the member resets their password', async ({
  page,
}) => {
  const resetToken = 'opaque-password-reset-token-for-camille'
  let resetConsumed = false
  await page.route('**/api/v1/setup/status', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ setupRequired: false, setupTokenConfigured: true }),
    }),
  )
  await page.route('**/api/v1/auth/me', (route) => {
    if (resetConsumed) {
      return route.fulfill({
        status: 401,
        contentType: 'application/problem+json',
        body: JSON.stringify({
          type: 'about:blank',
          title: 'Authentication required',
          status: 401,
        }),
      })
    }
    return route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: {
          id: 'admin-1',
          email: 'owner@example.test',
          nickname: 'Owner',
          avatarSeed: 'seed',
          theme: 'system',
          role: 'admin',
          createdAt: new Date().toISOString(),
        },
      }),
    })
  })
  await page.route('**/api/v1/admin/invitations', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ invitations: [] }),
    }),
  )
  await page.route('**/api/v1/admin/collections/**', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ enabled: false, latestRun: null }),
    }),
  )
  await page.route('**/api/v1/admin/password-resets', async (route) => {
    expect(await route.request().postDataJSON()).toEqual({ email: 'camille@example.test' })
    return route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({
        token: resetToken,
        expiresAt: new Date(Date.now() + 30 * 60 * 1000).toISOString(),
        user: { email: 'camille@example.test', nickname: 'Camille' },
      }),
    })
  })
  await page.route('**/api/v1/auth/password-reset', async (route) => {
    expect(await route.request().postDataJSON()).toEqual({
      token: resetToken,
      password: 'a-brand-new-safe-password',
    })
    resetConsumed = true
    return route.fulfill({ status: 204 })
  })

  await page.goto('/admin')
  await page.getByLabel('Member email').fill('camille@example.test')
  await page.getByRole('button', { name: 'Generate reset link' }).click()
  await expect(page.getByRole('status')).toContainText(`/reset-password?token=${resetToken}`)

  await page.goto(`/reset-password?token=${resetToken}`)
  await page.getByLabel('New password').fill('a-brand-new-safe-password')
  await page.getByLabel('Confirm password').fill('a-brand-new-safe-password')
  await page.getByRole('button', { name: 'Update password' }).click()
  await expect(page.getByRole('heading', { name: 'Password updated' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Return to sign in' })).toBeVisible()
})
