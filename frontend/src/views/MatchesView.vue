<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CalendarDays, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import ChannelMark from '@/components/ChannelMark.vue'
import MatchRow from '@/components/MatchRow.vue'
import { api } from '@/lib/api'
import { addDays, localDate, type FootballMatch } from '@/lib/football'
import type { Competition } from '@/stores/favorites'
import { broadcastChannelsByMatch, listingTime, type BroadcastListing } from '@/lib/broadcast'

const matches = ref<FootballMatch[]>([])
const broadcasts = ref<BroadcastListing[]>([])
const competitions = ref<Competition[]>([])
const from = ref(localDate(new Date()))
const to = ref(from.value)
const competition = ref('')
const favoritesOnly = ref(false)
const loading = ref(true)
const { locale, t } = useI18n()

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
const channelsByMatch = computed(() => broadcastChannelsByMatch(broadcasts.value))
const externalByDate = computed(() => {
  const map = new Map<string, BroadcastListing[]>()
  if (favoritesOnly.value) return map
  for (const listing of broadcasts.value.filter((item) => item.external)) {
    const date = localDate(new Date(listing.startsAt))
    map.set(date, [...(map.get(date) || []), listing])
  }
  return map
})
const dates = computed(() =>
  [
    ...new Set([...groupedMatches.value.map(([date]) => date), ...externalByDate.value.keys()]),
  ].sort(),
)
const matchesByDate = computed(() => new Map(groupedMatches.value))

async function load() {
  loading.value = true
  const query = new URLSearchParams({ from: from.value, to: to.value })
  if (competition.value) query.set('competition', competition.value)
  const [matchResponse, broadcastResponse] = await Promise.all([
    api<{ matches: FootballMatch[] }>(`/matches?${query}`),
    api<{ listings: BroadcastListing[] }>(
      `/broadcasts?${new URLSearchParams({ from: from.value, to: to.value })}`,
    ),
  ])
  matches.value = matchResponse.matches
  broadcasts.value = broadcastResponse.listings
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
  return new Intl.DateTimeFormat(locale.value, {
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
      <h1>{{ t('matches.title') }}</h1>
    </header>
    <section class="date-presets" aria-label="Match date range">
      <button @click="choose(addDays(new Date(), -1))">{{ t('matches.yesterday') }}</button>
      <button @click="choose(new Date())">{{ t('matches.today') }}</button>
      <button @click="choose(addDays(new Date(), 1))">{{ t('matches.tomorrow') }}</button>
      <button @click="choose(new Date(), 7)">{{ t('matches.next7') }}</button>
    </section>
    <section class="match-filters" aria-label="Match filters">
      <button class="icon-link" :aria-label="t('matches.previous')" @click="shift(-1)">
        <ChevronLeft :size="19" />
      </button>
      <label
        ><span class="sr-only">Start date</span
        ><input v-model="from" type="date" @change="changeStart"
      /></label>
      <button class="icon-link" :aria-label="t('matches.next')" @click="shift(1)">
        <ChevronRight :size="19" />
      </button>
      <label
        ><span class="sr-only">Competition</span
        ><select v-model="competition" @change="load">
          <option value="">{{ t('matches.allCompetitions') }}</option>
          <option v-for="item in competitions" :key="item.id" :value="item.code">
            {{ item.name }}
          </option>
        </select></label
      >
      <label class="check-filter"
        ><input v-model="favoritesOnly" type="checkbox" /> {{ t('matches.myClubs') }}</label
      >
    </section>
    <p v-if="loading" class="catalog-state">Loading matches…</p>
    <div v-else-if="dates.length" class="match-groups">
      <section v-for="date in dates" :key="date">
        <div class="match-date">
          <h2>{{ dateHeading(date) }}</h2>
          <span>{{ t('matches.count', { count: matchesByDate.get(date)?.length || 0 }) }}</span>
        </div>
        <MatchRow
          v-for="match in matchesByDate.get(date) || []"
          :key="match.id"
          :match="match"
          :channels="channelsByMatch.get(match.id)"
        />
        <article
          v-for="listing in externalByDate.get(date) || []"
          :key="listing.id"
          class="external-match-row"
        >
          <time :datetime="listing.startsAt">{{ listingTime(listing.startsAt) }}</time>
          <div>
            <strong>{{ listing.label }}</strong
            ><span>{{ listing.competitionName || 'External TV listing' }} · not enriched</span>
          </div>
          <span class="channel-list"
            ><ChannelMark v-for="channel in listing.channels" :key="channel" :channel="channel"
          /></span>
        </article>
      </section>
    </div>
    <section v-else class="empty-hero matches-empty">
      <CalendarDays :size="28" />
      <h2>{{ t('matches.empty') }}</h2>
      <p>{{ t('matches.emptyHelp') }}</p>
    </section>
    <p class="attribution">
      Football data provided by football-data.org. TV schedule provided by Footao. Times shown in
      your local timezone.
    </p>
  </main>
</template>
