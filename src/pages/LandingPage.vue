<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'
import { UiButton } from '@/components/ui/button'
import { UiCard } from '@/components/ui/card'
</script>

<template>
  <main class="mx-auto grid max-w-6xl gap-8 px-4 py-10 md:grid-cols-[1fr_420px] md:py-12">
    <section class="pt-5">
      <span class="inline-flex rounded-full bg-[#edf2f8] px-4 py-1 text-xs font-bold uppercase tracking-[0.08em] text-[#7388a4]">
        Быстрая запись на звонок
      </span>
      <h1 class="mt-5 text-6xl font-semibold tracking-tight text-[#12213c] md:text-7xl">Calendar</h1>
      <p class="mt-4 max-w-lg text-2xl leading-relaxed text-[#7388a4]">
        Забронируйте встречу за минуту: выберите тип события и удобное время.
      </p>
      <RouterLink to="/name-owner" class="mt-8 inline-flex">
        <UiButton size="lg" class="gap-2">
          Записаться
          <ArrowRight class="h-4 w-4" />
        </UiButton>
      </RouterLink>
    </section>

    <UiCard class="h-fit p-6 md:mt-4">
      <h2 class="text-4xl font-semibold text-[#13213d]">Возможности</h2>
      <ul class="mt-4 list-disc space-y-3 pl-5 text-base leading-relaxed text-[#6f84a1]">
        <li>Выбор типа события и удобного времени для встречи.</li>
        <li>Быстрое бронирование с подтверждением и дополнительными заметками.</li>
        <li>Управление типами встреч и просмотр предстоящих записей в админке.</li>
      </ul>
    </UiCard>
  </main>
</template>
