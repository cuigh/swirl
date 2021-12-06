<template>
  <n-layout :position="isMobile ? 'static' : 'absolute'">
    <n-layout-content
      :position="isMobile ? 'static' : 'absolute'"
      :style="isMobile ? '' : 'bottom: 64px'"
    >
      <div class="content">
        <router-view></router-view>
      </div>
    </n-layout-content>
    <n-layout-footer
      :position="isMobile ? 'static' : 'absolute'"
      :style="{ backgroundColor: themeVars.bodyColor }"
    >
      <span>{{ t('copyright') }}</span>
    </n-layout-footer>
  </n-layout>
</template>

<script setup lang="ts">
import {
  NLayout,
  NLayoutContent,
  NLayoutFooter,
  useThemeVars,
} from "naive-ui";
import { RouterView } from "vue-router";
import { useIsMobile } from "@/utils";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const isMobile = useIsMobile()
const themeVars = useThemeVars()
</script>

<style scoped>
.content {
  padding: 48px 24px;
  overflow-x: scroll;
}
</style>
