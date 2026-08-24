<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { CalendarDays, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import MatchRow from '@/components/MatchRow.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import type { Competition } from '@/stores/favorites'

const matches = ref<FootballMatch[]>([])
const competitions = ref<Competition[]>([])
const from = ref(localDate(new Date()))
const to = ref(from.value)
const competition = ref('')
const favoritesOnly = ref(false)
const loading = ref(true)

const visibleMatches = computed(() =>
  favoritesOnly.value ? matches.value.filter((match) => match.favorite) : matches.value,
)
const groupedMatches = computed(() => {
  const groups = new Map<string, FootballMatch[]>()
  for (const match of visibleMatches.value) {
    const key = localDate(new Date(match.utcDate))
    groups.set(key, [...(groups.get(key) || []), match])
  }
  return [...groups.entries()]
})

async function load() {
  loading.value = true
  const query = new URLSearchParams({ from: from.value, to: to.value })
  if (competition.value) query.set('competition', competition.value)
  matches.value = (await api<{ matches: FootballMatch[] }>(`/matches?${query}`)).matches
  loading.value = false
}
function choose(start: Date, days = 0) {
  from.value = localDate(start)
  to.value = localDate(addDays(start, days))
  void load()
}
function shift(days: number) {
  const start = addDays(new Date(`${from.value}T12:00:00`), days)
  const span = Math.round(
    (new Date(`${to.value}T12:00:00`).getTime() - new Date(`${from.value}T12:00:00`).getTime()) /
      86400000,
  )
  choose(start, span)
}
function changeStart() {
  to.value = from.value
  void load()
}
function dateHeading(value: string) {
  return new Intl.DateTimeFormat('en-GB', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  }).format(new Date(`${value}T12:00:00`))
}
onMounted(async () => {
  competitions.value = (await api<{ competitions: Competition[] }>('/competitions')).competitions
  await load()
})
</script>

<template>
  <main class="matches-page page-width">
    <header class="matches-heading">
      <p class="eyebrow">Fixtures & results</p>
      <h1>All matches.</h1>
      <p class="lede">The scoreline, the schedule, and nothing noisy around it.</p>
    </header>
    <section class="date-presets" aria-label="Match date range">
      <button @click="choose(addDays(new Date(), -1))">Yesterday</button>
      <button @click="choose(new Date())">Today</button>
      <button @click="choose(addDays(new Date(), 1))">Tomorrow</button>
      <button @click="choose(new Date(), 7)">Next 7 days</button>
    </section>
    <section class="match-filters" aria-label="Match filters">
      <button class="icon-link" aria-label="Previous period" @click="shift(-1)">
        <ChevronLeft :size="19" />
      </button>
      <label
        ><span class="sr-only">Start date</span
        ><input v-model="from" type="date" @change="changeStart"
      /></label>
      <button class="icon-link" aria-label="Next period" @click="shift(1)">
        <ChevronRight :size="19" />
      </button>
      <label
        ><span class="sr-only">Competition</span
        ><select v-model="competition" @change="load">
          <option value="">All competitions</option>
          <option v-for="item in competitions" :key="item.id" :value="item.code">
            {{ item.name }}
          </option>
        </select></label
      >
      <label class="check-filter"><input v-model="favoritesOnly" type="checkbox" /> My clubs</label>
    </section>
    <p v-if="loading" class="catalog-state">Loading matches…</p>
    <div v-else-if="groupedMatches.length" class="match-groups">
      <section v-for="[date, items] in groupedMatches" :key="date">
        <div class="match-date">
          <h2>{{ dateHeading(date) }}</h2>
          <span>{{ items.length }} matches</span>
        </div>
        <MatchRow
          v-for="(match, index) in items"
          :key="`${match.utcDate}-${index}`"
          :match="match"
        />
      </section>
    </div>
    <section v-else class="empty-hero matches-empty">
      <CalendarDays :size="28" />
      <h2>No matches in this window</h2>
      <p>Try another date or ask an administrator to refresh sports data.</p>
    </section>
    <p class="attribution">
      Football data provided by football-data.org. Times shown in your local timezone.
    </p>
  </main>
</template>
