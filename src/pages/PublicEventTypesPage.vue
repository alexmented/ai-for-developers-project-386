<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppErrorState from '@/components/common/AppErrorState.vue'
import { UiCard } from '@/components/ui/card'
import EventTypeCard from '@/components/booking/EventTypeCard.vue'
import { api } from '@/services/api'
import type { EventType, PublicOwnerProfile } from '@/types/api'

const route = useRoute()
const router = useRouter()
const ownerSlug = computed(() => String(route.params.ownerSlug ?? ''))

const owner = ref<PublicOwnerProfile | null>(null)
const eventTypes = ref<EventType[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const dayLabels: Record<string, string> = {
  monday: 'Пн',
  tuesday: 'Вт',
  wednesday: 'Ср',
  thursday: 'Чт',
  friday: 'Пт',
  saturday: 'Сб',
  sunday: 'Вс',
}

const weeklyScheduleSummary = computed(() => {
  const schedule = owner.value?.weeklySchedule ?? []
  const active = schedule.filter((item) => item.isActive)

  if (active.length === 0) {
    return 'Нет активных дней'
  }

  return active
    .map((item) => `${dayLabels[item.dayOfWeek] ?? item.dayOfWeek} ${String(item.startHour).padStart(2, '0')}:00–${String(item.endHour).padStart(2, '0')}:00`)
    .join(', ')
})

async function loadData() {
  loading.value = true
  error.value = null

  try {
    const [ownerResponse, eventTypesResponse] = await Promise.all([
      api.getOwnerProfile(ownerSlug.value),
      api.getPublicEventTypes(ownerSlug.value),
    ])

    owner.value = ownerResponse
    eventTypes.value = eventTypesResponse
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить страницу'
  } finally {
    loading.value = false
  }
}

function openEvent(eventTypeId: string) {
  void router.push(`/${ownerSlug.value}/${eventTypeId}`)
}

onMounted(loadData)
</script>

<template>
  <main class="mx-auto max-w-5xl px-4 py-8 md:py-12">
    <UiCard class="p-6 md:p-7">
      <div class="flex items-center gap-4">
        <img
          v-if="owner?.photoUrl"
          :src="owner.photoUrl"
          alt="Фото владельца"
          class="h-14 w-14 rounded-full object-cover border border-[#d2dbe7]"
        />
        <div v-else class="grid h-14 w-14 place-items-center rounded-full bg-gradient-to-b from-[#ffcd9d] to-[#23a79a] text-xl">
          🧑
        </div>
        <div>
          <h2 class="text-3xl font-semibold text-[#14233f]">{{ owner?.displayName ?? 'Tota' }}</h2>
          <p class="text-sm text-[#7084a1]">
            {{ owner?.timezone ?? 'Europe/Moscow' }} · {{ weeklyScheduleSummary }}
          </p>
        </div>
      </div>

      <h1 class="mt-6 text-5xl font-semibold text-[#12213c]">Выберите тип события</h1>
      <p class="mt-2 text-base text-[#6f84a1]">Нажмите на карточку, чтобы открыть календарь и выбрать удобный слот.</p>
    </UiCard>

    <section v-if="loading" class="mt-6 text-[#6f84a1]">Загрузка...</section>
    <AppErrorState v-else-if="error" class="mt-6" :message="error" title="Не удалось загрузить страницу" />

    <section v-else class="mt-6 grid gap-4 md:grid-cols-2">
      <button
        v-for="eventType in eventTypes"
        :key="eventType.id"
        type="button"
        class="text-left"
        @click="openEvent(eventType.id)"
      >
        <EventTypeCard :event-type="eventType" />
      </button>
    </section>
  </main>
</template>
