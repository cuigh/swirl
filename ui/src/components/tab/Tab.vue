<template>
  <div class="tab" :style="{borderColor: themeVars.dividerColor}">
    <n-space :size="24">
      <slot />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { NSpace, useThemeVars } from "naive-ui";

const themeVars = useThemeVars()
</script>

<style scoped>
.tab {
  border-bottom: 1px solid;
}
</style>