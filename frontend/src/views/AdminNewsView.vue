<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ArrowLeft, Plus, RefreshCw, Trash2 } from 'lucide-vue-next'
import { api } from '@/lib/api'
import { useI18n } from 'vue-i18n'
import type { NewsFeed } from '@/lib/news'
import type { Club } from '@/stores/favorites'
const { locale, t } = useI18n()

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
    error.value = caught instanceof Error ? caught.message : t('adminNews.saveFailed')
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
    error.value = caught instanceof Error ? caught.message : t('admin.syncFailed')
  } finally {
    syncing.value = false
  }
}
onMounted(load)
</script>

<template>
  <main class="admin-news page-width">
    <RouterLink to="/admin" class="back-link"
      ><ArrowLeft :size="16" /> {{ t('admin.eyebrow') }}</RouterLink
    >
    <header class="settings-title">
      <p class="eyebrow">{{ t('adminNews.eyebrow') }}</p>
      <h1>{{ t('adminNews.title') }}</h1>
      <p>{{ t('adminNews.intro') }}</p>
    </header>
    <section class="settings-card collection-card">
      <div>
        <p class="eyebrow">{{ t('adminNews.hourly') }}</p>
        <h2>{{ t('adminNews.feeds') }}</h2>
        <p class="quiet">
          {{ t('adminNews.enabledFeeds', { count: feeds.filter((feed) => feed.enabled).length }) }}
        </p>
        <p v-if="collection?.latestRun" class="collection-status">
          <span class="status">{{ collection.latestRun.status }}</span>
          {{ collection.latestRun.recordsCount }} {{ t('adminNews.items') }} ·
          {{ new Date(collection.latestRun.startedAt).toLocaleString(locale) }}
        </p>
      </div>
      <button
        class="button secondary"
        :disabled="!collection?.enabled || syncing"
        @click="synchronize"
      >
        <RefreshCw :size="17" :class="{ spinning: syncing }" />{{
          syncing ? t('admin.synchronizing') : t('adminNews.run')
        }}
      </button>
    </section>
    <section class="settings-card">
      <h2>{{ t('adminNews.add') }}</h2>
      <form class="news-feed-form" @submit.prevent="create">
        <label
          >{{ t('adminNews.club')
          }}<select v-model="form.clubId" required>
            <option value="" disabled>{{ t('adminNews.selectClub') }}</option>
            <option v-for="club in clubs" :key="club.id" :value="club.id">{{ club.name }}</option>
          </select></label
        ><label
          >{{ t('adminNews.source')
          }}<input
            v-model="form.sourceName"
            maxlength="120"
            placeholder="Official club news"
            required /></label
        ><label class="feed-url"
          >{{ t('adminNews.url')
          }}<input
            v-model="form.url"
            type="url"
            maxlength="2000"
            placeholder="https://club.example/feed.xml"
            required /></label
        ><label class="check-filter admin-hidden"
          ><input v-model="form.enabled" type="checkbox" /> {{ t('adminNews.enabled') }}</label
        ><button class="button primary"><Plus :size="16" /> {{ t('adminNews.addFeed') }}</button>
      </form>
      <p v-if="error" class="form-error" role="alert">{{ error }}</p>
    </section>
    <section class="settings-card">
      <h2>{{ t('adminNews.configured') }}</h2>
      <p v-if="!feeds.length" class="quiet">{{ t('adminNews.none') }}</p>
      <div v-else class="feed-list">
        <article v-for="feed in feeds" :key="feed.id">
          <div>
            <strong>{{ feed.clubName }}</strong
            ><span
              >{{ feed.sourceName }} ·
              {{ feed.enabled ? t('adminNews.enabled') : t('adminNews.disabled') }}</span
            ><a :href="feed.url" target="_blank" rel="noopener noreferrer">{{ feed.url }}</a>
          </div>
          <button
            class="icon-link"
            type="button"
            :aria-label="t('adminNews.delete', { source: feed.sourceName, club: feed.clubName })"
            @click="remove(feed.id)"
          >
            <Trash2 :size="17" />
          </button>
        </article>
      </div>
    </section>
  </main>
</template>
