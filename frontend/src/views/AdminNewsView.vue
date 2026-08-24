<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowLeft, Plus, RefreshCw, Trash2 } from 'lucide-vue-next'
import { api } from '@/lib/api'
import type { NewsFeed } from '@/lib/news'
import type { Club } from '@/stores/favorites'

interface CollectionRun {
  status: string
  recordsCount: number
  startedAt: string
  finishedAt: string | null
  errorMessage?: string
}
const clubs = ref<Club[]>([])
const feeds = ref<NewsFeed[]>([])
const collection = ref<{ enabled: boolean; latestRun: CollectionRun | null } | null>(null)
const syncing = ref(false)
const error = ref('')
const form = reactive({ clubId: '', sourceName: '', url: '', enabled: true })
async function load() {
  const [clubResponse, feedResponse, statusResponse] = await Promise.all([
    api<{ clubs: Club[] }>('/clubs'),
    api<{ feeds: NewsFeed[] }>('/admin/news/feeds'),
    api<{ enabled: boolean; latestRun: CollectionRun | null }>('/admin/collections/news'),
  ])
  clubs.value = clubResponse.clubs
  feeds.value = feedResponse.feeds
  collection.value = statusResponse
}
async function create() {
  error.value = ''
  try {
    await api('/admin/news/feeds', { method: 'POST', body: JSON.stringify(form) })
    form.clubId = ''
    form.sourceName = ''
    form.url = ''
    form.enabled = true
    await load()
  } catch (caught) {
    error.value = caught instanceof Error ? caught.message : 'Feed could not be saved.'
  }
}
async function remove(id: string) {
  await api(`/admin/news/feeds/${id}`, { method: 'DELETE' })
  await load()
}
async function synchronize() {
  syncing.value = true
  error.value = ''
  try {
    await api('/admin/collections/news', { method: 'POST' })
    await load()
  } catch (caught) {
    error.value = caught instanceof Error ? caught.message : 'Synchronization failed.'
  } finally {
    syncing.value = false
  }
}
onMounted(load)
</script>

<template>
  <main class="admin-news page-width">
    <RouterLink to="/admin" class="back-link"><ArrowLeft :size="16" /> Administration</RouterLink>
    <header class="settings-title">
      <p class="eyebrow">Official news administration</p>
      <h1>Feeds, kept deliberate.</h1>
      <p>
        Configure only official club RSS or Atom feeds. Pivot stores metadata and links, never
        article bodies.
      </p>
    </header>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">Hourly collection</p>
        <h2>Official club feeds</h2>
        <p class="quiet">
          {{ feeds.filter((feed) => feed.enabled).length }} enabled feeds · 30-day retention.
        </p>
        <p v-if="collection?.latestRun" class="collection-status">
          <span class="status">{{ collection.latestRun.status }}</span>
          {{ collection.latestRun.recordsCount }} items ·
          {{ new Date(collection.latestRun.startedAt).toLocaleString() }}
        </p>
      </div>
      <button
        class="button secondary"
        :disabled="!collection?.enabled || syncing"
        @click="synchronize"
      >
        <RefreshCw :size="17" :class="{ spinning: syncing }" />{{
          syncing ? 'Synchronizing…' : 'Run news collection'
        }}
      </button>
    </section>
    <section class="settings-card">
      <h2>Add an official feed</h2>
      <form class="news-feed-form" @submit.prevent="create">
        <label
          >Club<select v-model="form.clubId" required>
            <option value="" disabled>Select a club</option>
            <option v-for="club in clubs" :key="club.id" :value="club.id">{{ club.name }}</option>
          </select></label
        ><label
          >Source name<input
            v-model="form.sourceName"
            maxlength="120"
            placeholder="Official club news"
            required /></label
        ><label class="feed-url"
          >RSS or Atom URL<input
            v-model="form.url"
            type="url"
            maxlength="2000"
            placeholder="https://club.example/feed.xml"
            required /></label
        ><label class="check-filter admin-hidden"
          ><input v-model="form.enabled" type="checkbox" /> Enabled</label
        ><button class="button primary"><Plus :size="16" /> Add feed</button>
      </form>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
    </section>
    <section class="settings-card">
      <h2>Configured feeds</h2>
      <p v-if="!feeds.length" class="quiet">No official feeds configured yet.</p>
      <div v-else class="feed-list">
        <article v-for="feed in feeds" :key="feed.id">
          <div>
            <strong>{{ feed.clubName }}</strong
            ><span>{{ feed.sourceName }} · {{ feed.enabled ? 'enabled' : 'disabled' }}</span
            ><a :href="feed.url" target="_blank" rel="noopener noreferrer">{{ feed.url }}</a>
          </div>
          <button
            class="icon-link"
            type="button"
            :aria-label="`Delete ${feed.sourceName} for ${feed.clubName}`"
            @click="remove(feed.id)"
          >
            <Trash2 :size="17" />
          </button>
        </article>
      </div>
    </section>
  </main>
</template>
