<script setup lang="ts">
import { computed, ref } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { UiCard } from '@/components/ui/card'
import { getMonthLabel, toDateKey } from '@/lib/date'
import type { Slot } from '@/types/api'

const props = defineProps<{
  slots: Slot[]
  selectedDate: string | null
}>()

const emit = defineEmits<{
  select: [value: string]
}>()

const monthOffset = ref(0)

const activeMonth = computed(() => {
  const date = new Date()
  date.setDate(1)
  date.setHours(12, 0, 0, 0)
  date.setMonth(date.getMonth() + monthOffset.value)
  return date
})

const monthLabel = computed(() => getMonthLabel(activeMonth.value.toISOString()))

const dayHeaders = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const days = computed(() => {
  const month = activeMonth.value
  const year = month.getFullYear()
  const monthIndex = month.getMonth()
  const firstDay = new Date(year, monthIndex, 1)
  const startShift = (firstDay.getDay() + 6) % 7
  const daysInMonth = new Date(year, monthIndex + 1, 0).getDate()
  const toIso = (day: number, targetMonth = monthIndex) => new Date(Date.UTC(year, targetMonth, day, 12)).toISOString()

  const result: Array<{ iso: string; label: number; hasSlot: boolean }> = []
  for (let i = 0; i < startShift; i += 1) {
    const date = new Date(year, monthIndex, 1 - (startShift - i))
    const iso = toIso(date.getDate(), date.getMonth())
    result.push({ iso, label: date.getDate(), hasSlot: false })
  }

  for (let day = 1; day <= daysInMonth; day += 1) {
    const iso = toIso(day)
    const key = toDateKey(iso)
    const hasSlot = props.slots.some((slot) => toDateKey(slot.startAt) === key && slot.isAvailable)
    result.push({ iso, label: day, hasSlot })
  }

  let trailingDay = 1
  while (result.length % 7 !== 0) {
    const iso = toIso(trailingDay, monthIndex + 1)
    result.push({ iso, label: trailingDay, hasSlot: false })
    trailingDay += 1
  }

  return result
})

function goPrev() {
  monthOffset.value -= 1
}

function goNext() {
  monthOffset.value += 1
}
</script>

<template>
  <UiCard class="p-4 md:p-5">
    <div class="mb-5 flex items-center justify-between">
      <h2 class="text-3xl font-semibold text-[#1d2d4a]">Календарь</h2>
      <div class="flex gap-2">
        <button
          type="button"
          class="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-[#d2dbe7] text-[#627692]"
          @click="goPrev"
        >
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button
          type="button"
          class="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-[#d2dbe7] text-[#627692]"
          @click="goNext"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>

    <p class="mb-4 text-sm font-medium text-[#667b98]">{{ monthLabel }}</p>

    <div class="mb-2 grid grid-cols-7 gap-1 text-center text-xs font-semibold text-[#697f9d]">
      <span v-for="day in dayHeaders" :key="day">{{ day }}</span>
    </div>

    <div class="grid grid-cols-7 gap-1.5">
      <button
        v-for="day in days"
        :key="day.iso"
        type="button"
        class="flex h-12 items-center justify-center rounded-xl border text-sm font-medium"
        :class="[
          selectedDate && toDateKey(selectedDate) === toDateKey(day.iso)
            ? 'border-[#f7b788] bg-[#fff0e5] text-[#1d2d4a]'
            : 'border-[#d8e0eb] bg-[#f8fbff] text-[#6b7f9c]',
          day.hasSlot ? 'hover:border-[#b8c8dd] hover:bg-white' : 'opacity-60',
        ]"
        @click="emit('select', day.iso)"
      >
        {{ day.label }}
      </button>
    </div>
  </UiCard>
</template>
