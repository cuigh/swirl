<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" type="warning" @click="prune">
        <template #icon>
          <n-icon>
            <close-icon />
          </n-icon>
        </template>
        {{ t('buttons.prune') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-select
        size="small"
        :placeholder="t('fields.type')"
        v-model:value="filter.type"
        :options="types"
        style="width: 140px"
        clearable
      />
      <n-input
        size="small"
        v-model:value="filter.name"
        :placeholder="t('fields.object')"
        clearable
      />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="row => row.name"
      size="small"
      :columns="columns"
      :data="state.data"
      :pagination="pagination"
      :loading="state.loading"
      @update:page="fetchData"
      @update-page-size="changePageSize"
      scroll-x="max-content"
    />
  </n-space>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import {
  NSpace,
  NButton,
  NDataTable,
  NSelect,
  NInput,
  NIcon,
} from "naive-ui";
import { CloseOutline as CloseIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import eventApi from "@/api/event";
import type { Event } from "@/api/event";
import { useDataTable } from "@/utils/data-table";
import { renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'
import type { RouteLocationRaw } from "vue-router";

const { t } = useI18n()
const filter = reactive({
  type: undefined,
  name: "",
});
const types: any = [
  {
    type: 'group',
    label: 'System',
    key: 'system',
    children: [
      {
        label: 'User',
        value: 'User'
      },
      {
        label: 'Role',
        value: 'Role'
      },
      {
        label: 'Chart',
        value: 'Chart'
      },
      {
        label: 'Setting',
        value: 'Setting'
      },
    ],
  },
  {
    type: 'group',
    label: 'Swarm',
    key: 'swarm',
    children: [
      {
        label: 'Registry',
        value: 'Registry'
      },
      {
        label: 'Node',
        value: 'Node'
      },
      {
        label: 'Network',
        value: 'Network'
      },
      {
        label: 'Service',
        value: 'Service'
      },
      {
        label: 'Stack',
        value: 'Stack'
      },
      {
        label: 'Secret',
        value: 'Secret'
      },
      {
        label: 'Config',
        value: 'Config'
      },
    ],
  },
  {
    type: 'group',
    label: 'Local',
    key: 'local',
    children: [
      {
        label: 'Image',
        value: 'Image'
      },
      {
        label: 'Container',
        value: 'Container'
      },
      {
        label: 'Volume',
        value: 'Volume'
      },
    ],
  },
]
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    width: 210,
    fixed: "left" as const,
    // render: (e: Event) => renderLink(`/system/events/${e.id}`, e.id),
  },
  {
    title: t('fields.type'),
    key: "type",
    render(e: Event) {
      return renderTag(e.type)
    },
  },
  {
    title: t('fields.action'),
    key: "action",
    render(e: Event) {
      return renderTag(e.action)
    },
  },
  {
    title: t('fields.object'),
    key: "name",
    render(e: Event) {
      const u = url(e)
      return u ? renderLink(u, e.name) : e.name
    },
  },
  {
    title: t('fields.operator'),
    key: "name",
    render: (e: Event) => renderLink(`/system/users/${e.userId}`, e.username),
  },
  {
    title: t('fields.time'),
    key: "time",
    render: (e: Event) => renderTime(e.time),
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(eventApi.search, filter)

function url(e: Event): RouteLocationRaw | null {
  if (e.type === 'Setting') {
    return { name: 'setting' }
  } else if (!e.code) {
    return null
  }

  switch (e.type) {
    case "User":
      return { name: 'user_detail', params: { id: e.code } }
    case "Role":
      return { name: 'role_detail', params: { id: e.code } }
    case "Chart":
      return { name: 'chart_detail', params: { id: e.code } }
    case "Registry":
      return { name: 'registry_detail', params: { id: e.code } }
    case "Node":
      return { name: 'node_detail', params: { id: e.code } }
    case "Network":
      return { name: 'network_detail', params: { name: e.code } }
    case "Service":
      return { name: 'service_detail', params: { name: e.code } }
    case "Stack":
      return { name: 'stack_detail', params: { name: e.code } }
    case "Config":
      return { name: 'config_detail', params: { id: e.code } }
    case "Secret":
      return { name: 'secret_detail', params: { id: e.code } }
    case "Image":
      return { name: 'image_detail', params: { node: '-', id: e.code } }
    case "Container":
      return { name: 'container_detail', params: { node: '-', id: e.code } }
    case "Volume":
      return { name: 'volume_detail', params: { node: '-', name: e.code } }
  }
  return null
}

function prune() {
  window.message.info("TODO...")
}
</script>