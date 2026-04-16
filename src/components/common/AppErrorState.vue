<script setup lang="ts">
import { computed } from 'vue'
import { UiCard } from '@/components/ui/card'

const props = withDefaults(
  defineProps<{
    title?: string
    message?: string | null
    compact?: boolean
  }>(),
  {
    title: 'Что-то пошло не так',
    message: null,
    compact: false,
  },
)

const normalizedMessage = computed(() => {
  const raw = (props.message ?? '').trim()
  if (!raw) {
    return 'Попробуйте обновить страницу чуть позже.'
  }

  const lower = raw.toLowerCase()
  if (lower.includes('failed to fetch') || lower.includes('networkerror') || lower.includes('load failed')) {
    return 'Не удалось связаться с сервером. Проверьте интернет-соединение и попробуйте снова.'
  }

  return raw
})
</script>

<template>
  <UiCard :class="compact ? 'p-4' : 'p-6'">
    <div class="flex items-start gap-3">
      <div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-red-50 text-red-600">⚠️</div>
      <div>
        <h3 class="text-lg font-semibold text-[#1d2d4a]">{{ title }}</h3>
        <p class="mt-1 text-sm text-[#6f84a1]">{{ normalizedMessage }}</p>
      </div>
    </div>
  </UiCard>
</template>
