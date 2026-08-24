<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import AuthLayout from '@/components/AuthLayout.vue'
import { api, ApiError } from '@/lib/api'
import { useSessionStore } from '@/stores/session'

const route = useRoute()
const session = useSessionStore()
const token = computed(() => (typeof route.query.token === 'string' ? route.query.token : ''))
const form = reactive({ password: '', confirmation: '' })
const error = ref('')
const submitting = ref(false)
const complete = ref(false)

async function submit() {
  error.value = ''
  if (!token.value) {
    error.value = 'This reset link is incomplete.'
    return
  }
  if (form.password !== form.confirmation) {
    error.value = 'The passwords do not match.'
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
    error.value = caught instanceof ApiError ? caught.message : 'Unable to reset this password.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthLayout
    eyebrow="Account recovery"
    title="Choose a new private key to your football world."
    intro="Administrator-issued links expire quickly and can only be used once."
  >
    <div class="form-heading">
      <p class="eyebrow">Password reset</p>
      <h2>{{ complete ? 'Password updated' : 'Set a new password' }}</h2>
      <p v-if="complete">
        Your existing sessions have been closed. Sign in with your new password.
      </p>
      <p v-else>Use at least 12 characters.</p>
    </div>
    <div v-if="complete" class="form-stack">
      <RouterLink to="/login" class="button primary">Return to sign in</RouterLink>
    </div>
    <form v-else class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label
        >New password<input
          v-model="form.password"
          type="password"
          autocomplete="new-password"
          minlength="12"
          maxlength="128"
          required
      /></label>
      <label
        >Confirm password<input
          v-model="form.confirmation"
          type="password"
          autocomplete="new-password"
          minlength="12"
          maxlength="128"
          required
      /></label>
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? 'Updating…' : 'Update password' }}
      </button>
    </form>
  </AuthLayout>
</template>
