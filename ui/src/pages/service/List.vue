<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/services/new')">
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
    <n-space :size="12">
      <n-input size="small" v-model:value="filter.name" :placeholder="t('fields.name')" clearable />
      <n-select
        size="small"
        :placeholder="t('fields.mode')"
        v-model:value="filter.mode"
        :options="modeOptions"
        style="width: 120px"
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
import { AddOutline as AddIcon } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import serviceApi from "@/api/service";
import type { Service } from "@/api/service";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter();
const filter = reactive({
  name: "",
  mode: undefined as string | undefined,
});
const modeOptions = [
  { label: 'Replicated', value: 'replicated' },
  { label: 'Global', value: 'global' },
  { label: 'Replicated Job', value: 'replicated-job' },
  { label: 'Global Job', value: 'global-job' },
];
const columns = [
  {
    title: t('fields.name'),
    key: "name",
    fixed: "left" as const,
    render: (s: Service) => renderLink(`/swarm/services/${s.name}`, s.name),
  },
  {
    title: t('objects.image'),
    key: "image",
  },
  {
    title: t('fields.mode'),
    key: "mode",
    render: (s: Service) => renderTag(s.mode, s.mode === 'replicated' || s.mode === 'replicated-job' ? "info" : "error"),
  },
  {
    title: t('fields.task'),
    key: "tasks",
    render: (s: Service) => {
      const type = s.desiredTasks === 0 ? 'warning' : (s.runningTasks === s.desiredTasks ? "success" : "error")
      return renderTag(`${s.runningTasks}/${s.desiredTasks}`, type)
    },
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt"
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(s: Service, index: number) {
      return renderButtons([
        {
          type: 'error',
          text: t('buttons.delete'),
          action: () => deleteService(s.name, index),
          prompt: t('prompts.delete'),
        },
        {
          type: 'warning',
          text: t('buttons.edit'),
          action: () => router.push(`/swarm/services/${s.name}/edit`),
        },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(serviceApi.search, filter)

async function deleteService(name: string, index: number) {
  await serviceApi.delete(name);
  state.data.splice(index, 1)
}
</script>