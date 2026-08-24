<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Newspaper } from 'lucide-vue-next'
import NewsCard from '@/components/NewsCard.vue'
import { api } from '@/lib/api'
import type { NewsItem } from '@/lib/news'

const items = ref<NewsItem[]>([])
const loading = ref(true)
onMounted(async () => {
  items.value = (await api<{ news: NewsItem[] }>('/news?limit=60')).news
  loading.value = false
})
</script>

<template>
  <main class="news-page page-width">
    <header class="matches-heading news-heading">
      <p class="eyebrow">Official club sources</p>
      <h1>The morning papers.</h1>
      <p class="lede">
        Headlines from the clubs you follow, with every story opening at its original source.
      </p>
    </header>
    <p v-if="loading" class="catalog-state">Loading official news…</p>
    <section v-else-if="items.length" class="news-grid" aria-label="Official club news">
      <NewsCard v-for="item in items" :key="item.id" :item="item" show-club />
    </section>
    <section v-else class="empty-hero matches-empty">
      <Newspaper :size="28" />
      <h2>No official headlines yet</h2>
      <p>Follow clubs and ask an administrator to configure their official RSS or Atom feeds.</p>
    </section>
    <p class="attribution">
      Pivot stores titles, source names, dates and links for 30 days. Article content remains with
      each club.
    </p>
  </main>
</template>
