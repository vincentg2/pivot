<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import ClubMark from '@/components/ClubMark.vue'
import type { FootballMatch } from '@/lib/football'

const props = defineProps<{ match: FootballMatch }>()
const { locale, t } = useI18n()
const date = computed(() =>
  new Intl.DateTimeFormat(locale.value, { day: 'numeric', month: 'short' }).format(
    new Date(props.match.utcDate),
  ),
)
const scorers = computed(() =>
  (props.match.goals ?? [])
    .map((goal) => {
      const minute = `${goal.minute}${goal.injuryTime ? `+${goal.injuryTime}` : ''}′`
      const type =
        goal.type === 'PENALTY'
          ? ` (${t('matches.penalty')})`
          : goal.type === 'OWN'
            ? ` (${t('matches.ownGoal')})`
            : ''
      return `${goal.scorerName} ${minute}${type}`
    })
    .join(' · '),
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
    <p v-if="scorers" class="recent-scorers">{{ scorers }}</p>
  </article>
</template>
