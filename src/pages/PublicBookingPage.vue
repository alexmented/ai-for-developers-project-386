<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppErrorState from '@/components/common/AppErrorState.vue'
import MiniCalendarGrid from '@/components/booking/MiniCalendarGrid.vue'
import SlotStatusList from '@/components/booking/SlotStatusList.vue'
import { UiBadge } from '@/components/ui/badge'
import { UiButton } from '@/components/ui/button'
import { UiCard } from '@/components/ui/card'
import { UiInput } from '@/components/ui/input'
import { buildBookingWindow, formatFullDate, formatTimeRange, toDateKey } from '@/lib/date'
import { api } from '@/services/api'
import { ApiError, type EventType, type PublicOwnerProfile, type Slot } from '@/types/api'

const route = useRoute()
const router = useRouter()

const ownerSlug = computed(() => String(route.params.ownerSlug ?? ''))
const eventTypeId = computed(() => String(route.params.eventTypeId ?? ''))

const owner = ref<PublicOwnerProfile | null>(null)
const eventType = ref<EventType | null>(null)
const slots = ref<Slot[]>([])
const loading = ref(true)
const saving = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)

const selectedDate = ref<string | null>(null)
const selectedSlotStartAt = ref<string | null>(null)
const showGuestForm = ref(false)

const guestForm = reactive({
  name: '',
  email: '',
  comment: '',
})

const selectedSlot = computed(() =>
  slots.value.find((slot) => slot.startAt === selectedSlotStartAt.value) ?? null,
)

const slotsForSelectedDate = computed(() => {
  if (!selectedDate.value) {
    return slots.value
  }

  const key = toDateKey(selectedDate.value)
  return slots.value.filter((slot) => toDateKey(slot.startAt) === key)
})

async function loadData() {
  loading.value = true
  error.value = null

  try {
    const [ownerResponse, eventTypes, slotResponse] = await Promise.all([
      api.getOwnerProfile(ownerSlug.value),
      api.getPublicEventTypes(ownerSlug.value),
      api.getAvailableSlots(ownerSlug.value, eventTypeId.value, buildBookingWindow().from, buildBookingWindow().to),
    ])

    owner.value = ownerResponse
    eventType.value = eventTypes.find((item) => item.id === eventTypeId.value) ?? null
    slots.value = slotResponse

    if (slotResponse[0]) {
      selectedDate.value = slotResponse[0].startAt
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить слоты'
  } finally {
    loading.value = false
  }
}

function handleSelectDate(value: string) {
  selectedDate.value = value
  selectedSlotStartAt.value = null
  success.value = null
}

function handleBack() {
  void router.push(`/${ownerSlug.value}`)
}

function handleContinue() {
  if (!selectedSlotStartAt.value) {
    error.value = 'Выберите время'
    return
  }
  error.value = null
  showGuestForm.value = true
}

function handleBackToSlots() {
  showGuestForm.value = false
  error.value = null
}

async function handleBook() {
  if (!guestForm.name.trim() || !guestForm.email.trim()) {
    error.value = 'Укажите имя и email'
    return
  }

  error.value = null
  saving.value = true

  try {
    await api.createBooking(ownerSlug.value, {
      eventTypeId: eventTypeId.value,
      slotStartAt: selectedSlotStartAt.value!,
      guestName: guestForm.name.trim(),
      guestEmail: guestForm.email.trim(),
      guestComment: guestForm.comment.trim() || undefined,
    })

    success.value = 'Бронирование создано'
    showGuestForm.value = false
    guestForm.name = ''
    guestForm.email = ''
    guestForm.comment = ''
    await loadData()
  } catch (err) {
    if (err instanceof ApiError && err.code === 'SLOT_CONFLICT') {
      error.value = 'Выбранный слот уже занят. Обновите выбор времени.'
      showGuestForm.value = false
      await loadData()
    } else {
      error.value = err instanceof Error ? err.message : 'Не удалось создать бронирование'
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <main class="mx-auto max-w-6xl px-4 py-8 md:py-12">
    <h1 class="mb-6 text-6xl font-semibold text-[#12213c]">{{ eventType?.name ?? 'Встреча' }}</h1>

    <section v-if="loading" class="text-[#6f84a1]">Загрузка...</section>
    <AppErrorState v-else-if="error && !eventType" :message="error" title="Не удалось загрузить слоты" />

    <section v-else class="grid gap-4 xl:grid-cols-[320px_1fr_330px]">
      <UiCard class="p-4 md:p-5">
        <div class="mb-4 flex items-center gap-3">
          <img
            v-if="owner?.photoUrl"
            :src="owner.photoUrl"
            alt="Фото владельца"
            class="h-12 w-12 rounded-full object-cover border border-[#d2dbe7]"
          />
          <div v-else class="grid h-12 w-12 place-items-center rounded-full bg-gradient-to-b from-[#ffcd9d] to-[#23a79a] text-xl">
            🧑
          </div>
          <div>
            <p class="text-2xl font-semibold text-[#1d2d4a]">{{ owner?.displayName ?? 'Владелец' }}</p>
            <p class="text-sm text-[#7288a5]">{{ owner?.timezone ?? '' }}</p>
          </div>
        </div>

        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-4xl font-semibold text-[#1d2d4a]">{{ eventType?.name }}</h2>
          <UiBadge>{{ eventType?.durationMinutes ?? 0 }} мин</UiBadge>
        </div>
        <p class="text-base text-[#7185a2]">{{ eventType?.description }}</p>

        <div class="mt-5 space-y-3">
          <div class="rounded-xl bg-[#eef3f9] px-3 py-3 text-sm text-[#7084a1]">
            <p class="mb-1">Выбранная дата</p>
            <p class="font-semibold text-[#1d2d4a]">{{ selectedDate ? formatFullDate(selectedDate) : 'Дата не выбрана' }}</p>
          </div>
          <div class="rounded-xl bg-[#eef3f9] px-3 py-3 text-sm text-[#7084a1]">
            <p class="mb-1">Выбранное время</p>
            <p class="font-semibold text-[#1d2d4a]">
              {{ selectedSlot ? formatTimeRange(selectedSlot.startAt, selectedSlot.endAt) : 'Время не выбрано' }}
            </p>
          </div>
        </div>

        <p v-if="error" class="mt-3 text-sm text-red-500">{{ error }}</p>
        <p v-if="success" class="mt-3 text-sm text-green-600">{{ success }}</p>
      </UiCard>

      <MiniCalendarGrid
        :slots="slots"
        :selected-date="selectedDate"
        @select="handleSelectDate"
      />

      <template v-if="!showGuestForm">
        <SlotStatusList
          :slots="slotsForSelectedDate"
          :selected-slot-start-at="selectedSlotStartAt"
          :pending="saving"
          @select="selectedSlotStartAt = $event"
          @back="handleBack"
          @continue="handleContinue"
        />
      </template>

      <UiCard v-else class="flex h-full min-h-[540px] flex-col p-4">
        <h2 class="text-3xl font-semibold text-[#1d2d4a]">Ваши данные</h2>
        <p class="mt-2 text-sm text-[#7084a1]">Укажите имя и email для бронирования</p>

        <div class="mt-5 space-y-3">
          <UiInput v-model="guestForm.name" placeholder="Ваше имя" />
          <UiInput v-model="guestForm.email" type="email" placeholder="Ваш email" />
          <UiInput v-model="guestForm.comment" placeholder="Комментарий (необязательно)" />
        </div>

        <p v-if="error" class="mt-3 text-sm text-red-500">{{ error }}</p>

        <div class="mt-auto pt-4 grid grid-cols-2 gap-3">
          <UiButton variant="outline" @click="handleBackToSlots">Назад</UiButton>
          <UiButton :disabled="saving" @click="handleBook">
            {{ saving ? 'Бронируем...' : 'Забронировать' }}
          </UiButton>
        </div>
      </UiCard>
    </section>
  </main>
</template>
