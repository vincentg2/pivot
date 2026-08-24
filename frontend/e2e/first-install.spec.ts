import { expect, test } from '@playwright/test'

test('the first administrator completes the installation wizard', async ({ page }) => {
  let installed = false
  await page.route('**/api/v1/setup/status', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ setupRequired: !installed, setupTokenConfigured: true }),
    }),
  )
  await page.route('**/api/v1/auth/me', (route) =>
    route.fulfill({
      status: 401,
      contentType: 'application/problem+json',
      body: JSON.stringify({ type: 'about:blank', title: 'Authentication required', status: 401 }),
    }),
  )
  await page.route('**/api/v1/setup', (route) => {
    installed = true
    return route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({ admin: { role: 'admin' } }),
    })
  })
  const user = {
    id: 'admin-1',
    email: 'owner@example.test',
    nickname: 'Owner',
    avatarSeed: 'seed',
    theme: 'system',
    role: 'admin',
    createdAt: new Date().toISOString(),
  }
  await page.route('**/api/v1/auth/login', (route) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ user }) }),
  )
  await page.route('**/api/v1/favorites', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ favorites: [] }),
    }),
  )
  await page.route('**/api/v1/matches**', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ matches: [] }),
    }),
  )
  await page.route('**/api/v1/news**', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ news: [] }),
    }),
  )
  await page.goto('/')
  await expect(page.getByRole('heading', { name: 'Premier administrateur' })).toBeVisible()
  await page.getByLabel('Jeton de configuration').fill('a-long-one-time-setup-token')
  await page.getByLabel('Pseudo').fill('Owner')
  await page.getByLabel('E-mail').fill('owner@example.test')
  await page.getByLabel('Mot de passe').fill('secure-password')
  await page.getByRole('button', { name: /Créer l’administrateur/ }).click()
  await expect(page.getByRole('heading', { name: /Bonsoir, Owner/ })).toBeVisible()
})
