<script setup lang="ts">
import { onMounted } from 'vue'
import { ArrowRight, CalendarDays, Radio, Trophy } from 'lucide-vue-next'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import ClubMark from '@/components/ClubMark.vue'
import { useSessionStore } from '@/stores/session'
import { useFavoritesStore } from '@/stores/favorites'
const session = useSessionStore()
const favorites = useFavoritesStore()
onMounted(() => {
  if (!favorites.ready) void favorites.load()
})
</script>

<template>
  <main class="dashboard page-width">
    <section class="welcome-row">
      <div>
        <p class="eyebrow">Monday, 24 August</p>
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
      <article>
        <CalendarDays />
        <p class="eyebrow">Matches</p>
        <h3>The week ahead</h3>
        <p>Yesterday, today, tomorrow and the next seven days.</p>
      </article>
      <article>
        <Radio />
        <p class="eyebrow">On television</p>
        <h3>Where to watch</h3>
        <p>French listings, centrally collected with operator consent.</p>
      </article>
      <article>
        <Trophy />
        <p class="eyebrow">Tables</p>
        <h3>The wider picture</h3>
        <p>Domestic standings with clear source attribution.</p>
      </article>
    </section>
    <p class="attribution">
      Data connectors are currently off. Future match data: football-data.org.
    </p>
  </main>
</template>
