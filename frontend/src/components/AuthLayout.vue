<script setup lang="ts">
import PivotLogo from '@/components/PivotLogo.vue'
import LocaleSwitch from '@/components/LocaleSwitch.vue'
import { useI18n } from 'vue-i18n'

defineProps<{ eyebrow: string; title: string; intro: string }>()
const { t } = useI18n()
</script>

<template>
  <main class="auth-page">
    <section class="auth-story" aria-labelledby="auth-title">
      <div class="auth-brand-row">
        <RouterLink to="/" class="brand"
          ><PivotLogo class="brand-mark" /><span>Pivot</span></RouterLink
        >
        <LocaleSwitch />
      </div>
      <div>
        <p class="eyebrow">{{ eyebrow }}</p>
        <h1 id="auth-title">{{ title }}</h1>
        <p class="lede">{{ intro }}</p>
      </div>
      <p class="quiet">{{ t('auth.considered') }}</p>
    </section>
    <section class="auth-panel"><slot /></section>
  </main>
</template>
