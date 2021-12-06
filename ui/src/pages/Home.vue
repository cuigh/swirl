<template>
  <n-space class="page-body" vertical :size="12">
    <n-grid cols="1 s:2 m:4" x-gap="12" y-gap="12" responsive="screen">
      <n-gi>
        <x-statistic :title="t('objects.node', 2)">
          <template #icon>
            <server-outline />
          </template>
          <x-anchor url="/swarm/nodes">{{ summary.nodeCount }}</x-anchor>
        </x-statistic>
      </n-gi>
      <n-gi>
        <x-statistic :title="t('objects.network', 2)">
          <template #icon>
            <globe-outline />
          </template>
          <x-anchor url="/swarm/networks">{{ summary.networkCount }}</x-anchor>
        </x-statistic>
      </n-gi>
      <n-gi>
        <x-statistic :title="t('objects.service', 2)">
          <template #icon>
            <image-outline />
          </template>
          <x-anchor url="/swarm/services">{{ summary.serviceCount }}</x-anchor>
        </x-statistic>
      </n-gi>
      <n-gi>
        <x-statistic :title="t('objects.stack', 2)">
          <template #icon>
            <albums-outline />
          </template>
          <x-anchor url="/swarm/stacks">{{ summary.stackCount }}</x-anchor>
        </x-statistic>
      </n-gi>
      <!-- <n-gi>
        <n-statistic label="任务" value="125" />
      </n-gi>-->
    </n-grid>
    <n-hr style="margin: 4px 0" />
    <x-dashboard type="home" />
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NSpace,
  NGrid,
  NGi,
  NHr,
} from "naive-ui";
import {
  ServerOutline,
  GlobeOutline,
  ImageOutline,
  AlbumsOutline,
} from "@vicons/ionicons5";
import XStatistic from "@/components/Statistic.vue";
import XAnchor from "@/components/Anchor.vue";
import XDashboard from "@/components/Dashboard.vue";
import systemApi from "@/api/system";
import type { Summary } from "@/api/system";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const summary = ref({
  nodeCount: 0,
  networkCount: 0,
  serviceCount: 0,
  stackCount: 0,
} as Summary)

async function initData() {
  const r = await systemApi.summarize();
  summary.value = r.data as Summary;
}

onMounted(() => {
  initData()
});
</script>
