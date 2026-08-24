<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Newspaper } from 'lucide-vue-next'
import NewsCard from '@/components/NewsCard.vue'
import { api } from '@/lib/api'
import type { NewsItem } from '@/lib/news'

const items = ref<NewsItem[]>([])
const loading = ref(true)
const { t } = useI18n()
onMounted(async () => {
  items.value = (await api<{ news: NewsItem[] }>('/news?limit=60')).news
  loading.value = false
})
</script>

<template>
  <main class="news-page page-width">
    <header class="matches-heading news-heading">
      <h1>{{ t('news.title') }}</h1>
    </header>
    <p v-if="loading" class="catalog-state">{{ t('news.loading') }}</p>
    <section v-else-if="items.length" class="news-grid" aria-label="Official club news">
      <NewsCard v-for="item in items" :key="item.id" :item="item" show-club />
    </section>
    <section v-else class="empty-hero matches-empty">
      <Newspaper :size="28" />
      <h2>{{ t('news.empty') }}</h2>
      <p>{{ t('news.emptyHelp') }}</p>
    </section>
    <p class="attribution">
      Pivot stores titles, source names, dates and links for 30 days. Article content remains with
      each club.
    </p>
  </main>
</template>
