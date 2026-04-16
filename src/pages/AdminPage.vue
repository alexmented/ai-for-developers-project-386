<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { UiButton } from '@/components/ui/button'
import { UiCard } from '@/components/ui/card'
import { UiInput } from '@/components/ui/input'
import { UiTextarea } from '@/components/ui/textarea'
import { UiTable } from '@/components/ui/table'
import AppErrorState from '@/components/common/AppErrorState.vue'
import { formatFullDate, formatTimeRange } from '@/lib/date'
import { timezoneOptions } from '@/lib/timezones'
import { api } from '@/services/api'
import type { Booking, CalendarOwner, DayOfWeek, EventType, WorkDaySchedule } from '@/types/api'

const orderedDays: DayOfWeek[] = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday']
const dayLabels: Record<DayOfWeek, string> = {
  monday: 'Пн',
  tuesday: 'Вт',
  wednesday: 'Ср',
  thursday: 'Чт',
  friday: 'Пт',
  saturday: 'Сб',
  sunday: 'Вс',
}

const eventTypes = ref<EventType[]>([])
const bookings = ref<Booking[]>([])
const owner = ref<CalendarOwner | null>(null)
const loading = ref(true)
const savingOwner = ref(false)
const pageError = ref<string | null>(null)
const ownerError = ref<string | null>(null)
const eventTypeError = ref<string | null>(null)
const ownerSaved = ref(false)

const ownerForm = reactive({
  photoUrl: '',
  displayName: '',
  email: '',
  timezone: 'Europe/Moscow',
  weeklySchedule: createDefaultWeeklySchedule(),
})

const form = reactive({
  name: '',
  description: '',
  durationMinutes: 30,
})

function createDefaultWeeklySchedule(): WorkDaySchedule[] {
  return orderedDays.map((dayOfWeek) => ({
    dayOfWeek,
    isActive: true,
    startHour: 9,
    endHour: 18,
  }))
}

function normalizeWeeklySchedule(schedule: WorkDaySchedule[] | undefined): WorkDaySchedule[] {
  if (!schedule || schedule.length === 0) {
    return createDefaultWeeklySchedule()
  }

  const byDay = new Map(schedule.map((item) => [item.dayOfWeek, item]))
  return orderedDays.map((dayOfWeek) => {
    const item = byDay.get(dayOfWeek)
    if (!item) {
      return { dayOfWeek, isActive: true, startHour: 9, endHour: 18 }
    }

    return {
      dayOfWeek,
      isActive: item.isActive,
      startHour: item.startHour,
      endHour: item.endHour,
    }
  })
}

async function loadAdminData() {
  loading.value = true
  pageError.value = null
  ownerError.value = null
  eventTypeError.value = null
  ownerSaved.value = false

  try {
    const [ownerResponse, eventTypesResponse, bookingsResponse] = await Promise.all([
      api.getAdminOwner(),
      api.getAdminEventTypes(),
      api.getAdminUpcomingBookings(),
    ])

    owner.value = ownerResponse
    ownerForm.photoUrl = ownerResponse.photoUrl ?? ''
    ownerForm.displayName = ownerResponse.displayName
    ownerForm.email = ownerResponse.email
    ownerForm.timezone = ownerResponse.timezone
    ownerForm.weeklySchedule = normalizeWeeklySchedule(ownerResponse.weeklySchedule)
    eventTypes.value = eventTypesResponse
    bookings.value = bookingsResponse
  } catch (err) {
    pageError.value = err instanceof Error ? err.message : 'Ошибка загрузки админки'
  } finally {
    loading.value = false
  }
}

async function updateOwner() {
  savingOwner.value = true
  ownerError.value = null
  ownerSaved.value = false

  try {
    owner.value = await api.updateAdminOwner({
      photoUrl: ownerForm.photoUrl.trim() || undefined,
      displayName: ownerForm.displayName,
      email: ownerForm.email,
      timezone: ownerForm.timezone,
      weeklySchedule: ownerForm.weeklySchedule.map((item) => ({
        dayOfWeek: item.dayOfWeek,
        isActive: item.isActive,
        startHour: Number(item.startHour),
        endHour: Number(item.endHour),
      })),
    })

    if (owner.value) {
      ownerForm.photoUrl = owner.value.photoUrl ?? ''
      ownerForm.displayName = owner.value.displayName
      ownerForm.email = owner.value.email
      ownerForm.timezone = owner.value.timezone
      ownerForm.weeklySchedule = normalizeWeeklySchedule(owner.value.weeklySchedule)
    }

    ownerSaved.value = true
  } catch (err) {
    ownerError.value = err instanceof Error ? err.message : 'Не удалось обновить профиль'
  } finally {
    savingOwner.value = false
  }
}

const ownerPhotoPreview = computed(() => ownerForm.photoUrl.trim())

async function createEventType() {
  ownerSaved.value = false
  eventTypeError.value = null
  try {
    const created = await api.createAdminEventType({
      name: form.name,
      description: form.description,
      durationMinutes: Number(form.durationMinutes),
    })

    eventTypes.value = [...eventTypes.value, created]
    form.name = ''
    form.description = ''
    form.durationMinutes = 30
  } catch (err) {
    eventTypeError.value = err instanceof Error ? err.message : 'Не удалось создать тип события'
  }
}

onMounted(loadAdminData)
</script>

<template>
  <main class="mx-auto max-w-6xl space-y-6 px-4 py-8 md:py-12">
    <h1 class="text-5xl font-semibold text-[#12213c]">Админка</h1>

    <section v-if="loading" class="text-[#6f84a1]">Загрузка...</section>
    <AppErrorState v-else-if="pageError" :message="pageError" title="Не удалось загрузить админку" />

    <template v-else>
      <UiCard class="p-6">
        <h2 class="text-3xl font-semibold text-[#1d2d4a]">Профиль владельца</h2>
        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div class="md:col-span-2 flex items-center gap-4">
            <img
              v-if="ownerPhotoPreview"
              :src="ownerPhotoPreview"
              alt="Фото владельца"
              class="h-16 w-16 rounded-full object-cover border border-[#d2dbe7]"
            />
            <div v-else class="grid h-16 w-16 place-items-center rounded-full bg-gradient-to-b from-[#ffcd9d] to-[#23a79a] text-xl">🧑</div>
            <p class="text-sm text-[#6f84a1]">Фото используется на публичной странице владельца</p>
          </div>

          <UiInput v-model="ownerForm.photoUrl" placeholder="URL фото" />
          <UiInput v-model="ownerForm.displayName" placeholder="Имя" />
          <UiInput v-model="ownerForm.email" placeholder="Email" type="email" />
          <select
            v-model="ownerForm.timezone"
            class="flex h-11 w-full rounded-xl border border-[#d3dbe7] bg-white px-3 py-2 text-sm text-[#13213a] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#d3dbe7]"
          >
            <option v-for="tz in timezoneOptions" :key="tz.value" :value="tz.value">{{ tz.label }}</option>
          </select>

          <div class="md:col-span-2 rounded-xl border border-[#d2dbe7] p-4">
            <p class="mb-3 text-sm font-semibold text-[#1d2d4a]">Доступность по дням недели</p>
            <div class="space-y-2">
              <div v-for="day in ownerForm.weeklySchedule" :key="day.dayOfWeek" class="grid grid-cols-[56px_1fr_1fr] items-center gap-3">
                <label class="inline-flex items-center gap-2 text-sm text-[#1d2d4a]">
                  <input v-model="day.isActive" type="checkbox" class="h-4 w-4" />
                  {{ dayLabels[day.dayOfWeek] }}
                </label>
                <UiInput
                  v-model.number="day.startHour"
                  type="number"
                  min="0"
                  max="23"
                  :disabled="!day.isActive"
                  :placeholder="`${dayLabels[day.dayOfWeek]} начало`"
                />
                <UiInput
                  v-model.number="day.endHour"
                  type="number"
                  min="1"
                  max="24"
                  :disabled="!day.isActive"
                  :placeholder="`${dayLabels[day.dayOfWeek]} конец`"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="mt-4 flex items-center gap-3">
          <UiButton :disabled="savingOwner" @click="updateOwner">{{ savingOwner ? 'Сохранение...' : 'Сохранить профиль' }}</UiButton>
          <span v-if="ownerSaved" class="text-sm text-emerald-600">Профиль сохранён</span>
        </div>
        <p v-if="ownerError" class="mt-2 text-sm text-red-600">{{ ownerError }}</p>
      </UiCard>

      <UiCard class="p-6">
        <h2 class="text-3xl font-semibold text-[#1d2d4a]">Создать тип события</h2>
        <div class="mt-4 grid gap-3 md:grid-cols-2">
          <UiInput v-model="form.name" placeholder="Название" class="md:col-span-2" />
          <UiInput v-model.number="form.durationMinutes" type="number" placeholder="Длительность, мин" />
          <div class="md:col-span-2">
            <UiTextarea v-model="form.description" placeholder="Описание" />
          </div>
        </div>
        <UiButton class="mt-4" @click="createEventType">Создать</UiButton>
        <p v-if="eventTypeError" class="mt-2 text-sm text-red-600">{{ eventTypeError }}</p>
      </UiCard>

      <UiCard class="p-6">
        <h2 class="text-3xl font-semibold text-[#1d2d4a]">Типы событий</h2>
        <ul class="mt-4 space-y-2 text-[#5f7190]">
          <li v-for="eventType in eventTypes" :key="eventType.id" class="rounded-xl border border-[#d2dbe7] px-4 py-3">
            <span class="font-semibold text-[#1d2d4a]">{{ eventType.name }}</span>
            <span class="ml-2 text-sm">({{ eventType.durationMinutes }} мин)</span>
            <p class="mt-1 text-sm">{{ eventType.description }}</p>
          </li>
        </ul>
      </UiCard>

      <UiCard class="p-6">
        <h2 class="text-3xl font-semibold text-[#1d2d4a]">Предстоящие встречи</h2>
        <UiTable class="mt-4">
          <thead>
            <tr class="border-b border-[#d2dbe7] text-left text-sm text-[#7388a4]">
              <th class="pb-3">Тип</th>
              <th class="pb-3">Гость</th>
              <th class="pb-3">Дата</th>
              <th class="pb-3">Время</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="booking in bookings"
              :key="booking.id"
              class="border-b border-[#edf2f8] text-[#1d2d4a]"
            >
              <td class="py-3">{{ booking.eventTypeId }}</td>
              <td class="py-3">{{ booking.guestName }} ({{ booking.guestEmail }})</td>
              <td class="py-3">{{ formatFullDate(booking.startAt) }}</td>
              <td class="py-3">{{ formatTimeRange(booking.startAt, booking.endAt) }}</td>
            </tr>
          </tbody>
        </UiTable>
      </UiCard>
    </template>
  </main>
</template>
