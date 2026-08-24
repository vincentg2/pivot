<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
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
  new Intl.DateTimeFormat('en-GB', { weekday: 'long', day: 'numeric', month: 'long' }).format(
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
        <h1>Good evening, {{ session.user?.nickname }}.</h1>
        <p class="dashboard-date">{{ todayLabel }}</p>
      </div>
      <AvatarMonogram
        v-if="session.user"
        :name="session.user.nickname"
        :seed="session.user.avatarSeed"
        size="lg"
      />
      <div v-if="recentResults.length" class="recent-results">
        <p>Last results</p>
        <div class="recent-results-list">
          <RecentResult v-for="match in recentResults" :key="match.id" :match="match" />
        </div>
      </div>
    </section>
    <section v-if="news.length" class="dashboard-news">
      <div class="section-heading">
        <h2>Latest club news</h2>
        <RouterLink to="/news" class="text-link">All news</RouterLink>
      </div>
      <div class="news-grid dashboard-news-grid">
        <NewsCard v-for="item in news" :key="item.id" :item="item" show-club />
      </div>
    </section>
    <section v-if="favorites.clubs.length" class="dashboard-matches">
      <div class="section-heading">
        <h2>Upcoming matches <span>30 days</span></h2>
        <RouterLink to="/matches" class="text-link">All matches</RouterLink>
      </div>
      <div class="favorite-club-filters" role="group" aria-label="Filter matches by favorite club">
        <button
          type="button"
          :class="{ active: selectedClubId === null }"
          :aria-pressed="selectedClubId === null"
          @click="selectedClubId = null"
        >
          All clubs
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
        No matches synchronized for this selection in this period.
      </p>
      <div v-if="hasMoreMatches" class="dashboard-match-more">
        <button type="button" class="button secondary" @click="visibleMatchCount += matchBatchSize">
          Show next matches <ChevronDown :size="17" />
        </button>
        <span aria-live="polite">{{ favoriteMatches.length }} of {{ filteredMatches.length }}</span>
      </div>
    </section>
    <section v-if="favorites.ready && favorites.clubs.length" class="favorite-dashboard">
      <div class="section-heading">
        <h2>Favourite clubs</h2>
        <RouterLink to="/clubs" class="text-link">Edit favorites</RouterLink>
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
      <p class="eyebrow">Your starting eleven</p>
      <h2>Choose the clubs you care about</h2>
      <p>Browse the club catalog and select up to five favorites for your dashboard.</p>
      <RouterLink to="/clubs" class="button primary"
        >Explore clubs <ArrowRight :size="16"
      /></RouterLink>
    </section>
    <section class="preview-grid" aria-label="Upcoming Pivot sections">
      <RouterLink to="/matches" class="preview-card">
        <CalendarDays />
        <h3>Matches</h3>
        <p>Yesterday, today, tomorrow and the next seven days.</p>
      </RouterLink>
      <RouterLink to="/tv" class="preview-card">
        <Radio />
        <h3>TV schedule</h3>
        <p>French listings, centrally collected with operator consent.</p>
      </RouterLink>
      <RouterLink to="/standings" class="preview-card">
        <Trophy />
        <h3>Tables</h3>
        <p>Domestic standings with clear source attribution.</p>
      </RouterLink>
      <RouterLink to="/news" class="preview-card">
        <Newspaper />
        <h3>News</h3>
        <p>Club headlines and links, retained for thirty days.</p>
      </RouterLink>
    </section>
    <p class="attribution">Football data provided by football-data.org.</p>
  </main>
</template>
