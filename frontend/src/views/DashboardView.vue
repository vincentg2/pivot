<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowRight, CalendarDays, ChevronDown, Newspaper, Radio, Trophy } from 'lucide-vue-next'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import ClubMark from '@/components/ClubMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import NewsCard from '@/components/NewsCard.vue'
import RecentResult from '@/components/RecentResult.vue'
import { api } from '@/lib/api'
import { broadcastChannelsByMatch, type BroadcastListing } from '@/lib/broadcast'
import { favoriteMatchesWithin, latestFavoriteResults, matchesForClub } from '@/lib/dashboard'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import { useSessionStore } from '@/stores/session'
import { useFavoritesStore } from '@/stores/favorites'
import type { NewsItem } from '@/lib/news'
const session = useSessionStore()
const { locale, t } = useI18n()
const favorites = useFavoritesStore()
const loadedMatches = ref<FootballMatch[]>([])
const recentMatches = ref<FootballMatch[]>([])
const broadcasts = ref<BroadcastListing[]>([])
const news = ref<NewsItem[]>([])
const selectedClubId = ref<string | null>(null)
const visibleMatchCount = ref(8)
const matchBatchSize = 8
const channelsByMatch = computed(() => broadcastChannelsByMatch(broadcasts.value))
const matchesInHorizon = computed(() => favoriteMatchesWithin(loadedMatches.value, new Date(), 30))
const filteredMatches = computed(() => matchesForClub(matchesInHorizon.value, selectedClubId.value))
const favoriteMatches = computed(() => filteredMatches.value.slice(0, visibleMatchCount.value))
const recentResults = computed(() => latestFavoriteResults(recentMatches.value))
const hasMoreMatches = computed(() => visibleMatchCount.value < filteredMatches.value.length)
const todayLabel = computed(() =>
  new Intl.DateTimeFormat(locale.value, { weekday: 'long', day: 'numeric', month: 'long' }).format(
    new Date(),
  ),
)
onMounted(async () => {
  if (!favorites.ready) await favorites.load()
  const from = localDate(new Date())
  const to = localDate(addDays(new Date(), 30))
  const recentFrom = localDate(addDays(new Date(), -30))
  const recentTo = localDate(addDays(new Date(), -1))
  const [matchResponse, recentResponse, broadcastResponse, newsResponse] = await Promise.all([
    api<{ matches: FootballMatch[] }>(`/matches?from=${from}&to=${to}`),
    api<{ matches: FootballMatch[] }>(`/matches?from=${recentFrom}&to=${recentTo}`),
    api<{ listings: BroadcastListing[] }>(`/broadcasts?from=${from}&to=${to}`),
    api<{ news: NewsItem[] }>('/news?limit=6'),
  ])
  loadedMatches.value = matchResponse.matches
  recentMatches.value = recentResponse.matches
  broadcasts.value = broadcastResponse.listings
  news.value = newsResponse.news
})
watch(selectedClubId, () => (visibleMatchCount.value = matchBatchSize))
</script>

<template>
  <main class="dashboard page-width">
    <section class="welcome-row">
      <div>
        <h1>{{ t('dashboard.greeting', { name: session.user?.nickname }) }}</h1>
        <p class="dashboard-date">{{ todayLabel }}</p>
      </div>
      <AvatarMonogram
        v-if="session.user"
        :name="session.user.nickname"
        :seed="session.user.avatarSeed"
        size="lg"
      />
      <div v-if="recentResults.length" class="recent-results">
        <p>{{ t('dashboard.recent') }}</p>
        <div class="recent-results-list">
          <RecentResult v-for="match in recentResults" :key="match.id" :match="match" />
        </div>
      </div>
    </section>
    <section v-if="news.length" class="dashboard-news">
      <div class="section-heading">
        <h2>{{ t('dashboard.latestNews') }}</h2>
        <RouterLink to="/news" class="text-link">{{ t('dashboard.allNews') }}</RouterLink>
      </div>
      <div class="news-grid dashboard-news-grid">
        <NewsCard v-for="item in news" :key="item.id" :item="item" show-club />
      </div>
    </section>
    <section v-if="favorites.clubs.length" class="dashboard-matches">
      <div class="section-heading">
        <h2>
          {{ t('dashboard.upcoming') }} <span>{{ t('dashboard.days') }}</span>
        </h2>
        <RouterLink to="/matches" class="text-link">{{ t('dashboard.allMatches') }}</RouterLink>
      </div>
      <div class="favorite-club-filters" role="group" :aria-label="t('dashboard.filterFavorites')">
        <button
          type="button"
          :class="{ active: selectedClubId === null }"
          :aria-pressed="selectedClubId === null"
          @click="selectedClubId = null"
        >
          {{ t('dashboard.allClubs') }}
        </button>
        <button
          v-for="club in favorites.clubs"
          :key="club.id"
          type="button"
          :class="{ active: selectedClubId === club.id }"
          :aria-pressed="selectedClubId === club.id"
          @click="selectedClubId = club.id"
        >
          <ClubMark :name="club.name" :tla="club.tla" :crest-url="club.crestUrl" />
          {{ club.shortName || club.name }}
        </button>
      </div>
      <div v-if="favoriteMatches.length" class="compact-matches">
        <MatchRow
          v-for="(match, index) in favoriteMatches"
          :key="match.id || `${match.utcDate}-${index}`"
          :match="match"
          :channels="channelsByMatch.get(match.id)"
          show-channel-status
          show-date
        />
      </div>
      <p v-else class="quiet dashboard-empty-line">
        {{ t('dashboard.noMatches') }}
      </p>
      <div v-if="hasMoreMatches" class="dashboard-match-more">
        <button type="button" class="button secondary" @click="visibleMatchCount += matchBatchSize">
          {{ t('dashboard.showNext') }} <ChevronDown :size="17" />
        </button>
        <span aria-live="polite"
          >{{ favoriteMatches.length }} {{ t('common.of') }} {{ filteredMatches.length }}</span
        >
      </div>
    </section>
    <section v-if="favorites.ready && favorites.clubs.length" class="favorite-dashboard">
      <div class="section-heading">
        <h2>{{ t('dashboard.favoriteClubs') }}</h2>
        <RouterLink to="/clubs" class="text-link">{{ t('dashboard.editFavorites') }}</RouterLink>
      </div>
      <div class="favorite-strip">
        <RouterLink v-for="club in favorites.clubs" :key="club.id" :to="`/clubs/${club.id}`">
          <ClubMark :name="club.name" :tla="club.tla" :crest-url="club.crestUrl" />
          <span>{{ club.shortName || club.name }}</span
          ><ArrowRight :size="15" />
        </RouterLink>
      </div>
    </section>
    <section v-else class="empty-hero">
      <div class="empty-icon"><Trophy :size="26" /></div>
      <p class="eyebrow">{{ t('dashboard.startingEleven') }}</p>
      <h2>{{ t('dashboard.chooseClubs') }}</h2>
      <p>{{ t('dashboard.chooseIntro') }}</p>
      <RouterLink to="/clubs" class="button primary"
        >{{ t('dashboard.explore') }} <ArrowRight :size="16"
      /></RouterLink>
    </section>
    <section class="preview-grid" :aria-label="t('dashboard.sections')">
      <RouterLink to="/matches" class="preview-card">
        <CalendarDays />
        <h3>{{ t('nav.matches') }}</h3>
        <p>{{ t('dashboard.matchesIntro') }}</p>
      </RouterLink>
      <RouterLink to="/tv" class="preview-card">
        <Radio />
        <h3>{{ t('tv.title') }}</h3>
        <p>{{ t('dashboard.tvIntro') }}</p>
      </RouterLink>
      <RouterLink to="/standings" class="preview-card">
        <Trophy />
        <h3>{{ t('tables.title') }}</h3>
        <p>{{ t('dashboard.tablesIntro') }}</p>
      </RouterLink>
      <RouterLink to="/news" class="preview-card">
        <Newspaper />
        <h3>{{ t('news.title') }}</h3>
        <p>{{ t('dashboard.newsIntro') }}</p>
      </RouterLink>
    </section>
    <p class="attribution">{{ t('common.footballAttribution') }}</p>
  </main>
</template>
