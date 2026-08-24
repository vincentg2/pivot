<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const route = useRoute()
</script>

<template>
  <div class="app-shell">
    <AppHeader v-if="session.user && !route.meta.authLayout" />
    <RouterView />
  </div>
</template>
