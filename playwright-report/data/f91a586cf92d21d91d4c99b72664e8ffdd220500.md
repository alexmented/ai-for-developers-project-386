# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: calendar.spec.ts >> Calendar booking flows >> AS-01 + AS-02: owner creates a new event type in admin
- Location: tests/e2e/calendar.spec.ts:66:3

# Error details

```
Error: locator.fill: Test ended.
Call log:
  - waiting for getByPlaceholder('id')

```

# Test source

```ts
  1   | import { expect, test } from '@playwright/test'
  2   | 
  3   | test.describe('Calendar booking flows', () => {
  4   |   test('GS-01 + GS-02: guest sees event types and opens booking page', async ({ page }) => {
  5   |     await page.goto('/name-owner')
  6   | 
  7   |     await expect(page.getByRole('heading', { name: 'Выберите тип события' })).toBeVisible()
  8   |     await expect(page.getByText('Встреча 15 минут').first()).toBeVisible()
  9   |     await expect(page.getByText('Встреча 30 минут').first()).toBeVisible()
  10  | 
  11  |     await page.getByRole('button', { name: /Встреча 15 минут/i }).first().click()
  12  | 
  13  |     await expect(page.getByRole('heading', { name: 'Календарь' })).toBeVisible()
  14  |     await expect(page.getByRole('heading', { name: 'Статус слотов' })).toBeVisible()
  15  |   })
  16  | 
  17  |   test('GS-03: guest books a free slot', async ({ page }) => {
  18  |     await page.goto('/name-owner/meeting-15')
  19  | 
  20  |     const freeSlot = page.getByRole('button').filter({ hasText: 'Свободно' }).first()
  21  |     await expect(freeSlot).toBeVisible()
  22  |     await freeSlot.click()
  23  | 
  24  |     await page.getByRole('button', { name: 'Продолжить' }).click()
  25  | 
  26  |     await expect(page.getByText('Бронирование создано')).toBeVisible()
  27  |   })
  28  | 
  29  |   test('GS-04: second guest gets conflict on same slot', async ({ request }) => {
  30  |     const slotsResponse = await request.get('http://127.0.0.1:4020/public/name-owner/event-types/meeting-30/slots')
  31  |     expect(slotsResponse.ok()).toBeTruthy()
  32  | 
  33  |     const slots = (await slotsResponse.json()) as Array<{
  34  |       startAt: string
  35  |       isAvailable: boolean
  36  |     }>
  37  | 
  38  |     const freeSlot = slots.find((slot) => slot.isAvailable)
  39  |     expect(freeSlot).toBeTruthy()
  40  | 
  41  |     const firstResponse = await request.post('http://127.0.0.1:4020/public/name-owner/bookings', {
  42  |       data: {
  43  |         eventTypeId: 'meeting-30',
  44  |         slotStartAt: freeSlot!.startAt,
  45  |         guestName: 'first-guest',
  46  |         guestEmail: 'first@example.com',
  47  |       },
  48  |     })
  49  |     expect(firstResponse.ok()).toBeTruthy()
  50  | 
  51  |     const secondResponse = await request.post('http://127.0.0.1:4020/public/name-owner/bookings', {
  52  |       data: {
  53  |         eventTypeId: 'meeting-15',
  54  |         slotStartAt: freeSlot!.startAt,
  55  |         guestName: 'second-guest',
  56  |         guestEmail: 'second@example.com',
  57  |       },
  58  |     })
  59  | 
  60  |     expect(secondResponse.status()).toBe(409)
  61  |     await expect(secondResponse.json()).resolves.toMatchObject({
  62  |       code: 'SLOT_CONFLICT',
  63  |     })
  64  |   })
  65  | 
  66  |   test('AS-01 + AS-02: owner creates a new event type in admin', async ({ page }) => {
  67  |     const suffix = Date.now()
  68  |     const id = `meeting-${suffix}`
  69  |     const name = `Встреча ${suffix}`
  70  | 
  71  |     await page.goto('/admin')
  72  | 
  73  |     await expect(page.getByRole('heading', { name: 'Админка' })).toBeVisible()
  74  |     await expect(page.getByRole('heading', { name: 'Типы событий' })).toBeVisible()
  75  | 
> 76  |     await page.getByPlaceholder('id').fill(id)
      |                                       ^ Error: locator.fill: Test ended.
  77  |     await page.getByPlaceholder('Название').fill(name)
  78  |     await page.getByPlaceholder('Длительность, мин').fill('45')
  79  |     await page.getByPlaceholder('Описание').fill('Тип события для e2e проверки')
  80  |     await page.getByRole('button', { name: 'Создать' }).click()
  81  | 
  82  |     await expect(page.getByText(name)).toBeVisible()
  83  |   })
  84  | 
  85  |   test('AS-03: owner sees upcoming bookings in admin table', async ({ page }) => {
  86  |     await page.goto('/name-owner/meeting-15')
  87  | 
  88  |     const freeSlot = page.getByRole('button').filter({ hasText: 'Свободно' }).first()
  89  |     await expect(freeSlot).toBeVisible()
  90  |     await freeSlot.click()
  91  |     await page.getByRole('button', { name: 'Продолжить' }).click()
  92  |     await expect(page.getByText('Бронирование создано')).toBeVisible()
  93  | 
  94  |     await page.goto('/admin')
  95  | 
  96  |     await expect(page.getByRole('heading', { name: 'Предстоящие встречи' })).toBeVisible()
  97  |     await expect(page.getByRole('cell', { name: 'Guest User (guest@example.com)' }).first()).toBeVisible()
  98  |   })
  99  | 
  100 |   test('AS-04: owner updates profile fields and sees them on public page', async ({ page, request }) => {
  101 |     const suffix = Date.now()
  102 |     const displayName = `Alex ${suffix}`
  103 |     const email = `alex.${suffix}@example.com`
  104 |     const photoUrl = 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=300&q=80'
  105 | 
  106 |     await expect
  107 |       .poll(async () => {
  108 |         const response = await request.get('http://127.0.0.1:4020/admin/owner')
  109 |         return response.status()
  110 |       })
  111 |       .toBe(200)
  112 | 
  113 |     await page.goto('/admin')
  114 | 
  115 |     await page.getByPlaceholder('URL фото').fill(photoUrl)
  116 |     await page.getByPlaceholder('Имя').fill(displayName)
  117 |     await page.getByPlaceholder('Email').fill(email)
  118 |     await page.getByPlaceholder('Часовой пояс (например, Europe/Moscow)').fill('UTC')
  119 |     await page.getByPlaceholder('Начало рабочего дня (час)').fill('10')
  120 |     await page.getByPlaceholder('Конец рабочего дня (час)').fill('16')
  121 |     await page.getByRole('button', { name: 'Сохранить профиль' }).click()
  122 | 
  123 |     await expect(page.getByText('Профиль сохранён')).toBeVisible()
  124 | 
  125 |     await page.goto('/name-owner')
  126 | 
  127 |     await expect(page.getByRole('heading', { name: displayName })).toBeVisible()
  128 |     await expect(page.getByText('UTC · Рабочие часы: 10:00–16:00')).toBeVisible()
  129 |     await expect(page.getByAltText('Фото владельца')).toHaveAttribute('src', photoUrl)
  130 |   })
  131 | })
  132 | 
```