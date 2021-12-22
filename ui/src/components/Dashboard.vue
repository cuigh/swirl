<template>
  <n-space justify="space-between" style="margin-bottom: 6px">
    <n-form-item
      :label="t('fields.monitor')"
      label-placement="left"
      size="small"
      :show-feedback="false"
    >
      <n-select
        v-model:value="period"
        :options="options"
        @update-value="refreshData"
        style="width: 140px"
      />
    </n-form-item>
    <n-space :size="4">
      <n-button secondary size="small" @click="showDlg">
        <template #icon>
          <n-icon>
            <add-outline />
          </n-icon>
        </template>
        {{ t('buttons.add') }}
      </n-button>
      <n-button secondary @click="saveDashboard" size="small">
        <template #icon>
          <n-icon>
            <save-outline />
          </n-icon>
        </template>
        {{ t('buttons.save') }}
      </n-button>
    </n-space>
  </n-space>
  <n-grid id="charts" cols="1 640:12" x-gap="12" y-gap="12">
    <n-gi :key="c.id" :span="c.width" v-for="c in dashboard.charts">
      <x-chart
        :ref="setChartRefs"
        :info="c"
        :data-id="c.id"
        :data-width="c.width"
        :data-height="c.height"
        @remove="removeChart"
      />
    </n-gi>
  </n-grid>
  <n-modal
    v-model:show="modelDlg.visible"
    preset="card"
    :title="t('dialogs.add_chart.title')"
    size="small"
    style="width: 500px;"
  >
    <n-data-table
      size="small"
      :row-key="r => r.id"
      :columns="columns"
      :data="modelDlg.charts"
      @update:checked-row-keys="(ids: any) => modelDlg.checkedIds = ids"
      max-height="300"
    />
    <div style="display: flex; justify-content: flex-end; margin-top: 12px">
      <n-button
        type="primary"
        :disabled="modelDlg.checkedIds.length === 0"
        @click="addCharts"
      >{{ t('buttons.confirm') }}</n-button>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, PropType, reactive, ref } from "vue";
import {
  NSpace,
  NFormItem,
  NSelect,
  NButton,
  NIcon,
  NGrid,
  NGi,
  NModal,
  NDataTable,
} from "naive-ui";
import {
  AddOutline,
  SaveOutline,
} from "@vicons/ionicons5";
import { XChart } from "@/components/chart";
import dragula from 'dragula'
import 'dragula/dist/dragula.css'
import dashboardApi from "@/api/dashboard";
import type { Dashboard, ChartInfo } from "@/api/dashboard";
import chartApi from "@/api/chart";
import type { Chart } from "@/api/chart";
import { isEmpty, useTimer } from "@/utils";
import { useI18n } from 'vue-i18n'

const props = defineProps({
  type: {
    type: String as PropType<'home' | 'service' | 'task'>,
    required: true,
    default: 'home',
  },
  name: {
    type: String,
    default: '',
  },
})

const { t } = useI18n()
const options = [
  {
    label: t('texts.period_minute', { period: 30 }),
    value: '30',
  },
  {
    label: t('texts.period_hour', { period: 1 }, 1),
    value: '60',
  },
  {
    label: t('texts.period_hour', { period: 3 }, 3),
    value: '180',
  },
  {
    label: t('texts.period_hour', { period: 6 }, 6),
    value: '360',
  },
  {
    label: t('texts.period_hour', { period: 12 }, 12),
    value: '720',
  },
  {
    label: t('texts.period_hour', { period: 24 }, 24),
    value: '1440',
  },
]
const modelDlg = reactive({
  visible: false,
  checkedIds: [] as string[],
  charts: [] as Chart[],
})
const columns = [
  {
    type: 'selection' as const,
    disabled(row: ChartInfo) {
      return charts.has(row.id)
    }
  },
  {
    title: t('fields.title'),
    key: "title",
  },
  {
    title: t('fields.width'),
    key: "width",
    width: 60,
  },
  {
    title: t('fields.height'),
    key: "height",
    width: 70,
  },
];
const dashboard = ref({} as Dashboard)
const period = computed({
  get() { return dashboard.value?.period?.toString() || '30' },
  set(v: string) { dashboard.value.period = parseInt(v) },
})
const charts = new Map<string, any>()
const setChartRefs = (c: any) => c && charts.set(c.id, c)
var stop: () => void

async function saveDashboard() {
  await dashboardApi.save(dashboard.value);
  window.message.success(t('texts.action_success'))
}

function showDlg() {
  modelDlg.visible = true;
  chartApi.search({ dashboard: props.type, pageIndex: 1, pageSize: 100 }).then(r => {
    modelDlg.charts = r.data?.items as Chart[]
  })
}

function addCharts() {
  modelDlg.visible = false
  if (modelDlg.checkedIds.length > 0) {
    const set = new Set<string>(modelDlg.checkedIds)
    dashboard.value.charts = dashboard.value.charts || []
    modelDlg.charts.forEach(c => {
      if (set.has(c.id)) {
        dashboard.value.charts.push({
          id: c.id,
          title: c.title,
          type: c.type as 'line' | 'bar' | 'pie' | 'gauge',
          unit: c.unit,
          width: c.width,
          height: c.height,
        } as ChartInfo)
      }
    })
    refreshData()
  }
}

function removeChart(id: string) {
  window.dialog.warning({
    title: t('dialogs.remove_chart.title'),
    content: t('dialogs.remove_chart.content'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: () => {
      charts.delete(id)
      const i = dashboard.value.charts.findIndex(c => c.id === id)
      dashboard.value.charts.splice(i, 1)
    },
  })
}

async function refreshData() {
  if (dashboard.value != null && !isEmpty(dashboard.value.charts)) {
    const ids = dashboard.value.charts.map(i => i.id)
    const r = await dashboardApi.fetchData(props.name, ids, dashboard.value.period || 30);
    charts.forEach((v, k) => {
      r.data[k] && v.setData(r.data[k])
    })
  }
}

function initDrag() {
  dragula([document.getElementById('charts') as Element], {
    moves: function (el, container, handle): boolean {
      while (handle) {
        if (handle.classList.contains('drag-handle')) {
          return true
        }
        handle = handle.parentElement as Element
      }
      return false
    }
  }).on("dragend", (el: Element) => {
    sortCharts()
  });
}

function sortCharts() {
  const m = new Map<string, number>()
  const container = document.getElementById('charts') as Element
  for (var i = 0; i < container.children.length; i++) {
    const elem = container.children.item(i)?.firstElementChild as HTMLElement
    m.set(elem.dataset.id as string, i)
  }

  const charts = new Array<ChartInfo>(m.size)
  dashboard.value.charts.forEach(c => {
    charts[m.get(c.id) as number] = c
  });
  dashboard.value.charts = charts
}

onMounted(() => {
  dashboardApi.find(props.type, props.name).then(r => {
    dashboard.value = r.data as Dashboard
    initDrag()
    stop = useTimer(refreshData, (dashboard.value.interval || 10) * 1000)
  })
});

onUnmounted(() => stop && stop());
</script>
