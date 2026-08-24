<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/components/AuthLayout.vue'
import { ApiError } from '@/lib/api'
import { useSessionStore } from '@/stores/session'

const email = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)
const session = useSessionStore()
const router = useRouter()
const route = useRoute()
const { t } = useI18n()
async function submit() {
  submitting.value = true
  error.value = ''
  try {
    await session.login(email.value, password.value)
    await router.push(typeof route.query.next === 'string' ? route.query.next : '/')
  } catch (caught) {
    error.value = caught instanceof ApiError ? caught.message : 'Unable to sign in right now.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthLayout :eyebrow="t('auth.private')" :title="t('auth.story')" :intro="t('auth.intro')">
    <div class="form-heading">
      <p class="eyebrow">{{ t('auth.welcome') }}</p>
      <h2>{{ t('auth.signInTitle') }}</h2>
      <p>{{ t('auth.invitedOnly') }}</p>
    </div>
    <form class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label
        >{{ t('auth.email') }}<input v-model="email" type="email" autocomplete="email" required
      /></label>
      <label
        >{{ t('auth.password')
        }}<input v-model="password" type="password" autocomplete="current-password" required
      /></label>
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? t('auth.signingIn') : t('auth.signIn') }}
      </button>
    </form>
    <p class="form-foot">
      {{ t('auth.haveInvite') }}
      <RouterLink to="/register">{{ t('auth.createAccount') }}</RouterLink>
    </p>
  </AuthLayout>
</template>
