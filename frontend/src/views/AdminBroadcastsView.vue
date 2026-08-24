<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowLeft, EyeOff, RotateCcw, Save } from 'lucide-vue-next'
import { api } from '@/lib/api'
import { addDays, localDate } from '@/lib/football'
import { listingTime, type BroadcastAudit, type BroadcastListing } from '@/lib/broadcast'
import { useI18n } from 'vue-i18n'
const { locale, t } = useI18n()

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
    message.value = t('adminTv.saved')
  } finally {
    saving.value = false
  }
}
async function clearCorrection() {
  if (!selected.value) return
  await api(`/admin/broadcasts/${selected.value.id}/correction`, { method: 'DELETE' })
  message.value = t('adminTv.cleared')
  selected.value = null
  await load()
}
onMounted(load)
</script>

<template>
  <main class="admin-tv page-width">
    <RouterLink to="/admin" class="back-link"
      ><ArrowLeft :size="16" /> {{ t('admin.eyebrow') }}</RouterLink
    >
    <header class="settings-title">
      <p class="eyebrow">{{ t('adminTv.eyebrow') }}</p>
      <h1>{{ t('adminTv.title') }}</h1>
      <p>{{ t('adminTv.intro') }}</p>
    </header>
    <div class="admin-tv-layout">
      <section class="settings-card tv-review-list">
        <h2>{{ t('adminTv.next7') }}</h2>
        <p v-if="!listings.length" class="quiet">{{ t('adminTv.none') }}</p>
        <button
          v-for="item in listings"
          :key="item.id"
          type="button"
          :class="{ active: selected?.id === item.id }"
          @click="edit(item)"
        >
          <time :datetime="item.startsAt"
            >{{ new Date(item.startsAt).toLocaleDateString(locale) }} ·
            {{ listingTime(item.startsAt) }}</time
          >
          <strong>{{ item.label }}</strong
          ><span>{{ item.channels.join(', ') }}</span>
          <small v-if="item.corrected">{{ t('adminTv.corrected') }}</small
          ><small v-if="item.hidden"><EyeOff :size="12" /> {{ t('adminTv.hidden') }}</small>
        </button>
      </section>
      <section class="settings-card tv-editor">
        <template v-if="selected">
          <h2>{{ t('adminTv.edit') }}</h2>
          <form class="form-stack" @submit.prevent="save">
            <label
              >{{ t('adminTv.date') }}<input v-model="form.startsAt" type="datetime-local" required
            /></label>
            <label
              >{{ t('adminTv.label') }}<input v-model="form.label" maxlength="280" required
            /></label>
            <div class="two-fields">
              <label>{{ t('adminTv.home') }}<input v-model="form.homeName" maxlength="140" /></label
              ><label
                >{{ t('adminTv.away') }}<input v-model="form.awayName" maxlength="140"
              /></label>
            </div>
            <label
              >{{ t('adminTv.competition') }}<input v-model="form.competitionName" maxlength="140"
            /></label>
            <div class="two-fields">
              <label
                >{{ t('adminTv.type')
                }}<select v-model="form.kind">
                  <option value="live">{{ t('adminTv.live') }}</option>
                  <option value="delayed">{{ t('adminTv.delayed') }}</option>
                  <option value="replay">{{ t('adminTv.replay') }}</option>
                </select></label
              ><label
                >{{ t('adminTv.channels')
                }}<input v-model="form.channels" placeholder="Canal+, beIN Sports 1"
              /></label>
            </div>
            <label
              >{{ t('adminTv.note')
              }}<input
                v-model="form.note"
                maxlength="300"
                :placeholder="t('adminTv.notePlaceholder')"
            /></label>
            <label class="check-filter admin-hidden"
              ><input v-model="form.hidden" type="checkbox" /> {{ t('adminTv.hide') }}</label
            >
            <div class="save-row">
              <button class="button primary" :disabled="saving">
                <Save :size="16" /> {{ saving ? t('adminTv.saving') : t('adminTv.save') }}</button
              ><button
                v-if="selected.corrected"
                type="button"
                class="button secondary"
                @click="clearCorrection"
              >
                <RotateCcw :size="16" /> {{ t('adminTv.restore') }}
              </button>
            </div>
            <p v-if="message" class="collection-status" role="status">{{ message }}</p>
          </form>
        </template>
        <div v-else class="tv-editor-empty">
          <p>{{ t('adminTv.select') }}</p>
        </div>
      </section>
    </div>
    <section class="settings-card tv-audit">
      <h2>{{ t('adminTv.audit') }}</h2>
      <p v-if="!audit.length" class="quiet">{{ t('adminTv.noAudit') }}</p>
      <ol v-else>
        <li v-for="entry in audit" :key="entry.id">
          <span class="status">{{ entry.action }}</span
          ><code>{{ entry.listingId.slice(0, 8) }}</code
          ><time :datetime="entry.createdAt">{{
            new Date(entry.createdAt).toLocaleString(locale)
          }}</time>
        </li>
      </ol>
    </section>
  </main>
</template>
