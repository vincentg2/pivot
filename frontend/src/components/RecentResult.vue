<script setup lang="ts">
import { computed } from 'vue'
import ClubMark from '@/components/ClubMark.vue'
import type { FootballMatch } from '@/lib/football'

const props = defineProps<{ match: FootballMatch }>()
const date = computed(() =>
  new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short' }).format(
    new Date(props.match.utcDate),
  ),
)
const homeWon = computed(() => (props.match.homeScore ?? 0) > (props.match.awayScore ?? 0))
const awayWon = computed(() => (props.match.awayScore ?? 0) > (props.match.homeScore ?? 0))
</script>

<template>
  <article class="recent-result">
    <time :datetime="match.utcDate">{{ date }}</time>
    <div class="recent-scoreline">
      <span class="recent-team" :class="{ winner: homeWon }" :title="match.home.name">
        <ClubMark :name="match.home.name" :tla="match.home.tla" :crest-url="match.home.crestUrl" />
        <span>{{ match.home.shortName || match.home.name }}</span>
      </span>
      <strong>{{ match.homeScore ?? '–' }}–{{ match.awayScore ?? '–' }}</strong>
      <span class="recent-team" :class="{ winner: awayWon }" :title="match.away.name">
        <ClubMark :name="match.away.name" :tla="match.away.tla" :crest-url="match.away.crestUrl" />
        <span>{{ match.away.shortName || match.away.name }}</span>
      </span>
    </div>
  </article>
</template>
