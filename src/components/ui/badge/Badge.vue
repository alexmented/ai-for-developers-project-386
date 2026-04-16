<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from '@/lib/utils'

const attrs = useAttrs()
const classes = computed(() =>
  cn('inline-flex items-center justify-center rounded-full bg-[#eaf0f7] px-2.5 py-1 text-xs font-semibold text-[#7288a5] min-w-[70px]', attrs.class as string),
)
</script>

<template>
  <span :class="classes" v-bind="attrs">
    <slot />
  </span>
</template>
