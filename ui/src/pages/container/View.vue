<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/local/containers')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')" display-directive="show:lazy">
        <n-space vertical :size="16">
          <x-description label-placement="left" label-align="right">
            <x-description-item :label="t('fields.id')" :span="2">{{ model.id }}</x-description-item>
            <x-description-item :label="t('fields.name')" :span="2">{{ model.name }}</x-description-item>
            <x-description-item :label="t('objects.image')" :span="2">{{ model.image }}</x-description-item>
            <x-description-item label="PID">{{ model.pid }}</x-description-item>
            <x-description-item :label="t('fields.state')">
              <n-tag
                round
                size="small"
                :type="model.state === 'running' ? 'success' : 'error'"
              >{{ model.state }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ model.createdAt }}</x-description-item>
            <x-description-item :label="t('fields.started_at')">{{ model.startedAt }}</x-description-item>
          </x-description>
          <x-panel :title="t('fields.labels')" v-if="model.labels && model.labels.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in model.labels">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="raw" :tab="t('fields.raw')" display-directive="show:lazy">
        <x-code :code="raw" language="json" />
      </n-tab-pane>
      <n-tab-pane name="logs" :tab="t('fields.logs')" display-directive="show:lazy">
        <x-logs type="container" :node="node" :id="model.id" v-if="store.getters.allow('container.logs')"></x-logs>
        <n-alert type="info" v-else>{{ t('texts.403') }}</n-alert>
      </n-tab-pane>
      <n-tab-pane name="exec" :tab="t('fields.execute')" display-directive="show:lazy">
        <execute :node="node" :id="model.id" v-if="store.getters.allow('container.execute')"></execute>
        <n-alert type="info" v-else>{{ t('texts.403') }}</n-alert>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTable,
  NTabs,
  NTabPane,
  NAlert,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import { useStore } from "vuex";
import XPageHeader from "@/components/PageHeader.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import XLogs from "@/components/Logs.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import containerApi from "@/api/container";
import type { Container } from "@/api/container";
import { useRoute } from "vue-router";
import Execute from "./modules/Execute.vue";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const store = useStore();
const model = ref({} as Container);
const raw = ref('');
const node = route.params.node as string || '';

async function fetchData() {
  const id = route.params.id as string;
  let r = await containerApi.find(node, id);
  model.value = r.data?.container as Container;
  raw.value = r.data?.raw as string;
}

onMounted(fetchData);
</script>