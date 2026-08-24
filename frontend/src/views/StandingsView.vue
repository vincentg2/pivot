<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ClubMark from '@/components/ClubMark.vue'
import { ApiError, api } from '@/lib/api'
import type { Standing } from '@/lib/football'
import type { Competition } from '@/stores/favorites'

const competitions = ref<Competition[]>([])
const selected = ref('')
const standing = ref<Standing | null>(null)
const error = ref('')
async function load() {
  if (!selected.value) return
  error.value = ''
  try {
    standing.value = (
      await api<{ standing: Standing }>(`/standings?competition=${selected.value}`)
    ).standing
  } catch (caught) {
    standing.value = null
    error.value = caught instanceof ApiError ? caught.message : 'Standing unavailable.'
  }
}
onMounted(async () => {
  competitions.value = (await api<{ competitions: Competition[] }>('/competitions')).competitions
  if (competitions.value.length) {
    selected.value = competitions.value[0].code
    await load()
  }
})
</script>

<template>
  <main class="standings-page page-width">
    <header class="standings-heading">
      <h1>Tables</h1>
      <label
        >Competition<select v-model="selected" @change="load">
          <option v-for="item in competitions" :key="item.id" :value="item.code">
            {{ item.name }}
          </option>
        </select></label
      >
    </header>
    <p v-if="error" class="form-error" role="alert">{{ error }}</p>
    <section v-if="standing" class="table-card">
      <div class="table-title">
        <div>
          <p class="eyebrow">{{ standing.competition.code }}</p>
          <h2>{{ standing.competition.name }}</h2>
        </div>
        <span
          >{{ new Date(standing.season.startDate).getFullYear() }}–{{
            new Date(standing.season.endDate).getFullYear()
          }}</span
        >
      </div>
      <div class="standing-table" role="table" aria-label="League standing">
        <div class="standing-row standing-head" role="row">
          <span>#</span><span>Club</span><span>P</span><span class="wide-stat">W</span
          ><span class="wide-stat">D</span><span class="wide-stat">L</span><span>GD</span
          ><strong>Pts</strong>
        </div>
        <div
          v-for="row in standing.rows"
          :key="row.club.id || row.position"
          class="standing-row"
          role="row"
        >
          <span>{{ row.position }}</span
          ><span class="standing-club"
            ><ClubMark
              :name="row.club.name"
              :tla="row.club.tla"
              :crest-url="row.club.crestUrl"
            /><RouterLink v-if="row.club.id" :to="`/clubs/${row.club.id}`">{{
              row.club.shortName || row.club.name
            }}</RouterLink
            ><span v-else>{{ row.club.name }}</span></span
          ><span>{{ row.played }}</span
          ><span class="wide-stat">{{ row.won }}</span
          ><span class="wide-stat">{{ row.drawn }}</span
          ><span class="wide-stat">{{ row.lost }}</span
          ><span>{{ row.goalDifference > 0 ? '+' : '' }}{{ row.goalDifference }}</span
          ><strong>{{ row.points }}</strong>
        </div>
      </div>
    </section>
    <section v-else-if="!error" class="empty-hero">
      <h2>No table synchronized yet</h2>
      <p>Run the sports data collection from administration.</p>
    </section>
    <p class="attribution">Football data provided by football-data.org.</p>
  </main>
</template>
