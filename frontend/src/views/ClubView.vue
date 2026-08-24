<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ArrowLeft, ExternalLink, Heart, MapPin } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import ClubMark from '@/components/ClubMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import { useFavoritesStore, type Club } from '@/stores/favorites'

const route = useRoute()
const favorites = useFavoritesStore()
const club = ref<Club | null>(null)
const notice = ref('')
const matches = ref<FootballMatch[]>([])

async function toggle() {
  if (!club.value) return
  try {
    await favorites.toggle(club.value)
  } catch (caught) {
    notice.value = caught instanceof Error ? caught.message : 'Favorite could not be saved.'
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
})
</script>

<template>
  <main class="club-page page-width">
    <RouterLink to="/clubs" class="back-link"><ArrowLeft :size="16" /> All clubs</RouterLink>
    <p v-if="!club" class="catalog-state">Loading club…</p>
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
          {{ favorites.has(club.id) ? 'Following' : 'Follow club' }}
        </button>
      </header>
      <p v-if="notice" class="form-error" role="alert">{{ notice }}</p>
      <section class="club-details">
        <article>
          <p class="eyebrow">Home</p>
          <h2><MapPin :size="20" /> {{ club.venue || 'Venue not listed' }}</h2>
        </article>
        <article>
          <p class="eyebrow">Competitions</p>
          <h2>{{ club.competitions?.map((item) => item.name).join(', ') || 'Not listed' }}</h2>
        </article>
        <article v-if="club.websiteUrl">
          <p class="eyebrow">Official</p>
          <a :href="club.websiteUrl" target="_blank" rel="noopener noreferrer"
            >Visit club website <ExternalLink :size="15"
          /></a>
        </article>
      </section>
      <section class="club-fixtures">
        <div class="section-heading">
          <div>
            <p class="eyebrow">Recent & upcoming</p>
            <h2>Matches</h2>
          </div>
          <RouterLink to="/matches" class="text-link">All matches</RouterLink>
        </div>
        <div v-if="matches.length" class="compact-matches">
          <MatchRow
            v-for="(match, index) in matches"
            :key="`${match.utcDate}-${index}`"
            :match="match"
          />
        </div>
        <p v-else class="quiet dashboard-empty-line">No matches synchronized in this window.</p>
      </section>
      <p class="attribution">Football data provided by football-data.org.</p>
    </template>
  </main>
</template>
