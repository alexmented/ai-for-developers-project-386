export function buildBookingWindow(now = new Date()) {
  const fromDate = new Date(now)
  const toDate = new Date(now)
  toDate.setDate(toDate.getDate() + 14)

  return {
    from: fromDate.toISOString(),
    to: toDate.toISOString(),
  }
}

export function formatTimeRange(startAt: string, endAt: string, locale = 'ru-RU') {
  const start = new Date(startAt)
  const end = new Date(endAt)

  return `${start.toLocaleTimeString(locale, { hour: '2-digit', minute: '2-digit' })} - ${end.toLocaleTimeString(locale, { hour: '2-digit', minute: '2-digit' })}`
}

export function formatFullDate(value: string, locale = 'ru-RU') {
  return new Date(value).toLocaleDateString(locale, {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  })
}

export function getMonthLabel(value: string, locale = 'ru-RU') {
  return new Date(value).toLocaleDateString(locale, {
    month: 'long',
    year: 'numeric',
  })
}

export function toDateKey(value: string) {
  const date = new Date(value)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
}
