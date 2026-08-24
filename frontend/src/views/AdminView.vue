<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Copy, Plus, RefreshCw, X } from 'lucide-vue-next'
import { api } from '@/lib/api'

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
const syncing = ref(false)
const syncingSport = ref(false)
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
async function synchronize() {
  syncing.value = true
  syncError.value = ''
  try {
    await api('/admin/collections/football-data', { method: 'POST' })
    await loadCollection()
  } catch (error) {
    syncError.value = error instanceof Error ? error.message : 'Synchronization failed.'
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
    syncError.value = error instanceof Error ? error.message : 'Synchronization failed.'
  } finally {
    syncingSport.value = false
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
onMounted(() => {
  void Promise.all([load(), loadCollection(), loadSportCollection()])
})
</script>

<template>
  <main class="admin page-width">
    <header class="settings-title">
      <p class="eyebrow">Administration</p>
      <h1>Invitations</h1>
      <p>Control access with expiring, limited-use codes.</p>
    </header>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">Data collection</p>
        <h2>football-data.org</h2>
        <p v-if="collection?.enabled" class="quiet">Connector enabled for this installation.</p>
        <p v-else class="quiet">
          Connector disabled. Add <code>FOOTBALL_DATA_API_KEY</code> to the local environment to
          enable it.
        </p>
        <p v-if="collection?.latestRun" class="collection-status">
          <span class="status">{{ collection.latestRun.status }}</span>
          {{ collection.latestRun.recordsCount }} clubs ·
          {{ new Date(collection.latestRun.startedAt).toLocaleString() }}
        </p>
        <p v-if="syncError" class="form-error" role="alert">{{ syncError }}</p>
      </div>
      <button
        class="button secondary"
        :disabled="!collection?.enabled || syncing"
        @click="synchronize"
      >
        <RefreshCw :size="17" :class="{ spinning: syncing }" />{{
          syncing ? 'Synchronizing…' : 'Run now'
        }}
      </button>
    </section>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">Sports collection</p>
        <h2>Matches & standings</h2>
        <p class="quiet">
          Refreshes yesterday, today and the next seven days for the five major leagues.
        </p>
        <p v-if="sportCollection?.latestRun" class="collection-status">
          <span class="status">{{ sportCollection.latestRun.status }}</span>
          {{ sportCollection.latestRun.recordsCount }} records ·
          {{ new Date(sportCollection.latestRun.startedAt).toLocaleString() }}
        </p>
      </div>
      <button
        class="button secondary"
        :disabled="!sportCollection?.enabled || syncingSport"
        @click="synchronizeSport"
      >
        <RefreshCw :size="17" :class="{ spinning: syncingSport }" />{{
          syncingSport ? 'Synchronizing…' : 'Run sports data'
        }}
      </button>
    </section>
    <section class="settings-card admin-create">
      <h2>Create an invitation</h2>
      <form class="inline-form" @submit.prevent="create">
        <label
          >Label<input v-model="form.label" maxlength="100" placeholder="Friends, August" /></label
        ><label
          >Maximum uses<input
            v-model.number="form.maxUses"
            type="number"
            min="1"
            max="100"
            required /></label
        ><label>Expires<input v-model="form.expiresAt" type="datetime-local" /></label
        ><button class="button primary"><Plus :size="17" />Create</button>
      </form>
      <div v-if="revealedCode" class="one-time-code" role="status">
        <div>
          <strong>{{ revealedCode }}</strong
          ><span>Copy it now. It will not be shown again.</span>
        </div>
        <button class="icon-link" aria-label="Copy invitation code" @click="copyCode">
          <Copy :size="18" />
        </button>
      </div>
    </section>
    <section class="settings-card">
      <h2>Issued invitations</h2>
      <p v-if="loading">Loading…</p>
      <div v-else-if="!invitations.length" class="quiet">No invitations yet.</div>
      <div v-else class="invite-list">
        <article v-for="invite in invitations" :key="invite.id">
          <div>
            <strong>{{ invite.label || 'Untitled invitation' }}</strong
            ><span
              >{{ invite.uses }} of {{ invite.maxUses }} used ·
              {{
                invite.expiresAt
                  ? `expires ${new Date(invite.expiresAt).toLocaleDateString()}`
                  : 'no expiry'
              }}</span
            >
          </div>
          <span v-if="invite.revokedAt" class="status">Revoked</span
          ><button
            v-else
            class="icon-link"
            aria-label="Revoke invitation"
            @click="revoke(invite.id)"
          >
            <X :size="18" />
          </button>
        </article>
      </div>
    </section>
  </main>
</template>
