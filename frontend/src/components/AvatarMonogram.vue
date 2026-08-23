<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ name: string; seed: string; size?: 'sm' | 'lg' }>()
const initials = computed(() =>
  props.name
    .split(/\s+/)
    .map((part) => part[0])
    .join('')
    .slice(0, 2)
    .toUpperCase(),
)
const hue = computed(
  () => ([...props.seed].reduce((sum, char) => sum + char.charCodeAt(0), 0) % 45) + 8,
)
</script>

<template>
  <span class="avatar" :class="size" :style="{ '--avatar-hue': hue }" aria-hidden="true">{{
    initials
  }}</span>
</template>
