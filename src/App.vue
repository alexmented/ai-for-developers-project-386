<template>
  <div class="min-h-screen bg-[radial-gradient(circle_at_top_left,#d8e4f7_0%,#edf3fa_45%,#f7efe7_100%)]">
    <AppHeader />
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/layout/AppHeader.vue'
</script>
