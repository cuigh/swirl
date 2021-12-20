<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input
        size="small"
        v-model:value="filter.service"
        :placeholder="t('fields.name')"
        clearable
      />
      <n-select
        size="small"
        :placeholder="t('fields.state')"
        v-model:value="filter.state"
        :options="stateOptions"
        style="width: 120px"
        clearable
      />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="(t: Task) => t.id"
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
import { reactive, ref } from "vue";
import {
  NSpace,
  NButton,
  NSelect,
  NInput,
  NDataTable,
} from "naive-ui";
import XPageHeader from "@/components/PageHeader.vue";
import taskApi from "@/api/task";
import type { Task } from "@/api/task";
import type { SearchArgs } from "@/api/task";
import { renderLink, renderTag } from "@/utils/render";
import { useDataTable } from "@/utils/data-table";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  service: "",
  state: undefined as string | undefined,
} as SearchArgs);
const stateOptions = [
  { label: 'Running', value: 'running' },
  { label: 'Shutdown', value: 'shutdown' },
  { label: 'Accepted', value: 'accepted' },
];
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    fixed: "left" as const,
    render: (s: Task) => renderLink({ name: 'task_detail', params: { id: s.id } }, s.id),
  },
  {
    title: t('fields.service_id'),
    key: "service",
    render: (s: Task) => renderLink({ name: 'service_detail', params: { name: s.serviceId } }, s.serviceId),
  },
  {
    title: t('objects.image'),
    key: "image",
  },
  {
    title: t('fields.node_id'),
    key: "image",
    render: (s: Task) => renderLink({ name: 'node_detail', params: { id: s.nodeId } }, s.nodeName),
  },
  {
    title: t('fields.state'),
    key: "mode",
    render: (t: Task) => renderTag(t.state, t.state === 'running' || t.state === 'starting' ? "success" : "default"),
  },
  {
    title: t('fields.created_at'),
    key: "createdAt"
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(taskApi.search, filter)
</script>