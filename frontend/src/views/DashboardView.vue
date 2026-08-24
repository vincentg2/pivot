<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowRight, CalendarDays, Newspaper, Radio, Trophy } from 'lucide-vue-next'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import ClubMark from '@/components/ClubMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import NewsCard from '@/components/NewsCard.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import { useSessionStore } from '@/stores/session'
import { useFavoritesStore } from '@/stores/favorites'
import type { NewsItem } from '@/lib/news'
const session = useSessionStore()
const favorites = useFavoritesStore()
const favoriteMatches = ref<FootballMatch[]>([])
const news = ref<NewsItem[]>([])
const todayLabel = computed(() =>
  new Intl.DateTimeFormat('en-GB', { weekday: 'long', day: 'numeric', month: 'long' }).format(
    new Date(),
  ),
)
onMounted(async () => {
  if (!favorites.ready) await favorites.load()
  const from = localDate(new Date())
  const to = localDate(addDays(new Date(), 7))
  const [matchResponse, newsResponse] = await Promise.all([
    api<{ matches: FootballMatch[] }>(`/matches?from=${from}&to=${to}`),
    api<{ news: NewsItem[] }>('/news?limit=6'),
  ])
  favoriteMatches.value = matchResponse.matches.filter((match) => match.favorite).slice(0, 5)
  news.value = newsResponse.news
})
</script>

<template>
  <main class="dashboard page-width">
    <section class="welcome-row">
      <div>
        <p class="eyebrow">{{ todayLabel }}</p>
        <h1>Good evening, {{ session.user?.nickname }}.</h1>
        <p class="lede">Here’s the shape of your football week.</p>
      </div>
      <AvatarMonogram
        v-if="session.user"
        :name="session.user.nickname"
        :seed="session.user.avatarSeed"
        size="lg"
      />
    </section>
    <section v-if="news.length" class="dashboard-news">
      <div class="section-heading">
        <div>
          <p class="eyebrow">From official sources</p>
          <h2>Club dispatches</h2>
        </div>
        <RouterLink to="/news" class="text-link">All news</RouterLink>
      </div>
      <div class="news-grid dashboard-news-grid">
        <NewsCard v-for="item in news" :key="item.id" :item="item" show-club />
      </div>
    </section>
    <section v-if="favorites.clubs.length" class="dashboard-matches">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Next for your clubs</p>
          <h2>The week ahead</h2>
        </div>
        <RouterLink to="/matches" class="text-link">All matches</RouterLink>
      </div>
      <div v-if="favoriteMatches.length" class="compact-matches">
        <MatchRow
          v-for="(match, index) in favoriteMatches"
          :key="`${match.utcDate}-${index}`"
          :match="match"
        />
      </div>
      <p v-else class="quiet dashboard-empty-line">
        No matches synchronized for your favorites this week.
      </p>
    </section>
    <section v-if="favorites.ready && favorites.clubs.length" class="favorite-dashboard">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Your clubs</p>
          <h2>Close to home</h2>
        </div>
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
        <p class="eyebrow">Matches</p>
        <h3>The week ahead</h3>
        <p>Yesterday, today, tomorrow and the next seven days.</p>
      </RouterLink>
      <RouterLink to="/tv" class="preview-card">
        <Radio />
        <p class="eyebrow">On television</p>
        <h3>Where to watch</h3>
        <p>French listings, centrally collected with operator consent.</p>
      </RouterLink>
      <RouterLink to="/standings" class="preview-card">
        <Trophy />
        <p class="eyebrow">Tables</p>
        <h3>The wider picture</h3>
        <p>Domestic standings with clear source attribution.</p>
      </RouterLink>
      <RouterLink to="/news" class="preview-card">
        <Newspaper />
        <p class="eyebrow">Official news</p>
        <h3>From the source</h3>
        <p>Club headlines and links, retained for thirty days.</p>
      </RouterLink>
    </section>
    <p class="attribution">Football data provided by football-data.org.</p>
  </main>
</template>
