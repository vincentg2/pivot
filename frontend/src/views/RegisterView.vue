<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '@/components/AuthLayout.vue'
import { ApiError } from '@/lib/api'
import { useSessionStore } from '@/stores/session'

const form = reactive({ invitationCode: '', nickname: '', email: '', password: '' })
const error = ref('')
const submitting = ref(false)
const session = useSessionStore()
const router = useRouter()
async function submit() {
  submitting.value = true
  error.value = ''
  try {
    await session.register(form)
    await router.push('/')
  } catch (caught) {
    error.value = caught instanceof ApiError ? caught.message : 'Unable to create the account.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthLayout
    eyebrow="By invitation"
    title="A better view of the week in football."
    intro="Fixtures, results and the stories around your clubs — shaped around you."
  >
    <div class="form-heading">
      <p class="eyebrow">Join Pivot</p>
      <h2>Create your account</h2>
      <p>Your invitation may have an expiry or usage limit.</p>
    </div>
    <form class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label
        >Invitation code<input v-model="form.invitationCode" autocomplete="one-time-code" required
      /></label>
      <label
        >Nickname<input
          v-model="form.nickname"
          autocomplete="nickname"
          minlength="2"
          maxlength="40"
          required
      /></label>
      <label>Email<input v-model="form.email" type="email" autocomplete="email" required /></label>
      <label
        >Password<input
          v-model="form.password"
          type="password"
          autocomplete="new-password"
          minlength="12"
          required
        /><small>At least 12 characters.</small></label
      >
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? 'Creating…' : 'Create account' }}
      </button>
    </form>
    <p class="form-foot">Already a member? <RouterLink to="/login">Sign in</RouterLink></p>
  </AuthLayout>
</template>
