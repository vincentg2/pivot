<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { applyLocale, type Locale } from '@/i18n'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const { locale, t } = useI18n()

async function select(value: Locale) {
  if (value === locale.value) return
  applyLocale(value)
  if (session.user) await session.updateProfile(session.user.nickname, session.user.theme, value)
}
</script>

<template>
  <div class="locale-switch" role="group" :aria-label="t('locale.label')">
    <button type="button" :class="{ active: locale === 'fr' }" @click="select('fr')">FR</button>
    <button type="button" :class="{ active: locale === 'en' }" @click="select('en')">EN</button>
  </div>
</template>
