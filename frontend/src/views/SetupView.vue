<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowRight, KeyRound, ShieldCheck } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/components/AuthLayout.vue'
import { api } from '@/lib/api'
import { markInstalled } from '@/lib/installation'
import { useSessionStore } from '@/stores/session'

const router = useRouter()
const session = useSessionStore()
const { t } = useI18n()
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
    error.value = caught instanceof Error ? caught.message : t('auth.installFailed')
  }
}
</script>

<template>
  <AuthLayout
    :eyebrow="t('auth.setupEyebrow')"
    :title="t('auth.setupTitle')"
    :intro="t('auth.setupIntro')"
  >
    <div v-if="loading" class="quiet">{{ t('auth.checking') }}</div>
    <div v-else>
      <div class="form-heading">
        <div class="setup-icon"><ShieldCheck :size="24" /></div>
        <h2>{{ t('auth.firstAdmin') }}</h2>
        <p>{{ t('auth.setupHelp') }}</p>
      </div>
      <p v-if="!configured" class="form-error" role="alert">
        {{ t('auth.setupMissing') }}
      </p>
      <form v-else class="form-stack" @submit.prevent="install">
        <label
          >{{ t('auth.setupToken') }}
          <div class="input-with-icon">
            <KeyRound :size="17" /><input
              v-model="form.token"
              type="password"
              autocomplete="one-time-code"
              required
            /></div></label
        ><label
          >{{ t('auth.nickname')
          }}<input
            v-model="form.nickname"
            minlength="2"
            maxlength="40"
            autocomplete="nickname"
            required /></label
        ><label
          >{{ t('auth.email')
          }}<input
            v-model="form.email"
            type="email"
            maxlength="254"
            autocomplete="email"
            required /></label
        ><label
          >{{ t('auth.password')
          }}<input
            v-model="form.password"
            type="password"
            minlength="12"
            maxlength="128"
            autocomplete="new-password"
            required
          /><small>{{ t('auth.passwordHelp') }}</small></label
        >
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
        <button class="button primary">
          {{ t('auth.createAdmin') }} <ArrowRight :size="17" />
        </button>
      </form>
    </div>
  </AuthLayout>
</template>
