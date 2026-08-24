<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/components/AuthLayout.vue'
import { api, ApiError } from '@/lib/api'
import { useSessionStore } from '@/stores/session'

const route = useRoute()
const session = useSessionStore()
const { t } = useI18n()
const token = computed(() => (typeof route.query.token === 'string' ? route.query.token : ''))
const form = reactive({ password: '', confirmation: '' })
const error = ref('')
const submitting = ref(false)
const complete = ref(false)

async function submit() {
  error.value = ''
  if (!token.value) {
    error.value = t('auth.resetIncomplete')
    return
  }
  if (form.password !== form.confirmation) {
    error.value = t('auth.mismatch')
    return
  }
  submitting.value = true
  try {
    await api('/auth/password-reset', {
      method: 'POST',
      body: JSON.stringify({ token: token.value, password: form.password }),
    })
    await session.restore()
    complete.value = true
  } catch (caught) {
    error.value = caught instanceof ApiError ? caught.message : t('auth.resetFailed')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthLayout
    :eyebrow="t('auth.resetEyebrow')"
    :title="t('auth.resetTitle')"
    :intro="t('auth.resetIntro')"
  >
    <div class="form-heading">
      <p class="eyebrow">{{ t('auth.passwordReset') }}</p>
      <h2>{{ complete ? t('auth.passwordUpdated') : t('auth.newPasswordTitle') }}</h2>
      <p v-if="complete">
        {{ t('auth.sessionsClosed') }}
      </p>
      <p v-else>{{ t('auth.passwordHelp') }}</p>
    </div>
    <div v-if="complete" class="form-stack">
      <RouterLink to="/login" class="button primary">{{ t('auth.backToLogin') }}</RouterLink>
    </div>
    <form v-else class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label
        >{{ t('auth.newPassword')
        }}<input
          v-model="form.password"
          type="password"
          autocomplete="new-password"
          minlength="12"
          maxlength="128"
          required
      /></label>
      <label
        >{{ t('auth.confirmPassword')
        }}<input
          v-model="form.confirmation"
          type="password"
          autocomplete="new-password"
          minlength="12"
          maxlength="128"
          required
      /></label>
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? t('auth.updating') : t('auth.updatePassword') }}
      </button>
    </form>
  </AuthLayout>
</template>
