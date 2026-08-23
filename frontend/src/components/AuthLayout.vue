<script setup lang="ts">
defineProps<{ eyebrow: string; title: string; intro: string }>()
</script>

<template>
  <main class="auth-page">
    <section class="auth-story" aria-labelledby="auth-title">
      <RouterLink to="/" class="brand"
        ><span class="brand-mark">P</span><span>Pivot</span></RouterLink
      >
      <div>
        <p class="eyebrow">{{ eyebrow }}</p>
        <h1 id="auth-title">{{ title }}</h1>
        <p class="lede">{{ intro }}</p>
      </div>
      <p class="quiet">Football, considered.</p>
    </section>
    <section class="auth-panel"><slot /></section>
  </main>
</template>
