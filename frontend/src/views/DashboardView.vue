<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowRight, CalendarDays, Radio, Trophy } from 'lucide-vue-next'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import ClubMark from '@/components/ClubMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import { useSessionStore } from '@/stores/session'
import { useFavoritesStore } from '@/stores/favorites'
const session = useSessionStore()
const favorites = useFavoritesStore()
const favoriteMatches = ref<FootballMatch[]>([])
const todayLabel = computed(() =>
  new Intl.DateTimeFormat('en-GB', { weekday: 'long', day: 'numeric', month: 'long' }).format(
    new Date(),
  ),
)
onMounted(async () => {
  if (!favorites.ready) await favorites.load()
  const from = localDate(new Date())
  const to = localDate(addDays(new Date(), 7))
  const response = await api<{ matches: FootballMatch[] }>(`/matches?from=${from}&to=${to}`)
  favoriteMatches.value = response.matches.filter((match) => match.favorite).slice(0, 5)
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
      <article>
        <Radio />
        <p class="eyebrow">On television</p>
        <h3>Where to watch</h3>
        <p>French listings, centrally collected with operator consent.</p>
      </article>
      <RouterLink to="/standings" class="preview-card">
        <Trophy />
        <p class="eyebrow">Tables</p>
        <h3>The wider picture</h3>
        <p>Domestic standings with clear source attribution.</p>
      </RouterLink>
    </section>
    <p class="attribution">Football data provided by football-data.org.</p>
  </main>
</template>
