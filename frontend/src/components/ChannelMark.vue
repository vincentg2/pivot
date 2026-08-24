<script setup lang="ts">
import { computed, ref } from 'vue'
import { channelIconUrl } from '@/lib/channel'

const props = defineProps<{ channel: string }>()
const imageFailed = ref(false)
const remoteLogosEnabled = import.meta.env.VITE_REMOTE_LOGOS_ENABLED === 'true'
const logoUrl = computed(() => channelIconUrl(props.channel))
</script>

<template>
  <span class="channel-mark">
    <span v-if="remoteLogosEnabled && logoUrl && !imageFailed" class="channel-logo-wrap">
      <img
        :src="logoUrl"
        alt=""
        loading="lazy"
        referrerpolicy="no-referrer"
        @error="imageFailed = true"
      />
    </span>
    <span>{{ channel }}</span>
  </span>
</template>
