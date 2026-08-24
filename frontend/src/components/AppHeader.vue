<script setup lang="ts">
import { LogOut, Settings } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
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
      ><span class="brand-mark">P</span><span>Pivot</span></RouterLink
    >
    <nav aria-label="Account navigation">
      <RouterLink to="/clubs" class="text-link">Clubs</RouterLink>
      <RouterLink to="/matches" class="text-link">Matches</RouterLink>
      <RouterLink to="/standings" class="text-link">Tables</RouterLink>
      <RouterLink to="/tv" class="text-link">TV</RouterLink>
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
