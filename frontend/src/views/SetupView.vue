<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowRight, KeyRound, ShieldCheck } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import AuthLayout from '@/components/AuthLayout.vue'
import { api } from '@/lib/api'
import { markInstalled } from '@/lib/installation'
import { useSessionStore } from '@/stores/session'

const router = useRouter()
const session = useSessionStore()
const configured = ref(true)
const loading = ref(true)
const error = ref('')
const form = reactive({ token: '', nickname: '', email: '', password: '' })
onMounted(async () => {
  const status = await api<{ setupRequired: boolean; setupTokenConfigured: boolean }>(
    '/setup/status',
  )
  configured.value = status.setupTokenConfigured
  loading.value = false
  if (!status.setupRequired) await router.replace('/login')
})
async function install() {
  error.value = ''
  try {
    await api('/setup', { method: 'POST', body: JSON.stringify(form) })
    markInstalled()
    await session.login(form.email, form.password)
    await router.push('/')
  } catch (caught) {
    error.value = caught instanceof Error ? caught.message : 'Installation could not be completed.'
  }
}
</script>

<template>
  <AuthLayout
    eyebrow="First installation"
    title="Make Pivot yours."
    intro="Create the first administrator, then invite the people who belong here."
  >
    <div v-if="loading" class="quiet">Checking this installation…</div>
    <div v-else>
      <div class="form-heading">
        <div class="setup-icon"><ShieldCheck :size="24" /></div>
        <h2>First administrator</h2>
        <p>Use the one-time setup token configured on the server.</p>
      </div>
      <p v-if="!configured" class="form-error" role="alert">
        Set a secret <code>SETUP_TOKEN</code> of at least 20 characters on the server, then restart
        Pivot.
      </p>
      <form v-else class="form-stack" @submit.prevent="install">
        <label
          >Setup token
          <div class="input-with-icon">
            <KeyRound :size="17" /><input
              v-model="form.token"
              type="password"
              autocomplete="one-time-code"
              required
            /></div></label
        ><label
          >Nickname<input
            v-model="form.nickname"
            minlength="2"
            maxlength="40"
            autocomplete="nickname"
            required /></label
        ><label
          >Email<input
            v-model="form.email"
            type="email"
            maxlength="254"
            autocomplete="email"
            required /></label
        ><label
          >Password<input
            v-model="form.password"
            type="password"
            minlength="12"
            maxlength="128"
            autocomplete="new-password"
            required
          /><small>At least 12 characters.</small></label
        >
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <button class="button primary">Create administrator <ArrowRight :size="17" /></button>
      </form>
    </div>
  </AuthLayout>
</template>
