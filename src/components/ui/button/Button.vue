<script setup lang="ts">
import { cva } from 'class-variance-authority'
import { computed, useAttrs } from 'vue'
import { cn } from '@/lib/utils'

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap rounded-[12px] text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'bg-[#ff7a1a] text-white hover:bg-[#f36e0e]',
        secondary: 'bg-[#eaf0f7] text-[#55657c] hover:bg-[#dde6f0]',
        ghost: 'bg-transparent text-[#50617c] hover:bg-[#eaf0f7]',
        outline: 'border border-[#d2dbe7] bg-white text-[#1d2d4a] hover:bg-[#f5f7fb]',
      },
      size: {
        default: 'h-11 px-6',
        sm: 'h-9 rounded-[10px] px-4 text-xs',
        lg: 'h-12 rounded-[14px] px-8 text-base',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

const props = withDefaults(
  defineProps<{
    variant?: 'default' | 'secondary' | 'ghost' | 'outline'
    size?: 'default' | 'sm' | 'lg'
    as?: 'button' | 'a'
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    as: 'button',
    type: 'button',
  },
)

const attrs = useAttrs()

const classes = computed(() => cn(buttonVariants({ variant: props.variant, size: props.size }), attrs.class as string))
</script>

<template>
  <a v-if="as === 'a'" :class="classes" v-bind="attrs">
    <slot />
  </a>
  <button v-else :class="classes" :type="type" v-bind="attrs">
    <slot />
  </button>
</template>
