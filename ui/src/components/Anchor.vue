<template>
  <router-link :to="url" #="{ navigate, href }" custom>
    <n-a :href="href" @click="navigate">
      <slot />
    </n-a>
  </router-link>
</template>

<script setup lang="ts">
import { NA } from "naive-ui";
import { RouterLink } from "vue-router";
import type { RouteLocationRaw } from "vue-router";

interface Props {
  url: RouteLocationRaw
}

const props = defineProps<Props>()
</script>
