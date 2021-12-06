<template>
  <n-result size="huge" status="404" :title="t('titles.404')" :description="t('texts.404')">
    <template #footer>
      <n-button type="primary" @click="$router.push('/')">{{ t('buttons.home') }}</n-button>
    </template>
  </n-result>
</template>

<script setup lang="ts">
import { NResult, NButton } from "naive-ui";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
  