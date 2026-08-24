<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowLeft, EyeOff, RotateCcw, Save } from 'lucide-vue-next'
import { api } from '@/lib/api'
import { addDays, localDate } from '@/lib/football'
import { listingTime, type BroadcastAudit, type BroadcastListing } from '@/lib/broadcast'

const listings = ref<BroadcastListing[]>([])
const audit = ref<BroadcastAudit[]>([])
const selected = ref<BroadcastListing | null>(null)
const saving = ref(false)
const message = ref('')
const form = reactive({
  startsAt: '',
  label: '',
  homeName: '',
  awayName: '',
  competitionName: '',
  kind: 'live',
  channels: '',
  hidden: false,
  note: '',
})

async function load() {
  const from = localDate(new Date())
  const to = localDate(addDays(new Date(), 7))
  const [listingResponse, auditResponse] = await Promise.all([
    api<{ listings: BroadcastListing[] }>(`/admin/broadcasts?${new URLSearchParams({ from, to })}`),
    api<{ audit: BroadcastAudit[] }>('/admin/broadcasts/audit'),
  ])
  listings.value = listingResponse.listings
  audit.value = auditResponse.audit
}
function edit(item: BroadcastListing) {
  selected.value = item
  const date = new Date(item.startsAt)
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  form.startsAt = local.toISOString().slice(0, 16)
  form.label = item.label
  form.homeName = item.homeName
  form.awayName = item.awayName
  form.competitionName = item.competitionName
  form.kind = item.kind
  form.channels = item.channels.join(', ')
  form.hidden = Boolean(item.hidden)
  form.note = ''
  message.value = ''
}
async function save() {
  if (!selected.value) return
  saving.value = true
  message.value = ''
  try {
    await api(`/admin/broadcasts/${selected.value.id}/correction`, {
      method: 'PUT',
      body: JSON.stringify({
        ...form,
        startsAt: new Date(form.startsAt).toISOString(),
        channels: form.channels
          .split(',')
          .map((value) => value.trim())
          .filter(Boolean),
      }),
    })
    await load()
    const refreshed = listings.value.find((item) => item.id === selected.value?.id)
    if (refreshed) edit(refreshed)
    message.value = 'Correction saved and recorded in the audit trail.'
  } finally {
    saving.value = false
  }
}
async function clearCorrection() {
  if (!selected.value) return
  await api(`/admin/broadcasts/${selected.value.id}/correction`, { method: 'DELETE' })
  message.value = 'Correction cleared. The provider value is active again.'
  selected.value = null
  await load()
}
onMounted(load)
</script>

<template>
  <main class="admin-tv page-width">
    <RouterLink to="/admin" class="back-link"><ArrowLeft :size="16" /> Administration</RouterLink>
    <header class="settings-title">
      <p class="eyebrow">Television administration</p>
      <h1>Review the schedule.</h1>
      <p>
        Operator corrections always take priority over imported values. Every change is auditable.
      </p>
    </header>
    <div class="admin-tv-layout">
      <section class="settings-card tv-review-list">
        <h2>Next seven days</h2>
        <p v-if="!listings.length" class="quiet">No TV listings have been collected.</p>
        <button
          v-for="item in listings"
          :key="item.id"
          type="button"
          :class="{ active: selected?.id === item.id }"
          @click="edit(item)"
        >
          <time :datetime="item.startsAt"
            >{{ new Date(item.startsAt).toLocaleDateString() }} ·
            {{ listingTime(item.startsAt) }}</time
          >
          <strong>{{ item.label }}</strong
          ><span>{{ item.channels.join(', ') }}</span> <small v-if="item.corrected">Corrected</small
          ><small v-if="item.hidden"><EyeOff :size="12" /> Hidden</small>
        </button>
      </section>
      <section class="settings-card tv-editor">
        <template v-if="selected">
          <h2>Edit listing</h2>
          <form class="form-stack" @submit.prevent="save">
            <label
              >Date and time<input v-model="form.startsAt" type="datetime-local" required
            /></label>
            <label>Display label<input v-model="form.label" maxlength="280" required /></label>
            <div class="two-fields">
              <label>Home team<input v-model="form.homeName" maxlength="140" /></label
              ><label>Away team<input v-model="form.awayName" maxlength="140" /></label>
            </div>
            <label>Competition<input v-model="form.competitionName" maxlength="140" /></label>
            <div class="two-fields">
              <label
                >Type<select v-model="form.kind">
                  <option value="live">Live</option>
                  <option value="delayed">Delayed</option>
                  <option value="replay">Replay</option>
                </select></label
              ><label
                >Channels<input v-model="form.channels" placeholder="Canal+, beIN Sports 1"
              /></label>
            </div>
            <label
              >Audit note<input
                v-model="form.note"
                maxlength="300"
                placeholder="Reason for the correction"
            /></label>
            <label class="check-filter admin-hidden"
              ><input v-model="form.hidden" type="checkbox" /> Hide this listing from members</label
            >
            <div class="save-row">
              <button class="button primary" :disabled="saving">
                <Save :size="16" /> {{ saving ? 'Saving…' : 'Save correction' }}</button
              ><button
                v-if="selected.corrected"
                type="button"
                class="button secondary"
                @click="clearCorrection"
              >
                <RotateCcw :size="16" /> Restore import
              </button>
            </div>
            <p v-if="message" class="collection-status" role="status">{{ message }}</p>
          </form>
        </template>
        <div v-else class="tv-editor-empty">
          <p>Select a listing to correct its public presentation.</p>
        </div>
      </section>
    </div>
    <section class="settings-card tv-audit">
      <h2>Recent audit</h2>
      <p v-if="!audit.length" class="quiet">No manual corrections yet.</p>
      <ol v-else>
        <li v-for="entry in audit" :key="entry.id">
          <span class="status">{{ entry.action }}</span
          ><code>{{ entry.listingId.slice(0, 8) }}</code
          ><time :datetime="entry.createdAt">{{ new Date(entry.createdAt).toLocaleString() }}</time>
        </li>
      </ol>
    </section>
  </main>
</template>
