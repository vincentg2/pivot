<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Copy, Plus, X } from 'lucide-vue-next'
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
async function load() {
  loading.value = true
  invitations.value = (await api<{ invitations: Invitation[] }>('/admin/invitations')).invitations
  loading.value = false
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
onMounted(load)
</script>

<template>
  <main class="admin page-width">
    <header class="settings-title">
      <p class="eyebrow">Administration</p>
      <h1>Invitations</h1>
      <p>Control access with expiring, limited-use codes.</p>
    </header>
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
