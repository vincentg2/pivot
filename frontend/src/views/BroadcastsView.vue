<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ExternalLink, Radio, Tv } from 'lucide-vue-next'
import ChannelMark from '@/components/ChannelMark.vue'
import { api } from '@/lib/api'
import { addDays, localDate } from '@/lib/football'
import { listingTime, type BroadcastListing } from '@/lib/broadcast'

const listings = ref<BroadcastListing[]>([])
const from = ref(localDate(new Date()))
const to = ref(localDate(addDays(new Date(), 7)))
const loading = ref(true)
const { locale, t } = useI18n()
const grouped = computed(() => {
  const groups = new Map<string, BroadcastListing[]>()
  for (const listing of listings.value) {
    const date = localDate(new Date(listing.startsAt))
    groups.set(date, [...(groups.get(date) || []), listing])
  }
  return [...groups.entries()]
})

async function load() {
  loading.value = true
  const query = new URLSearchParams({ from: from.value, to: to.value })
  listings.value = (await api<{ listings: BroadcastListing[] }>(`/broadcasts?${query}`)).listings
  loading.value = false
}
function choose(start: Date, days: number) {
  from.value = localDate(start)
  to.value = localDate(addDays(start, days))
  void load()
}
function chooseTwoMonths() {
  const end = new Date()
  end.setMonth(end.getMonth() + 2)
  from.value = localDate(new Date())
  to.value = localDate(end)
  void load()
}
function dateHeading(value: string) {
  return new Intl.DateTimeFormat(locale.value, {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  }).format(new Date(`${value}T12:00:00`))
}
onMounted(load)
</script>

<template>
  <main class="broadcast-page page-width">
    <header class="matches-heading broadcast-heading">
      <h1>{{ t('tv.title') }}</h1>
      <Tv :size="44" aria-hidden="true" />
    </header>
    <section class="date-presets" aria-label="TV schedule date range">
      <button @click="choose(new Date(), 0)">{{ t('tv.today') }}</button>
      <button @click="choose(addDays(new Date(), 1), 0)">{{ t('tv.tomorrow') }}</button>
      <button @click="choose(new Date(), 7)">{{ t('tv.next7') }}</button>
      <button @click="chooseTwoMonths">{{ t('tv.next2Months') }}</button>
    </section>
    <p v-if="loading" class="catalog-state">{{ t('tv.loading') }}</p>
    <div v-else-if="grouped.length" class="broadcast-groups">
      <section v-for="[date, items] in grouped" :key="date">
        <div class="match-date">
          <h2>{{ dateHeading(date) }}</h2>
          <span>{{ t('tv.count', { count: items.length }) }}</span>
        </div>
        <article v-for="listing in items" :key="listing.id" class="broadcast-row">
          <time :datetime="listing.startsAt">{{ listingTime(listing.startsAt) }}</time>
          <div class="broadcast-copy">
            <div class="broadcast-title">
              <strong>{{ listing.label }}</strong>
              <span v-if="listing.kind !== 'live'" class="status">{{ listing.kind }}</span>
              <span v-if="listing.external" class="external-data"
                ><Radio :size="13" /> External listing</span
              >
            </div>
            <span>{{ listing.competitionName || t('tv.competitionUnknown') }}</span>
          </div>
          <div class="channel-list">
            <ChannelMark v-for="channel in listing.channels" :key="channel" :channel="channel" />
          </div>
          <a
            v-if="listing.sourceUrl"
            :href="listing.sourceUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="icon-link"
            :aria-label="t('tv.openFootao')"
            ><ExternalLink :size="17"
          /></a>
        </article>
      </section>
    </div>
    <section v-else class="empty-hero matches-empty">
      <Tv :size="28" />
      <h2>{{ t('tv.empty') }}</h2>
      <p>{{ t('tv.emptyHelp') }}</p>
    </section>
    <p class="attribution">
      French TV schedule provided by
      <a href="https://www.footao.tv/" target="_blank" rel="noopener noreferrer">Footao</a>.
      External listings are not enriched with Pivot club data.
    </p>
  </main>
</template>
