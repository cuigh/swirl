<template>
  <x-page-header :subtitle="t('texts.records', { total: model.length }, model.length)">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'registry_new' })">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>
        {{ t('buttons.new') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-table size="small" :bordered="true" :single-line="false">
      <thead>
        <tr>
          <th>{{ t('fields.name') }}</th>
          <th>{{ t('fields.address') }}</th>
          <th>{{ t('fields.login_name') }}</th>
          <th>{{ t('fields.created_at') }}</th>
          <th>{{ t('fields.updated_at') }}</th>
          <th>{{ t('fields.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, index) of model" :key="r.id">
          <td>
            <x-anchor :url="{ name: 'registry_detail', params: { id: r.id } }">{{ r.name }}</x-anchor>
          </td>
          <td>{{ r.url }}</td>
          <td>{{ r.username }}</td>
          <td>
            <n-time :time="r.createdAt" format="y-MM-dd HH:mm:ss" />
          </td>
          <td>
            <n-time :time="r.updatedAt" format="y-MM-dd HH:mm:ss" />
          </td>
          <td>
            <n-button
              size="tiny"
              quaternary
              type="warning"
              @click="$router.push({ name: 'registry_edit', params: { id: r.id } })"
            >{{ t('buttons.edit') }}</n-button>
            <n-popconfirm :show-icon="false" @positive-click="deleteRegistry(r.id, index)">
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
  NIcon,
  NTime,
} from "naive-ui";
import { AddOutline as AddIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import registryApi from "@/api/registry";
import type { Registry } from "@/api/registry";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref([] as Registry[]);

async function deleteRegistry(id: string, index: number) {
  await registryApi.delete(id);
  model.value.splice(index, 1)
}

async function fetchData() {
  let r = await registryApi.search();
  model.value = r.data || [];
}

onMounted(fetchData);
</script>