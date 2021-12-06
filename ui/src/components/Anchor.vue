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

const props = defineProps({
  url: {
    type: String,
    required: true,
  }
})
</script>
