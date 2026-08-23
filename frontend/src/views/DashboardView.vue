<script setup lang="ts">
import { ArrowRight, CalendarDays, Radio, Trophy } from 'lucide-vue-next'
import AvatarMonogram from '@/components/AvatarMonogram.vue'
import { useSessionStore } from '@/stores/session'
const session = useSessionStore()
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
    <section class="empty-hero">
      <div class="empty-icon"><Trophy :size="26" /></div>
      <p class="eyebrow">Your starting eleven</p>
      <h2>Choose the clubs you care about</h2>
      <p>
        Club catalog and up to five favorites arrive in milestone two. Your private profile is ready
        for them.
      </p>
      <span class="button disabled">Coming next <ArrowRight :size="16" /></span>
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
