export interface TimezoneOption {
  value: string
  label: string
}

function gmtOffset(offsetHours: number): string {
  const sign = offsetHours >= 0 ? '+' : '-'
  const abs = Math.abs(offsetHours)
  const h = Math.floor(abs)
  const m = (abs - h) * 60
  const mm = m > 0 ? `:${String(m).padStart(2, '0')}` : ''
  return `GMT${sign}${h}${mm}`
}

const timezoneDefinitions: Array<{ value: string; offset: number }> = [
  { value: 'Pacific/Midway', offset: -11 },
  { value: 'Pacific/Honolulu', offset: -10 },
  { value: 'America/Anchorage', offset: -9 },
  { value: 'America/Los_Angeles', offset: -8 },
  { value: 'America/Denver', offset: -7 },
  { value: 'America/Chicago', offset: -6 },
  { value: 'America/New_York', offset: -5 },
  { value: 'America/Sao_Paulo', offset: -3 },
  { value: 'Atlantic/Reykjavik', offset: 0 },
  { value: 'Europe/London', offset: 0 },
  { value: 'Europe/Berlin', offset: 1 },
  { value: 'Europe/Kyiv', offset: 2 },
  { value: 'Europe/Moscow', offset: 3 },
  { value: 'Europe/Istanbul', offset: 3 },
  { value: 'Asia/Dubai', offset: 4 },
  { value: 'Asia/Karachi', offset: 5 },
  { value: 'Asia/Kolkata', offset: 5.5 },
  { value: 'Asia/Dhaka', offset: 6 },
  { value: 'Asia/Bangkok', offset: 7 },
  { value: 'Asia/Shanghai', offset: 8 },
  { value: 'Asia/Singapore', offset: 8 },
  { value: 'Asia/Tokyo', offset: 9 },
  { value: 'Asia/Seoul', offset: 9 },
  { value: 'Australia/Sydney', offset: 11 },
  { value: 'Pacific/Auckland', offset: 12 },
]

export const timezoneOptions: TimezoneOption[] = timezoneDefinitions.map((tz) => ({
  value: tz.value,
  label: `${gmtOffset(tz.offset)}  ${tz.value.replace(/_/g, ' ')}`,
}))
