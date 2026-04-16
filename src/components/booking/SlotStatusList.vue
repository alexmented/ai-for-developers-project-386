<script setup lang="ts">
import { computed } from 'vue'
import { UiButton } from '@/components/ui/button'
import { UiCard } from '@/components/ui/card'
import { formatTimeRange } from '@/lib/date'
import type { Slot } from '@/types/api'

const props = defineProps<{
  slots: Slot[]
  selectedSlotStartAt: string | null
  pending?: boolean
}>()

const emit = defineEmits<{
  select: [value: string]
  back: []
  continue: []
}>()

const sortedSlots = computed(() =>
  [...props.slots].sort((a, b) => new Date(a.startAt).getTime() - new Date(b.startAt).getTime()),
)
</script>

<template>
  <UiCard class="flex h-full min-h-[540px] flex-col p-4">
    <h2 class="text-3xl font-semibold text-[#1d2d4a]">Статус слотов</h2>

    <div class="mt-4 max-h-[420px] flex-1 space-y-2 overflow-y-auto pr-1">
      <button
        v-for="slot in sortedSlots"
        :key="slot.startAt"
        type="button"
        class="flex w-full items-center justify-between rounded-xl border px-3 py-3 text-left"
        :class="[
          selectedSlotStartAt === slot.startAt ? 'border-[#f7b788] bg-[#fff0e5]' : 'border-[#d2dbe7] bg-white',
          slot.isAvailable ? 'text-[#1d2d4a]' : 'cursor-not-allowed text-[#8a99ad]'
        ]"
        :disabled="!slot.isAvailable"
        @click="emit('select', slot.startAt)"
      >
        <span class="font-medium">{{ formatTimeRange(slot.startAt, slot.endAt) }}</span>
        <span class="font-semibold">{{ slot.isAvailable ? 'Свободно' : 'Занято' }}</span>
      </button>
    </div>

    <div class="mt-4 grid grid-cols-2 gap-3">
      <UiButton variant="outline" @click="emit('back')">Назад</UiButton>
      <UiButton :disabled="!selectedSlotStartAt || pending" @click="emit('continue')">
        {{ pending ? 'Бронируем...' : 'Продолжить' }}
      </UiButton>
    </div>
  </UiCard>
</template>
