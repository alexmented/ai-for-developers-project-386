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
        eventTypeId: 'meeting-15',
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
    const id = `meeting-${suffix}`
    const name = `Встреча ${suffix}`

    await page.goto('/admin')

    await expect(page.getByRole('heading', { name: 'Админка' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Типы событий' })).toBeVisible()

    await page.getByPlaceholder('id').fill(id)
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
    await expect(page.getByText('Бронирование создано')).toBeVisible()

    await page.goto('/admin')

    await expect(page.getByRole('heading', { name: 'Предстоящие встречи' })).toBeVisible()
    await expect(page.getByRole('cell', { name: 'Guest User (guest@example.com)' }).first()).toBeVisible()
  })

  test('AS-04: owner updates profile fields and sees them on public page', async ({ page, request }) => {
    const suffix = Date.now()
    const displayName = `Alex ${suffix}`
    const email = `alex.${suffix}@example.com`
    const photoUrl = 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=300&q=80'

    await expect
      .poll(async () => {
        const response = await request.get('http://127.0.0.1:4020/admin/owner')
        return response.status()
      })
      .toBe(200)

    await page.goto('/admin')

    await page.getByPlaceholder('URL фото').fill(photoUrl)
    await page.getByPlaceholder('Имя').fill(displayName)
    await page.getByPlaceholder('Email').fill(email)
    await page.getByPlaceholder('Часовой пояс (например, Europe/Moscow)').fill('UTC')
    await page.getByPlaceholder('Начало рабочего дня (час)').fill('10')
    await page.getByPlaceholder('Конец рабочего дня (час)').fill('16')
    await page.getByRole('button', { name: 'Сохранить профиль' }).click()

    await expect(page.getByText('Профиль сохранён')).toBeVisible()

    await page.goto('/name-owner')

    await expect(page.getByRole('heading', { name: displayName })).toBeVisible()
    await expect(page.getByText('UTC · Рабочие часы: 10:00–16:00')).toBeVisible()
    await expect(page.getByAltText('Фото владельца')).toHaveAttribute('src', photoUrl)
  })
})
