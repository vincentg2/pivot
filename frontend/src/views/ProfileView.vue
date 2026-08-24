<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from 'reka-ui'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import type { Theme } from '@/stores/session'
import type { Locale } from '@/i18n'
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const router = useRouter()
const nickname = ref(session.user?.nickname ?? '')
const theme = ref<Theme>(session.user?.theme ?? 'system')
const localePreference = ref<Locale>(session.user?.locale ?? 'fr')
const { t } = useI18n()
const saved = ref(false)
const busy = ref(false)
async function save() {
  busy.value = true
  await session.updateProfile(nickname.value, theme.value, localePreference.value)
  saved.value = true
  busy.value = false
  setTimeout(() => {
    saved.value = false
  }, 2500)
}
async function remove() {
  busy.value = true
  await session.deleteAccount()
  await router.push('/login')
}
</script>

<template>
  <main class="settings page-width">
    <header class="settings-title">
      <p class="eyebrow">{{ t('profile.eyebrow') }}</p>
      <h1>{{ t('profile.title') }}</h1>
      <p>{{ t('profile.intro') }}</p>
    </header>
    <section v-if="session.user" class="settings-card">
      <div class="profile-summary">
        <AvatarMonogram :name="session.user.nickname" :seed="session.user.avatarSeed" size="lg" />
        <div>
          <strong>{{ session.user.nickname }}</strong
          ><span>{{ session.user.email }}</span
          ><small>{{ session.user.role }}</small>
        </div>
      </div>
      <form class="form-stack" @submit.prevent="save">
        <label
          >{{ t('profile.nickname')
          }}<input v-model="nickname" minlength="2" maxlength="40" required
        /></label>
        <fieldset>
          <legend>{{ t('profile.appearance') }}</legend>
          <div class="theme-options">
            <label v-for="option in ['system', 'light', 'dark'] as Theme[]" :key="option"
              ><input v-model="theme" type="radio" :value="option" />{{
                t(`profile.${option}`)
              }}</label
            >
          </div>
        </fieldset>
        <label
          >{{ t('profile.language')
          }}<select v-model="localePreference">
            <option value="fr">Français</option>
            <option value="en">English</option>
          </select></label
        >
        <div class="save-row">
          <button class="button primary" :disabled="busy">{{ t('common.save') }}</button
          ><span v-if="saved" role="status">{{ t('common.saved') }}</span>
        </div>
      </form>
    </section>
    <section class="danger-card">
      <div>
        <h2>{{ t('profile.delete') }}</h2>
        <p>{{ t('profile.deleteHelp') }}</p>
      </div>
      <DialogRoot
        ><DialogTrigger class="button danger">{{ t('profile.delete') }}</DialogTrigger
        ><DialogPortal
          ><DialogOverlay class="dialog-overlay" /><DialogContent class="dialog-content"
            ><DialogTitle>{{ t('profile.confirmTitle') }}</DialogTitle
            ><DialogDescription>{{ t('profile.confirmHelp') }}</DialogDescription>
            <div class="dialog-actions">
              <DialogClose class="button secondary">{{ t('profile.keep') }}</DialogClose
              ><button class="button danger" @click="remove">
                {{ t('profile.deleteForever') }}
              </button>
            </div></DialogContent
          ></DialogPortal
        ></DialogRoot
      >
    </section>
  </main>
</template>
