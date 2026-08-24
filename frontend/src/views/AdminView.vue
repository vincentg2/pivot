<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Copy, KeyRound, Newspaper, Plus, RefreshCw, Tv, X } from 'lucide-vue-next'
import { api, ApiError } from '@/lib/api'
import { useI18n } from 'vue-i18n'

const { locale, t } = useI18n()

interface Invitation {
  id: string
  label: string
  expiresAt: string | null
  maxUses: number
  uses: number
  createdAt: string
  revokedAt: string | null
}
const invitations = ref<Invitation[]>([])
const revealedCode = ref('')
const revealedResetLink = ref('')
const resetError = ref('')
const resetEmail = ref('')
const loading = ref(true)
const form = reactive({ label: '', maxUses: 1, expiresAt: '' })
interface CollectionRun {
  status: string
  recordsCount: number
  startedAt: string
  finishedAt: string | null
  errorMessage?: string
}
const collection = ref<{ enabled: boolean; latestRun: CollectionRun | null } | null>(null)
const sportCollection = ref<{ enabled: boolean; latestRun: CollectionRun | null } | null>(null)
const footaoCollection = ref<{ enabled: boolean; latestRun: CollectionRun | null } | null>(null)
const syncing = ref(false)
const syncingSport = ref(false)
const syncingFootao = ref(false)
const syncError = ref('')
async function load() {
  loading.value = true
  invitations.value = (await api<{ invitations: Invitation[] }>('/admin/invitations')).invitations
  loading.value = false
}
async function loadCollection() {
  collection.value = await api<{ enabled: boolean; latestRun: CollectionRun | null }>(
    '/admin/collections/football-data',
  )
}
async function loadSportCollection() {
  sportCollection.value = await api<{ enabled: boolean; latestRun: CollectionRun | null }>(
    '/admin/collections/football-data/sport',
  )
}
async function loadFootaoCollection() {
  footaoCollection.value = await api<{ enabled: boolean; latestRun: CollectionRun | null }>(
    '/admin/collections/footao',
  )
}
async function synchronize() {
  syncing.value = true
  syncError.value = ''
  try {
    await api('/admin/collections/football-data', { method: 'POST' })
    await loadCollection()
  } catch (error) {
    syncError.value = error instanceof Error ? error.message : t('admin.syncFailed')
  } finally {
    syncing.value = false
  }
}
async function synchronizeSport() {
  syncingSport.value = true
  syncError.value = ''
  try {
    await api('/admin/collections/football-data/sport', { method: 'POST' })
    await loadSportCollection()
  } catch (error) {
    syncError.value = error instanceof Error ? error.message : t('admin.syncFailed')
  } finally {
    syncingSport.value = false
  }
}
async function synchronizeFootao() {
  syncingFootao.value = true
  syncError.value = ''
  try {
    await api('/admin/collections/footao', { method: 'POST' })
    await loadFootaoCollection()
  } catch (error) {
    syncError.value = error instanceof Error ? error.message : t('admin.syncFailed')
  } finally {
    syncingFootao.value = false
  }
}
async function create() {
  const response = await api<{ code: string }>('/admin/invitations', {
    method: 'POST',
    body: JSON.stringify({
      label: form.label,
      maxUses: form.maxUses,
      expiresAt: form.expiresAt ? new Date(form.expiresAt).toISOString() : null,
    }),
  })
  revealedCode.value = response.code
  form.label = ''
  form.maxUses = 1
  form.expiresAt = ''
  await load()
}
async function revoke(id: string) {
  await api(`/admin/invitations/${id}`, { method: 'DELETE' })
  await load()
}
async function copyCode() {
  await navigator.clipboard.writeText(revealedCode.value)
}
async function createPasswordReset() {
  resetError.value = ''
  revealedResetLink.value = ''
  try {
    const response = await api<{ token: string; expiresAt: string }>('/admin/password-resets', {
      method: 'POST',
      body: JSON.stringify({ email: resetEmail.value }),
    })
    revealedResetLink.value = `${window.location.origin}/reset-password?token=${encodeURIComponent(response.token)}`
    resetEmail.value = ''
  } catch (caught) {
    resetError.value = caught instanceof ApiError ? caught.message : t('admin.resetFailed')
  }
}
async function copyResetLink() {
  await navigator.clipboard.writeText(revealedResetLink.value)
}
onMounted(() => {
  void Promise.all([load(), loadCollection(), loadSportCollection(), loadFootaoCollection()])
})
</script>

<template>
  <main class="admin page-width">
    <header class="settings-title">
      <p class="eyebrow">{{ t('admin.eyebrow') }}</p>
      <h1>{{ t('admin.invitations') }}</h1>
      <p>{{ t('admin.accessIntro') }}</p>
    </header>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">{{ t('admin.dataCollection') }}</p>
        <h2>football-data.org</h2>
        <p v-if="collection?.enabled" class="quiet">{{ t('admin.connectorEnabled') }}</p>
        <p v-else class="quiet">
          {{ t('admin.connectorDisabled') }}
        </p>
        <p v-if="collection?.latestRun" class="collection-status">
          <span class="status">{{ collection.latestRun.status }}</span>
          {{ collection.latestRun.recordsCount }} clubs ·
          {{ new Date(collection.latestRun.startedAt).toLocaleString(locale) }}
        </p>
        <p v-if="syncError" class="form-error" role="alert">{{ syncError }}</p>
      </div>
      <button
        class="button secondary"
        :disabled="!collection?.enabled || syncing"
        @click="synchronize"
      >
        <RefreshCw :size="17" :class="{ spinning: syncing }" />{{
          syncing ? t('admin.synchronizing') : t('admin.runNow')
        }}
      </button>
    </section>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">{{ t('admin.officialNews') }}</p>
        <h2>{{ t('admin.feeds') }}</h2>
        <p class="quiet">{{ t('admin.feedsIntro') }}</p>
        <RouterLink to="/admin/news" class="admin-detail-link"
          ><Newspaper :size="16" /> {{ t('admin.configureFeeds') }}</RouterLink
        >
      </div>
      <RouterLink to="/admin/news" class="button secondary">{{ t('admin.manageNews') }}</RouterLink>
    </section>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">{{ t('admin.sportCollection') }}</p>
        <h2>{{ t('admin.matchesTables') }}</h2>
        <p class="quiet">{{ t('admin.sportIntro') }}</p>
        <p v-if="sportCollection?.latestRun" class="collection-status">
          <span class="status">{{ sportCollection.latestRun.status }}</span>
          {{ sportCollection.latestRun.recordsCount }} {{ t('admin.records') }} ·
          {{ new Date(sportCollection.latestRun.startedAt).toLocaleString(locale) }}
        </p>
      </div>
      <button
        class="button secondary"
        :disabled="!sportCollection?.enabled || syncingSport"
        @click="synchronizeSport"
      >
        <RefreshCw :size="17" :class="{ spinning: syncingSport }" />{{
          syncingSport ? t('admin.synchronizing') : t('admin.runSport')
        }}
      </button>
    </section>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">{{ t('admin.tvCollection') }}</p>
        <h2>Footao</h2>
        <p v-if="footaoCollection?.enabled" class="quiet">
          {{ t('admin.footaoEnabled') }}
        </p>
        <p v-else class="quiet">
          {{ t('admin.footaoDisabled') }}
        </p>
        <p v-if="footaoCollection?.latestRun" class="collection-status">
          <span class="status">{{ footaoCollection.latestRun.status }}</span>
          {{ footaoCollection.latestRun.recordsCount }} {{ t('admin.listings') }} ·
          {{ new Date(footaoCollection.latestRun.startedAt).toLocaleString(locale) }}
        </p>
        <RouterLink to="/admin/tv" class="admin-detail-link"
          ><Tv :size="16" /> {{ t('admin.reviewTv') }}</RouterLink
        >
      </div>
      <button
        class="button secondary"
        :disabled="!footaoCollection?.enabled || syncingFootao"
        @click="synchronizeFootao"
      >
        <RefreshCw :size="17" :class="{ spinning: syncingFootao }" />{{
          syncingFootao ? t('admin.synchronizing') : t('admin.runTv')
        }}
      </button>
    </section>
    <section class="settings-card admin-create">
      <p class="eyebrow">{{ t('admin.recovery') }}</p>
      <h2>{{ t('admin.resetTitle') }}</h2>
      <p class="quiet">{{ t('admin.resetIntro') }}</p>
      <form class="inline-form password-reset-form" @submit.prevent="createPasswordReset">
        <label
          >{{ t('admin.memberEmail')
          }}<input
            v-model="resetEmail"
            type="email"
            autocomplete="off"
            placeholder="friend@example.com"
            required
        /></label>
        <button class="button secondary">
          <KeyRound :size="17" />{{ t('admin.generateReset') }}
        </button>
      </form>
      <p v-if="resetError" class="form-error" role="alert">{{ resetError }}</p>
      <div v-if="revealedResetLink" class="one-time-code reset-link" role="status">
        <div>
          <strong>{{ revealedResetLink }}</strong
          ><span>{{ t('admin.copyNowLink') }}</span>
        </div>
        <button class="icon-link" :aria-label="t('admin.copyReset')" @click="copyResetLink">
          <Copy :size="18" />
        </button>
      </div>
    </section>
    <section class="settings-card admin-create">
      <h2>{{ t('admin.createInvitation') }}</h2>
      <form class="inline-form" @submit.prevent="create">
        <label
          >{{ t('admin.label')
          }}<input
            v-model="form.label"
            maxlength="100"
            :placeholder="t('admin.labelPlaceholder')" /></label
        ><label
          >{{ t('admin.maxUses')
          }}<input v-model.number="form.maxUses" type="number" min="1" max="100" required /></label
        ><label
          >{{ t('admin.expires') }}<input v-model="form.expiresAt" type="datetime-local" /></label
        ><button class="button primary"><Plus :size="17" />{{ t('admin.create') }}</button>
      </form>
      <div v-if="revealedCode" class="one-time-code" role="status">
        <div>
          <strong>{{ revealedCode }}</strong
          ><span>{{ t('admin.copyNowCode') }}</span>
        </div>
        <button class="icon-link" :aria-label="t('admin.copyCode')" @click="copyCode">
          <Copy :size="18" />
        </button>
      </div>
    </section>
    <section class="settings-card">
      <h2>{{ t('admin.issued') }}</h2>
      <p v-if="loading">{{ t('common.loading') }}</p>
      <div v-else-if="!invitations.length" class="quiet">{{ t('admin.noInvitations') }}</div>
      <div v-else class="invite-list">
        <article v-for="invite in invitations" :key="invite.id">
          <div>
            <strong>{{ invite.label || t('admin.untitled') }}</strong
            ><span
              >{{ t('admin.used', { uses: invite.uses, max: invite.maxUses }) }} ·
              {{
                invite.expiresAt
                  ? t('admin.expiresOn', {
                      date: new Date(invite.expiresAt).toLocaleDateString(locale),
                    })
                  : t('admin.noExpiry')
              }}</span
            >
          </div>
          <span v-if="invite.revokedAt" class="status">{{ t('admin.revoked') }}</span
          ><button
            v-else
            class="icon-link"
            :aria-label="t('admin.revoke')"
            @click="revoke(invite.id)"
          >
            <X :size="18" />
          </button>
        </article>
      </div>
    </section>
  </main>
</template>
