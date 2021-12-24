<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="importChart">
        <template #icon>
          <n-icon>
            <archive-icon />
          </n-icon>
        </template>
        {{ t('buttons.import') }}
      </n-button>
      <n-button secondary size="small" @click="$router.push({ name: 'chart_new' })">
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
      <n-input
        size="small"
        v-model:value="filter.title"
        :placeholder="t('fields.title')"
        clearable
      />
      <n-select
        size="small"
        :placeholder="t('fields.dashboard')"
        v-model:value="filter.dashboard"
        :options="[{ label: 'Home', value: 'home' }, { label: 'Service', value: 'service' }]"
        style="width: 140px"
        clearable
      />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="(r: Chart) => r.id"
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
import { h, onMounted, reactive } from "vue";
import {
  NSpace,
  NButton,
  NDataTable,
  NInput,
  NIcon,
  NSelect,
} from "naive-ui";
import {
  ArchiveOutline as ArchiveIcon,
  AddOutline as AddIcon,
  PieChartOutline as PieChart,
  BarChartOutline as BarChart,
} from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import XCode from "@/components/Code.vue";
import XIcon from "@/components/Icon.vue";
import chartApi from "@/api/chart";
import type { Chart } from "@/api/chart";
import { renderButtons, renderLink, renderTag, renderTime } from "@/utils/render";
import { useDataTable } from "@/utils/data-table";
import { toTitle } from "@/utils";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const filter = reactive({
  title: "",
  dashboard: undefined,
});
const columns = [
  {
    title: t('fields.title'),
    key: "title",
    fixed: "left" as const,
    render: (c: Chart) => renderLink({ name: 'chart_detail', params: { id: c.id } }, c.title),
  },
  {
    title: t('fields.type'),
    key: "type",
    render: (c: Chart) => {
      switch (c.type) {
        case 'line':
          return h(XIcon, { size: 20, viewBox: '0 0 32 32', path: 'M4.67 28l6.39-12l7.3 6.49a2 2 0 0 0 1.7.47a2 2 0 0 0 1.42-1.07L27 10.9l-1.82-.9l-5.49 11l-7.3-6.49a2 2 0 0 0-1.68-.51a2 2 0 0 0-1.42 1L4 25V2H2v26a2 2 0 0 0 2 2h26v-2z' })
        case 'bar':
          return h(NIcon, { size: 20, style: 'vertical-align: middle' }, { default: () => h(BarChart) })
        case 'pie':
          return h(NIcon, { size: 20, style: 'vertical-align: middle' }, { default: () => h(PieChart) })
        case 'gauge':
          return h(XIcon, { size: 20, viewBox: '0 0 24 24', path: 'M7.934 16.066a.75.75 0 1 1-1.06 1.06a7.25 7.25 0 0 1 6.798-12.181a.75.75 0 1 1-.344 1.46a5.75 5.75 0 0 0-5.393 9.661zm9.954-6.924a.75.75 0 0 1 .955.46a7.25 7.25 0 0 1-1.716 7.524a.75.75 0 1 1-1.061-1.06a5.75 5.75 0 0 0 1.362-5.969a.75.75 0 0 1 .46-.955zm-2.009-2.475a.625.625 0 0 1 .962.761l-.13.25a354.691 354.691 0 0 1-1.415 2.713a154.8 154.8 0 0 1-1.156 2.157c-.171.31-.326.586-.452.803a4.964 4.964 0 0 1-.32.5a1.875 1.875 0 0 1-2.94-2.327c.086-.109.244-.265.413-.425c.182-.173.414-.387.678-.625a154.39 154.39 0 0 1 1.832-1.62a375.175 375.175 0 0 1 2.314-2.003l.214-.184zM22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10zM3.5 12a8.5 8.5 0 1 0 17 0a8.5 8.5 0 0 0-17 0z' })
        default:
          return renderTag(toTitle(c.type))
      }
    },
  },
  {
    title: t('fields.dashboard'),
    key: "dashboard",
    render: (c: Chart) => renderTag(toTitle(c.dashboard || 'any')),
  },
  {
    title: t('fields.width'),
    key: "width",
  },
  {
    title: t('fields.height'),
    key: "height",
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt",
    render: (c: Chart) => renderTime(c.updatedAt),
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(c: Chart, index: number) {
      return renderButtons([
        {
          type: 'error',
          text: t('buttons.delete'),
          action: () => deleteChart(c, index),
          prompt: t('prompts.delete'),
        },
        {
          type: 'warning',
          text: t('buttons.edit'),
          action: () => router.push({ name: 'chart_edit', params: { id: c.id } }),
        },
        {
          type: 'info',
          text: t('buttons.export'),
          action: () => exportChart(c),
        },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(chartApi.search, filter)

function importChart() {
  var text = ''
  window.dialog.success({
    showIcon: false,
    title: t('dialogs.import_chart.title'),
    content: () => h(NInput, {
      type: 'textarea',
      rows: 10,
      placeholder: t('dialogs.import_chart.tip'),
      onInput(v: string) { text = v },
    }),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    async onPositiveClick() {
      try {
        const c = JSON.parse(text)
        await chartApi.save(c)
        fetchData()
      } catch (e: any) {
        window.message.error(e.message)
        return false
      }
    }
  })
}

function exportChart(c: Chart) {
  const { id, createdAt, updatedAt, createdBy, updatedBy, ...chart } = c
  window.dialog.success({
    showIcon: false,
    title: t('dialogs.export_chart.title'),
    content: () => h(XCode, { language: 'javascript', code: JSON.stringify(chart, null, 2) }),
  })
}

async function deleteChart(c: Chart, index: number) {
  await chartApi.delete(c.id, c.title);
  state.data.splice(index, 1)
}

onMounted(fetchData);
</script>