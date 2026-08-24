import { expect, test, type Page } from '@playwright/test'

const member = {
  id: 'user-1',
  email: 'camille@example.test',
  nickname: 'Camille',
  avatarSeed: 'seed',
  theme: 'system',
  role: 'user',
  createdAt: new Date().toISOString(),
}

async function mockMemberSession(page: Page) {
  await page.route('**/api/v1/setup/status', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ setupRequired: false, setupTokenConfigured: true }),
    }),
  )
  await page.route('**/api/v1/auth/me', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ user: member }),
    }),
  )
}

test('a member adds a club to their favorites', async ({ page }) => {
  await mockMemberSession(page)
  const club = {
    id: 'club-1',
    name: 'Paris Saint-Germain',
    shortName: 'PSG',
    tla: 'PSG',
    crestUrl: null,
    websiteUrl: null,
    venue: 'Parc des Princes',
  }
  await page.route('**/api/v1/competitions', (route) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: '{"competitions":[]}' }),
  )
  await page.route('**/api/v1/clubs**', (route) =>
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ clubs: [club] }),
    }),
  )
  await page.route('**/api/v1/favorites', async (route) => {
    if (route.request().method() === 'PUT') {
      expect(await route.request().postDataJSON()).toEqual({ clubIds: ['club-1'] })
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ favorites: [club] }),
      })
    }
    return route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ favorites: [] }),
    })
  })

  await page.goto('/clubs')
  await page.getByRole('button', { name: 'Add Paris Saint-Germain to favorites' }).click()
  await expect(
    page.getByRole('button', { name: 'Remove Paris Saint-Germain from favorites' }),
  ).toBeVisible()
  await expect(page.getByText('/ 5 favoris')).toContainText('1')
})

test('a member confirms permanent account deletion', async ({ page }) => {
  await mockMemberSession(page)
  let deleted = false
  await page.route('**/api/v1/profile', (route) => {
    if (route.request().method() === 'DELETE') {
      deleted = true
      return route.fulfill({ status: 204 })
    }
    return route.fallback()
  })

  await page.goto('/profile')
  await page.getByRole('button', { name: 'Supprimer le compte' }).click()
  await expect(page.getByRole('dialog')).toContainText('définitive et irréversible')
  await page.getByRole('button', { name: 'Supprimer définitivement' }).click()
  await expect(page).toHaveURL(/\/login$/)
  expect(deleted).toBe(true)
})
