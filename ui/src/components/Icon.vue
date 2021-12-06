<template>
  <n-icon :size="size" :color="color" style="vertical-align: middle">
    <svg xmlns="http://www.w3.org/2000/svg" :viewBox="viewBox">
      <path :d="p" fill="currentColor" v-if="Array.isArray(path)" v-for="p in path"/>
      <path :d="path" fill="currentColor" v-else/>
    </svg>
  </n-icon>
</template>

<script setup lang="ts">
import { NIcon } from "naive-ui";
import { PropType } from "vue";

const props = defineProps({
  size: {
    type: [Number, String] as PropType<number | string>,
  },
  color: {
    type: String,
  },  
  viewBox: {
    type: String,
    required: true,
  },
  path: {
    type: [String, Array] as PropType<string | string[]>,
    // type: String,
    required: true,
  },
})
</script>
