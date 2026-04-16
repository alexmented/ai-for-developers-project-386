import { expect, test } from '@playwright/test'

test.describe('Calendar booking flows', () => {
  test('GS-01 + GS-02: guest sees event types and opens booking page', async ({ page }) => {
    await page.goto('/name-owner')

    await expect(page.getByRole('heading', { name: 'Выберите тип события' })).toBeVisible()
    await expect(page.getByText('Встреча 15 минут').first()).toBeVisible()
    await expect(page.getByText('Встреча 30 минут').first()).toBeVisible()

    await page.getByRole('button', { name: /Встреча 15 минут/i }).first().click()

    await expect(page.getByRole('heading', { name: 'Календарь' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Статус слотов' })).toBeVisible()
  })

  test('GS-03: guest books a free slot', async ({ page }) => {
    await page.goto('/name-owner/meeting-15')

    const freeSlot = page.getByRole('button').filter({ hasText: 'Свободно' }).first()
    await expect(freeSlot).toBeVisible()
    await freeSlot.click()

    await page.getByRole('button', { name: 'Продолжить' }).click()

    // Guest form step
    await expect(page.getByText('Ваши данные')).toBeVisible()
    await page.getByPlaceholder('Ваше имя').fill('E2E Guest')
    await page.getByPlaceholder('Ваш email').fill('e2e@example.com')
    await page.getByRole('button', { name: 'Забронировать' }).click()

    await expect(page.getByText('Бронирование создано')).toBeVisible()
  })

  test('GS-04: second guest gets conflict on same slot', async ({ request }) => {
    const slotsResponse = await request.get('http://127.0.0.1:4020/public/name-owner/event-types/meeting-30/slots')
    expect(slotsResponse.ok()).toBeTruthy()

    const slots = (await slotsResponse.json()) as Array<{
      startAt: string
      isAvailable: boolean
    }>

    const freeSlot = slots.find((slot) => slot.isAvailable)
    expect(freeSlot).toBeTruthy()

    const firstResponse = await request.post('http://127.0.0.1:4020/public/name-owner/bookings', {
      data: {
        eventTypeId: 'meeting-30',
        slotStartAt: freeSlot!.startAt,
        guestName: 'first-guest',
        guestEmail: 'first@example.com',
      },
    })
    expect(firstResponse.ok()).toBeTruthy()

    const secondResponse = await request.post('http://127.0.0.1:4020/public/name-owner/bookings', {
      data: {
        eventTypeId: 'meeting-30',
        slotStartAt: freeSlot!.startAt,
        guestName: 'second-guest',
        guestEmail: 'second@example.com',
      },
    })

    expect(secondResponse.status()).toBe(409)
    await expect(secondResponse.json()).resolves.toMatchObject({
      code: 'SLOT_CONFLICT',
    })
  })

  test('AS-01 + AS-02: owner creates a new event type in admin', async ({ page }) => {
    const suffix = Date.now()
    const name = `Встреча ${suffix}`

    await page.goto('/admin')

    await expect(page.getByRole('heading', { name: 'Админка' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Типы событий' })).toBeVisible()

    await page.getByPlaceholder('Название').fill(name)
    await page.getByPlaceholder('Длительность, мин').fill('45')
    await page.getByPlaceholder('Описание').fill('Тип события для e2e проверки')
    await page.getByRole('button', { name: 'Создать' }).click()

    await expect(page.getByText(name)).toBeVisible()
  })

  test('AS-03: owner sees upcoming bookings in admin table', async ({ page }) => {
    await page.goto('/name-owner/meeting-15')

    const freeSlot = page.getByRole('button').filter({ hasText: 'Свободно' }).first()
    await expect(freeSlot).toBeVisible()
    await freeSlot.click()
    await page.getByRole('button', { name: 'Продолжить' }).click()

    // Guest form step
    await expect(page.getByText('Ваши данные')).toBeVisible()
    await page.getByPlaceholder('Ваше имя').fill('E2E Guest')
    await page.getByPlaceholder('Ваш email').fill('e2e@example.com')
    await page.getByRole('button', { name: 'Забронировать' }).click()

    await expect(page.getByText('Бронирование создано')).toBeVisible()

    await page.goto('/admin')

    await expect(page.getByRole('heading', { name: 'Предстоящие встречи' })).toBeVisible()
    await expect(page.getByRole('cell', { name: /E2E Guest/ }).first()).toBeVisible()
  })
})
