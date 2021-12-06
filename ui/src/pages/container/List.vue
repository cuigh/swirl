<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input size="small" v-model:value="filter.name" :placeholder="t('fields.name')" clearable />
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
  NInput,
} from "naive-ui";
import XPageHeader from "@/components/PageHeader.vue";
import containerApi from "@/api/container";
import type { Container } from "@/api/container";
import { useDataTable } from "@/utils/data-table";
import { renderButton, renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  name: "",
});
const columns = [
  {
    title: t('fields.name'),
    key: "name",
    fixed: "left" as const,
    render: (c: Container) => renderLink(`/local/containers/${c.id}`, c.name),
  },
  {
    title: t('objects.image'),    
    key: "image",
  },
  {
    title: t('fields.state'),
    key: "state",
    render(c: Container) {
      return renderTag(c.state, c.state === 'running' ? 'success' : 'error')
    }
  },
  {
    title: t('fields.created_at'),
    key: "createdAt"
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(i: Container, index: number) {
      return renderButton('error', t('buttons.delete'), () => deleteContainer(i.id, index), t('prompts.delete'))
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(containerApi.search, filter)

async function deleteContainer(id: string, index: number) {
  await containerApi.delete(id, "");
  state.data.splice(index, 1)
}
</script>