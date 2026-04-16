<script setup lang="ts">
import { useAttrs, computed } from 'vue'
import { cn } from '@/lib/utils'

const attrs = useAttrs()
const classes = computed(() =>
  cn(
    'rounded-2xl border border-[#cfdae7] bg-white shadow-[0_1px_0_rgba(20,40,80,0.04)]',
    attrs.class as string,
  ),
)
</script>

<template>
  <section :class="classes" v-bind="attrs">
    <slot />
  </section>
</template>
