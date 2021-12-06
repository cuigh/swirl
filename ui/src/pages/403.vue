<template>
  <n-result size="huge" status="403" :title="t('titles.403')" :description="t('texts.403')">
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
