<script setup lang="ts">
import { House, LogOut, Settings } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import PivotLogo from '@/components/PivotLogo.vue'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const router = useRouter()
async function logout() {
  await session.logout()
  await router.push('/login')
}
</script>

<template>
  <header class="site-header">
    <RouterLink to="/" class="brand" aria-label="Pivot home"
      ><PivotLogo class="brand-mark" /><span>Pivot</span></RouterLink
    >
    <nav aria-label="Account navigation">
      <RouterLink to="/" class="text-link home-link" aria-label="Home"
        ><House :size="16" /><span class="home-label">Home</span></RouterLink
      >
      <RouterLink to="/clubs" class="text-link">Clubs</RouterLink>
      <RouterLink to="/matches" class="text-link">Matches</RouterLink>
      <RouterLink to="/tv" class="text-link">TV</RouterLink>
      <RouterLink to="/standings" class="text-link">Tables</RouterLink>
      <RouterLink to="/news" class="text-link">News</RouterLink>
      <RouterLink v-if="session.user?.role === 'admin'" to="/admin" class="text-link"
        >Admin</RouterLink
      >
      <RouterLink to="/profile" class="icon-link"
        ><Settings :size="18" /><span class="sr-only">Profile settings</span></RouterLink
      >
      <button class="icon-link" type="button" @click="logout">
        <LogOut :size="18" /><span class="sr-only">Sign out</span>
      </button>
    </nav>
  </header>
</template>
