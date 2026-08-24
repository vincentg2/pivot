<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{ name: string; tla?: string; crestUrl?: string | null; size?: 'sm' | 'lg' }>(),
  { tla: '', crestUrl: null, size: 'sm' },
)
const imageFailed = ref(false)
const remoteLogosEnabled = import.meta.env.VITE_REMOTE_LOGOS_ENABLED === 'true'
const monogram = (props.tla || props.name.slice(0, 3)).toUpperCase()
</script>

<template>
  <span class="club-mark" :class="size" aria-hidden="true">
    <img
      v-if="remoteLogosEnabled && crestUrl && !imageFailed"
      :src="crestUrl"
      alt=""
      loading="lazy"
      referrerpolicy="no-referrer"
      @error="imageFailed = true"
    />
    <span v-else>{{ monogram }}</span>
  </span>
</template>
