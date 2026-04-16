<script setup lang="ts">
import { UiBadge } from '@/components/ui/badge'
import { UiCard } from '@/components/ui/card'
import type { EventType } from '@/types/api'

defineProps<{
  eventType: EventType
}>()
</script>

<template>
  <UiCard class="cursor-pointer p-5 transition hover:border-[#b7c8dd] hover:shadow-sm">
    <div class="mb-3 flex items-center justify-between gap-3">
      <h3 class="text-2xl font-semibold text-[#1d2d4a]">{{ eventType.name }}</h3>
      <UiBadge>{{ eventType.durationMinutes }} мин</UiBadge>
    </div>
    <p class="text-base text-[#7185a2]">{{ eventType.description }}</p>
  </UiCard>
</template>
