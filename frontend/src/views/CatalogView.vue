<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Heart, Search } from 'lucide-vue-next'
import ClubMark from '@/components/ClubMark.vue'
import { ApiError, api } from '@/lib/api'
import { useFavoritesStore, type Club, type Competition } from '@/stores/favorites'

const favorites = useFavoritesStore()
const competitions = ref<Competition[]>([])
const clubs = ref<Club[]>([])
const selectedCompetition = ref('')
const search = ref('')
const loading = ref(true)
const error = ref('')
const notice = ref('')

const filteredClubs = computed(() => {
  const query = search.value.trim().toLocaleLowerCase()
  if (!query) return clubs.value
  return clubs.value.filter((club) =>
    [club.name, club.shortName, club.tla].some((value) =>
      value.toLocaleLowerCase().includes(query),
    ),
  )
})

async function loadClubs() {
  loading.value = true
  error.value = ''
  try {
    const query = selectedCompetition.value
      ? `?competition=${encodeURIComponent(selectedCompetition.value)}`
      : ''
    clubs.value = (await api<{ clubs: Club[] }>(`/clubs${query}`)).clubs
  } catch (caught) {
    error.value =
      caught instanceof ApiError ? caught.message : 'The club catalog could not be loaded.'
  } finally {
    loading.value = false
  }
}

function selectCompetition(code: string) {
  selectedCompetition.value = code
  void loadClubs()
}

async function toggle(club: Club) {
  notice.value = ''
  try {
    await favorites.toggle(club)
  } catch (caught) {
    notice.value = caught instanceof Error ? caught.message : 'Favorite could not be saved.'
  }
}

onMounted(async () => {
  const [competitionResponse] = await Promise.all([
    api<{ competitions: Competition[] }>('/competitions'),
    favorites.ready ? Promise.resolve() : favorites.load(),
  ])
  competitions.value = competitionResponse.competitions
  await loadClubs()
})
</script>

<template>
  <main class="catalog page-width">
    <header class="catalog-heading">
      <h1>Clubs</h1>
      <p class="favorite-count">
        <strong>{{ favorites.clubs.length }}</strong> / 5 favorites
      </p>
    </header>

    <section class="catalog-tools" aria-label="Catalog filters">
      <div class="search-field">
        <Search :size="18" />
        <label class="sr-only" for="club-search">Search clubs</label>
        <input id="club-search" v-model="search" type="search" placeholder="Search a club" />
      </div>
      <div class="competition-tabs" aria-label="Filter by competition">
        <button :class="{ active: !selectedCompetition }" @click="selectCompetition('')">
          All
        </button>
        <button
          v-for="competition in competitions"
          :key="competition.id"
          :class="{ active: selectedCompetition === competition.code }"
          @click="selectCompetition(competition.code)"
        >
          {{ competition.name }}
        </button>
      </div>
    </section>

    <p v-if="notice || error" class="form-error" role="alert">{{ notice || error }}</p>
    <p v-if="loading" class="catalog-state">Loading the local catalog…</p>
    <section v-else-if="filteredClubs.length" class="club-grid" aria-label="Clubs">
      <article v-for="club in filteredClubs" :key="club.id" class="club-card">
        <RouterLink :to="`/clubs/${club.id}`" class="club-card-main">
          <ClubMark :name="club.name" :tla="club.tla" :crest-url="club.crestUrl" size="lg" />
          <span
            ><strong>{{ club.shortName || club.name }}</strong
            ><small>{{ club.venue || club.tla }}</small></span
          >
        </RouterLink>
        <button
          class="favorite-button"
          :class="{ selected: favorites.has(club.id) }"
          :aria-label="
            favorites.has(club.id)
              ? `Remove ${club.name} from favorites`
              : `Add ${club.name} to favorites`
          "
          :aria-pressed="favorites.has(club.id)"
          :disabled="favorites.saving"
          @click="toggle(club)"
        >
          <Heart :size="18" :fill="favorites.has(club.id) ? 'currentColor' : 'none'" />
        </button>
      </article>
    </section>
    <section v-else class="empty-hero catalog-empty">
      <p class="eyebrow">Local catalog</p>
      <h2>No clubs synchronized yet</h2>
      <p>
        An administrator can enable football-data.org and start the first collection. Pivot remains
        fully usable without an external key.
      </p>
    </section>
    <p class="attribution">
      Football data provided by football-data.org. Remote crests are disabled by default.
    </p>
  </main>
</template>
