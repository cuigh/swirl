<template>
  <x-page-header :subtitle="model.name || model.id">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/tasks')">
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
          <x-description label-placement="left" label-align="right" :label-width="90">
            <x-description-item :label="t('fields.id')">{{ model.id }}</x-description-item>
            <x-description-item :label="t('objects.image')">{{ model.image }}</x-description-item>
            <x-description-item :label="t('objects.service')" :span="2">
              <x-anchor :url="`/swarm/services/${model.serviceName}`">{{ model.serviceName }}</x-anchor>
            </x-description-item>
            <x-description-item :label="t('objects.container')" :span="2">
              <x-anchor :url="`/local/containers/${model.nodeId}/${model.containerId}`">{{ model.containerId }}</x-anchor>
            </x-description-item>
            <x-description-item :label="t('objects.node')" :span="2">
              <x-anchor :url="`/swarm/nodes/${model.nodeId}`">{{ model.nodeId }}</x-anchor>
            </x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ model.createdAt }}</x-description-item>
            <x-description-item :label="t('fields.updated_at')">{{ model.updatedAt }}</x-description-item>
            <x-description-item label="PID">{{ model.pid }}</x-description-item>
            <x-description-item :label="t('fields.exit_code')">{{ model.exitCode }}</x-description-item>
            <x-description-item :label="t('fields.state')">
              <n-tag
                round
                size="small"
                :type="model.state === 'running' ? 'success' : 'default'"
              >{{ model.state }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.message')">{{ model.message }}</x-description-item>
            <x-description-item :label="t('fields.error')" v-if="model.error" :span="2">
              <n-text type="error">{{ model.error }}</n-text>
            </x-description-item>
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
          <x-panel :title="t('objects.network', 2)" v-if="model.networks && model.networks.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.address') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="n in model.networks">
                  <td>
                    <x-anchor :url="`/swarm/networks/${n.name}`">{{ n.name }}</x-anchor>
                  </td>
                  <td>
                    <n-space :size="4">
                      <n-tag round size="small" v-for="ip in n.ips">{{ ip }}</n-tag>
                    </n-space>
                  </td>
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
        <x-logs type="task" :id="route.params.id as string"></x-logs>
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
  NText,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import XLogs from "@/components/Logs.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import taskApi from "@/api/task";
import type { Task } from "@/api/task";
import { useRoute } from "vue-router";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Task);
const raw = ref('');

async function fetchData() {
  const id = route.params.id as string;
  let r = await taskApi.find(id);
  model.value = r.data?.task as Task;
  raw.value = r.data?.raw as string;
}

onMounted(fetchData);
</script>
