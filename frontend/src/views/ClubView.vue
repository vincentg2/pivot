<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ArrowLeft, ExternalLink, Heart, MapPin } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import ClubMark from '@/components/ClubMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import NewsCard from '@/components/NewsCard.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import { useFavoritesStore, type Club } from '@/stores/favorites'
import type { NewsItem } from '@/lib/news'

const route = useRoute()
const { t } = useI18n()
const favorites = useFavoritesStore()
const club = ref<Club | null>(null)
const notice = ref('')
const matches = ref<FootballMatch[]>([])
const news = ref<NewsItem[]>([])

async function toggle() {
  if (!club.value) return
  try {
    await favorites.toggle(club.value)
  } catch (caught) {
    notice.value = caught instanceof Error ? caught.message : t('clubs.favoriteFailed')
  }
}

onMounted(async () => {
  if (!favorites.ready) await favorites.load()
  club.value = (await api<{ club: Club }>(`/clubs/${route.params.id}`)).club
  const from = localDate(addDays(new Date(), -1))
  const to = localDate(addDays(new Date(), 7))
  matches.value = (
    await api<{ matches: FootballMatch[] }>(
      `/matches?from=${from}&to=${to}&club=${route.params.id}`,
    )
  ).matches
  news.value = (await api<{ news: NewsItem[] }>(`/news?club=${route.params.id}&limit=6`)).news
})
</script>

<template>
  <main class="club-page page-width">
    <RouterLink to="/clubs" class="back-link"
      ><ArrowLeft :size="16" /> {{ t('clubs.allClubs') }}</RouterLink
    >
    <p v-if="!club" class="catalog-state">{{ t('clubs.loadingClub') }}</p>
    <template v-else>
      <header class="club-hero">
        <ClubMark :name="club.name" :tla="club.tla" :crest-url="club.crestUrl" size="lg" />
        <div>
          <p class="eyebrow">{{ club.tla || 'Club' }}</p>
          <h1>{{ club.name }}</h1>
        </div>
        <button
          class="button secondary"
          :class="{ 'favorite-active': favorites.has(club.id) }"
          @click="toggle"
        >
          <Heart :size="17" :fill="favorites.has(club.id) ? 'currentColor' : 'none'" />
          {{ favorites.has(club.id) ? t('clubs.following') : t('clubs.follow') }}
        </button>
      </header>
      <p v-if="notice" class="form-error" role="alert">{{ notice }}</p>
      <section class="club-details">
        <article>
          <p class="eyebrow">{{ t('clubs.home') }}</p>
          <h2><MapPin :size="20" /> {{ club.venue || t('clubs.venueMissing') }}</h2>
        </article>
        <article>
          <p class="eyebrow">{{ t('clubs.competitions') }}</p>
          <h2>
            {{ club.competitions?.map((item) => item.name).join(', ') || t('clubs.notListed') }}
          </h2>
        </article>
        <article v-if="club.websiteUrl">
          <p class="eyebrow">{{ t('clubs.official') }}</p>
          <a :href="club.websiteUrl" target="_blank" rel="noopener noreferrer"
            >{{ t('clubs.website') }} <ExternalLink :size="15"
          /></a>
        </article>
      </section>
      <section class="club-news">
        <div class="section-heading">
          <div>
            <p class="eyebrow">{{ t('clubs.dispatches') }}</p>
            <h2>{{ t('clubs.latestNews') }}</h2>
          </div>
          <RouterLink to="/news" class="text-link">{{ t('clubs.allNews') }}</RouterLink>
        </div>
        <div v-if="news.length" class="news-grid club-news-grid">
          <NewsCard v-for="item in news" :key="item.id" :item="item" />
        </div>
        <p v-else class="quiet dashboard-empty-line">
          {{ t('clubs.noNews') }}
        </p>
      </section>
      <section class="club-fixtures">
        <div class="section-heading">
          <div>
            <p class="eyebrow">{{ t('clubs.recentUpcoming') }}</p>
            <h2>{{ t('clubs.matches') }}</h2>
          </div>
          <RouterLink to="/matches" class="text-link">{{ t('clubs.allMatches') }}</RouterLink>
        </div>
        <div v-if="matches.length" class="compact-matches">
          <MatchRow
            v-for="(match, index) in matches"
            :key="`${match.utcDate}-${index}`"
            :match="match"
          />
        </div>
        <p v-else class="quiet dashboard-empty-line">{{ t('clubs.noMatches') }}</p>
      </section>
      <p class="attribution">{{ t('common.footballAttribution') }}</p>
    </template>
  </main>
</template>
