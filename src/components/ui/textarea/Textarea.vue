<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from '@/lib/utils'

const modelValue = defineModel<string | null>()

const attrs = useAttrs()
const classes = computed(() =>
  cn(
    'flex min-h-24 w-full rounded-xl border border-[#d3dbe7] bg-white px-3 py-2 text-sm text-[#13213a] placeholder:text-[#8da0ba] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#d3dbe7]',
    attrs.class as string,
  ),
)

function onInput(event: Event) {
  modelValue.value = (event.target as HTMLTextAreaElement).value
}
</script>

<template>
  <textarea :class="classes" :value="modelValue ?? ''" v-bind="attrs" @input="onInput" />
</template>
