<script setup lang="ts">
import { CalendarClock } from 'lucide-vue-next'
</script>

<template>
  <header class="border-b border-[#d2dbe7] bg-white/90 backdrop-blur">
    <div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 md:px-8">
      <RouterLink to="/" class="inline-flex items-center gap-2 text-lg font-semibold text-[#1d2d4a]">
        <CalendarClock class="h-4 w-4 text-[#ff7a1a]" />
        <span>Calendar</span>
      </RouterLink>

      <nav class="flex items-center gap-2 text-sm">
        <RouterLink
          to="/name-owner"
          class="rounded-lg px-4 py-2 text-[#5f6f86] transition hover:bg-[#edf2f8] hover:text-[#1f2f49]"
        >
          Записаться
        </RouterLink>
        <RouterLink
          to="/admin"
          class="rounded-lg px-4 py-2 text-[#5f6f86] transition hover:bg-[#edf2f8] hover:text-[#1f2f49]"
        >
          Админка
        </RouterLink>
      </nav>
    </div>
  </header>
</template>
