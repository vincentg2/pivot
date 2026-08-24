import { expect, test } from '@playwright/test'

test('an invited friend creates an account and reaches the private dashboard', async ({ page }) => {
	await page.route('**/api/v1/setup/status', (route) =>
	  route.fulfill({
		status: 200,
		contentType: 'application/json',
		body: JSON.stringify({ setupRequired: false, setupTokenConfigured: false }),
	  }),
	)
  const user = {
    id: 'user-1',
    email: 'camille@example.test',
    nickname: 'Camille',
    avatarSeed: 'seed',
    theme: 'system',
    role: 'user',
    createdAt: new Date().toISOString(),
  }
  await page.route('**/api/v1/auth/me', (route) =>
    route.fulfill({
      status: 401,
      contentType: 'application/problem+json',
      body: JSON.stringify({ type: 'about:blank', title: 'Authentication required', status: 401 }),
    }),
  )
  await page.route('**/api/v1/auth/register', (route) =>
    route.fulfill({ status: 201, contentType: 'application/json', body: JSON.stringify({ user }) }),
  )
  await page.route('**/api/v1/auth/login', (route) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ user }) }),
  )
  await page.goto('/register')
  await page.getByLabel('Invitation code').fill('PIVOT-EXAMPLECODE')
  await page.getByLabel('Nickname').fill('Camille')
  await page.getByLabel('Email').fill('camille@example.test')
  await page.getByLabel('Password').fill('correct horse battery staple')
  await page.getByRole('button', { name: 'Create account' }).click()
  await expect(page.getByRole('heading', { name: /Good evening, Camille/ })).toBeVisible()
})
