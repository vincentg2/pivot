<script setup lang="ts">
import { computed } from 'vue'
import ClubMark from '@/components/ClubMark.vue'
import type { FootballMatch } from '@/lib/football'

const props = defineProps<{ match: FootballMatch; channels?: string[] }>()
const finished = computed(() => props.match.status === 'FINISHED')
const time = computed(() =>
  new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit' }).format(
    new Date(props.match.utcDate),
  ),
)
const status = computed(() => {
  if (finished.value) return 'Full time'
  if (['IN_PLAY', 'PAUSED'].includes(props.match.status)) return 'Live'
  if (props.match.status === 'POSTPONED') return 'Postponed'
  return time.value
})
</script>

<template>
  <article class="match-row" :class="{ featured: match.favorite }">
    <div class="match-meta">
      <span>{{ match.competition.code }}</span
      ><time :datetime="match.utcDate">{{ status }}</time>
    </div>
    <div class="match-team home-team">
      <span>{{ match.home.shortName || match.home.name }}</span>
      <ClubMark :name="match.home.name" :tla="match.home.tla" :crest-url="match.home.crestUrl" />
    </div>
    <div class="match-score" :class="{ scheduled: !finished }">
      <template v-if="finished || match.status === 'IN_PLAY'"
        >{{ match.homeScore ?? '–' }}<span>:</span>{{ match.awayScore ?? '–' }}</template
      >
      <span v-else>vs</span>
    </div>
    <div class="match-team away-team">
      <ClubMark :name="match.away.name" :tla="match.away.tla" :crest-url="match.away.crestUrl" />
      <span>{{ match.away.shortName || match.away.name }}</span>
    </div>
    <div v-if="channels?.length" class="match-channels" aria-label="TV channels">
      <span v-for="channel in channels" :key="channel">{{ channel }}</span>
    </div>
  </article>
</template>
