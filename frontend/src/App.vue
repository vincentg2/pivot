<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
</script>

<template>
  <div class="app-shell">
    <AppHeader v-if="session.user" />
    <RouterView />
  </div>
</template>
