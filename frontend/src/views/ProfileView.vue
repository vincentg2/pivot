<script setup lang="ts">
import { ref } from 'vue'
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
import { useSessionStore } from '@/stores/session'

const session = useSessionStore()
const router = useRouter()
const nickname = ref(session.user?.nickname ?? '')
const theme = ref<Theme>(session.user?.theme ?? 'system')
const saved = ref(false)
const busy = ref(false)
async function save() {
  busy.value = true
  await session.updateProfile(nickname.value, theme.value)
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
      <p class="eyebrow">Account</p>
      <h1>Your profile</h1>
      <p>Identity and appearance follow you across devices.</p>
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
        <label>Nickname<input v-model="nickname" minlength="2" maxlength="40" required /></label>
        <fieldset>
          <legend>Appearance</legend>
          <div class="theme-options">
            <label v-for="option in ['system', 'light', 'dark'] as Theme[]" :key="option"
              ><input v-model="theme" type="radio" :value="option" />{{ option }}</label
            >
          </div>
        </fieldset>
        <div class="save-row">
          <button class="button primary" :disabled="busy">Save changes</button
          ><span v-if="saved" role="status">Saved.</span>
        </div>
      </form>
    </section>
    <section class="danger-card">
      <div>
        <h2>Delete account</h2>
        <p>Permanently removes your profile, sessions, and future favorites.</p>
      </div>
      <DialogRoot
        ><DialogTrigger class="button danger">Delete account</DialogTrigger
        ><DialogPortal
          ><DialogOverlay class="dialog-overlay" /><DialogContent class="dialog-content"
            ><DialogTitle>Delete your Pivot account?</DialogTitle
            ><DialogDescription>This action is permanent and cannot be undone.</DialogDescription>
            <div class="dialog-actions">
              <DialogClose class="button secondary">Keep account</DialogClose
              ><button class="button danger" @click="remove">Delete permanently</button>
            </div></DialogContent
          ></DialogPortal
        ></DialogRoot
      >
    </section>
  </main>
</template>
