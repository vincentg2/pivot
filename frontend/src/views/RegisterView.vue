<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/components/AuthLayout.vue'
import { ApiError } from '@/lib/api'
import { useSessionStore } from '@/stores/session'

const form = reactive({ invitationCode: '', nickname: '', email: '', password: '' })
const error = ref('')
const submitting = ref(false)
const session = useSessionStore()
const router = useRouter()
const { t } = useI18n()
async function submit() {
  submitting.value = true
  error.value = ''
  try {
    await session.register(form)
    await router.push('/')
  } catch (caught) {
    error.value = caught instanceof ApiError ? caught.message : t('auth.createFailed')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <AuthLayout
    :eyebrow="t('auth.invitationEyebrow')"
    :title="t('auth.invitationTitle')"
    :intro="t('auth.invitationIntro')"
  >
    <div class="form-heading">
      <p class="eyebrow">{{ t('auth.join') }}</p>
      <h2>{{ t('auth.accountTitle') }}</h2>
      <p>{{ t('auth.invitationHelp') }}</p>
    </div>
    <form class="form-stack" @submit.prevent="submit">
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      <label
        >{{ t('auth.invitationCode')
        }}<input v-model="form.invitationCode" autocomplete="one-time-code" required
      /></label>
      <label
        >{{ t('auth.nickname')
        }}<input
          v-model="form.nickname"
          autocomplete="nickname"
          minlength="2"
          maxlength="40"
          required
      /></label>
      <label
        >{{ t('auth.email')
        }}<input v-model="form.email" type="email" autocomplete="email" required
      /></label>
      <label
        >{{ t('auth.password')
        }}<input
          v-model="form.password"
          type="password"
          autocomplete="new-password"
          minlength="12"
          required
        /><small>{{ t('auth.passwordHelp') }}</small></label
      >
      <button class="button primary" type="submit" :disabled="submitting">
        {{ submitting ? t('auth.creating') : t('auth.createAccount') }}
      </button>
    </form>
    <p class="form-foot">
      {{ t('auth.alreadyMember') }} <RouterLink to="/login">{{ t('auth.signIn') }}</RouterLink>
    </p>
  </AuthLayout>
</template>
