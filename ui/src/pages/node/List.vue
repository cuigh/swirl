<template>
  <x-page-header :subtitle="t('texts.records', { total: model.length }, model.length)" />
  <n-space class="page-body" vertical :size="12">
    <n-table size="small" :bordered="true" :single-line="false">
      <thead>
        <tr>
          <th>{{ t('fields.name') }}</th>
          <th>{{ t('fields.role') }}</th>
          <th>{{ t('fields.version') }}</th>
          <th>{{ t('fields.cpu') }}</th>
          <th>{{ t('fields.memory') }}</th>
          <th>{{ t('fields.address') }}</th>
          <th>{{ t('fields.state') }}</th>
          <th>{{ t('fields.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, index) of model" :key="r.id">
          <td>
            <x-anchor :url="`/swarm/nodes/${r.id}`">{{ r.name || r.hostname }}</x-anchor>
          </td>
          <td>
            <n-tag
              round
              size="small"
              :type="r.role === 'manager' ? (r.manager?.leader ? 'error' : 'primary') : 'default'"
            >{{ r.role }}</n-tag>
          </td>
          <td>{{ r.engineVersion }}</td>
          <td>{{ r.cpu }}</td>
          <td>{{ r.memory.toFixed(2) }} GB</td>
          <td>{{ r.address }}</td>
          <td>
            <n-tag
              round
              size="small"
              :type="r.state === 'ready' ? 'success' : 'error'"
            >{{ r.state }}</n-tag>
          </td>
          <td>
            <n-button
              size="tiny"
              quaternary
              type="warning"
              @click="$router.push({ name: 'node_edit', params: { id: r.id } })"
            >{{ t('buttons.edit') }}</n-button>
            <n-popconfirm :show-icon="false" @positive-click="deleteNode(r.id, index)">
              <template #trigger>
                <n-button size="tiny" quaternary type="error">{{ t('buttons.delete') }}</n-button>
              </template>
              {{ t('prompts.delete') }}
            </n-popconfirm>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NSpace,
  NButton,
  NTable,
  NPopconfirm,
  NTag,
} from "naive-ui";
import XAnchor from "@/components/Anchor.vue";
import XPageHeader from "@/components/PageHeader.vue";
import nodeApi from "@/api/node";
import type { Node } from "@/api/node";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref([] as Node[]);

async function deleteNode(id: string, index: number) {
  await nodeApi.delete(id);
  model.value.splice(index, 1)
}

async function fetchData() {
  let r = await nodeApi.search();
  model.value = r.data || [];
}

onMounted(fetchData);
</script>