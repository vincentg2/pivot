<script setup lang="ts">
import { ArrowUpRight } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import type { NewsItem } from '@/lib/news'

defineProps<{ item: NewsItem; showClub?: boolean }>()
const { locale, t } = useI18n()
</script>

<template>
  <article class="news-card">
    <div class="news-meta">
      <span v-if="showClub">{{ item.clubTla || item.clubName }}</span>
      <span>{{ item.sourceName }}</span>
      <time :datetime="item.publishedAt">{{
        new Date(item.publishedAt).toLocaleDateString(locale)
      }}</time>
    </div>
    <h3>
      <a :href="item.linkUrl" target="_blank" rel="noopener noreferrer">{{ item.title }}</a>
    </h3>
    <a
      :href="item.linkUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="news-open"
      :aria-label="t('news.readOn', { title: item.title, source: item.sourceName })"
      >{{ t('news.readSource') }} <ArrowUpRight :size="15"
    /></a>
  </article>
</template>
