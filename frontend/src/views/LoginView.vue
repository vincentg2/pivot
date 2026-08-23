<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
  <AuthLayout
    eyebrow="Private by design"
    title="Your football world, in one quiet place."
    intro="Follow the clubs that matter. See what’s next without the noise."
  >
    <div class="form-heading">
      <p class="eyebrow">Welcome back</p>
      <h2>Sign in to Pivot</h2>
      <p>Access is limited to invited members.</p>
    </div>
    <form class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label>Email<input v-model="email" type="email" autocomplete="email" required /></label>
      <label
        >Password<input v-model="password" type="password" autocomplete="current-password" required
      /></label>
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? 'Signing in…' : 'Sign in' }}
      </button>
    </form>
    <p class="form-foot">
      Have an invitation? <RouterLink to="/register">Create an account</RouterLink>
    </p>
  </AuthLayout>
</template>
