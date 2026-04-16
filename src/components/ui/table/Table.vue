<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from '@/lib/utils'

const attrs = useAttrs()
const classes = computed(() => cn('w-full caption-bottom text-sm', attrs.class as string))
</script>

<template>
  <div class="w-full overflow-auto">
    <table :class="classes" v-bind="attrs">
      <slot />
    </table>
  </div>
</template>
