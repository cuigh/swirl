<template>
  <div class="page-header">
    <div :class="['title', store.state.preference.theme === 'dark' ? 'dark' : 'light']">
      <div>
        <n-h3>{{ title ?? t('titles.' + ($route.name as string)) }}</n-h3>
        <n-text depth="3" v-if="subtitle">{{ subtitle }}</n-text>
      </div>
      <n-space :size="6" justify="end" v-if="$slots.action">
        <slot name="action"></slot>
      </n-space>
    </div>
    <div
      :class="['summary', store.state.preference.theme === 'dark' ? 'dark' : 'light']"
      v-if="$slots.default"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NSpace,
  NH3,
  NText,
} from "naive-ui";
import type { PropType } from "vue";
import { useStore } from "vuex";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps({
  title: {
    type: String as PropType<unknown | string>,
  },
  subtitle: {
    type: String,
  }
})
const store = useStore()
</script>

<style lang="scss" scoped>
.page-header {
  box-shadow: 0 2px 4px -2px rgb(10 10 10 / 10%);
  .title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 16px 6px 16px;
    min-height: 30px;
    .n-h3 {
      margin-right: 8px;
      display: inline;
    }
  }
  .summary {
    padding: 12px 16px;
  }
  .title.light,
  .summary.light {
    // background-color: rgba(250, 250, 252, 0.75);
    border-bottom: 1px solid rgb(239, 239, 245);
  }
  .title.dark,
  .summary.dark {
    background-color: rgba(38, 38, 42, 1);
    border-bottom: 1px solid rgba(255, 255, 255, 0.09);
  }
}
</style>