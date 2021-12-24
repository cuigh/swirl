<template>
  <div class="tab-pane" :class="{ active: active }" :style="{ borderColor: themeVars.primaryColor }">
    <router-link
      :to="href"
      :style="{ color: active ? themeVars.primaryColor : themeVars.textColorBase }"
    >
      <slot />
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { useThemeVars } from "naive-ui";
import { RouterLink, RouteLocationRaw } from "vue-router";

interface Props {
  active?: boolean;
  title?: string;
  href: RouteLocationRaw;
}

const props = defineProps<Props>()
const themeVars = useThemeVars();
</script>

<style scoped>
.tab-pane {
  padding-bottom: 4px;
  margin-bottom: -1px;
}
.active {
  border-bottom: 2px solid;
}
a {
  text-decoration: none;
}
</style>