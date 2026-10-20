<template>
  <AppLayoutV2 v-if="uiVersionStore.isV2" />
  <AppLayoutV1 v-else />
</template>

<script setup lang="ts">
import AppLayoutV1 from './AppLayoutV1.vue'
import AppLayoutV2 from '@/layouts/v2/AppLayoutV2.vue'
import { useUiVersionStore } from '@/stores/uiVersion.store'

const uiVersionStore = useUiVersionStore()
</script>
