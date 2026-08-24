<script setup lang="ts">
import { House, LogOut, Settings } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LocaleSwitch from '@/components/LocaleSwitch.vue'
import PivotLogo from '@/components/PivotLogo.vue'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const router = useRouter()
const { t } = useI18n()
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
      <RouterLink to="/" class="text-link home-link" :aria-label="t('nav.home')"
        ><House :size="16" /><span class="home-label">{{ t('nav.home') }}</span></RouterLink
      >
      <RouterLink to="/clubs" class="text-link">{{ t('nav.clubs') }}</RouterLink>
      <RouterLink to="/matches" class="text-link">{{ t('nav.matches') }}</RouterLink>
      <RouterLink to="/tv" class="text-link">{{ t('nav.tv') }}</RouterLink>
      <RouterLink to="/standings" class="text-link">{{ t('nav.tables') }}</RouterLink>
      <RouterLink to="/news" class="text-link">{{ t('nav.news') }}</RouterLink>
      <RouterLink v-if="session.user?.role === 'admin'" to="/admin" class="text-link">{{
        t('nav.admin')
      }}</RouterLink>
      <LocaleSwitch />
      <RouterLink to="/profile" class="icon-link"
        ><Settings :size="18" /><span class="sr-only">{{ t('nav.profile') }}</span></RouterLink
      >
      <button class="icon-link" type="button" @click="logout">
        <LogOut :size="18" /><span class="sr-only">{{ t('nav.logout') }}</span>
      </button>
    </nav>
  </header>
</template>
