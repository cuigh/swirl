<template>
  <n-space :size="24">
    <n-form-item :label="t('fields.rows')" label-placement="left">
      <n-input-number v-model:value="filters.lines" />
    </n-form-item>
    <n-form-item :label="t('fields.add_timestamps')" label-placement="left">
      <n-switch v-model:value="filters.timestamps" @update:value="fetchData" />
    </n-form-item>
    <n-form-item :label="t('fields.auto_refresh')" label-placement="left">
      <n-switch v-model:value="filters.refresh" @update:value="changeRefresh" />
    </n-form-item>
  </n-space>
  <n-tabs type="card" style="margin-top: -18px">
    <n-tab-pane name="stdout" tab="Stdout">
      <n-input
        type="textarea"
        :autosize="{ minRows: 5, maxRows: 30 }"
        :placeholder="''"
        readonly
        :value="logs.stdout"
      />
    </n-tab-pane>
    <n-tab-pane name="stderr" tab="Stderr">
      <n-input
        type="textarea"
        :autosize="{ minRows: 5, maxRows: 30 }"
        :placeholder="''"
        readonly
        :value="logs.stderr"
      />
    </n-tab-pane>
  </n-tabs>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, PropType, reactive, ref } from "vue";
import {
  NSpace,
  NInput,
  NInputNumber,
  NFormItem,
  NSwitch,
  NTabs,
  NTabPane,
} from "naive-ui";
import containerApi from "@/api/container";
import taskApi from "@/api/task";
import serviceApi from "@/api/service";
import { Result } from "@/api/ajax";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
interface Logs {
  stdout: string;
  stderr: string;
}

const props = defineProps({
  type: {
    type: String as PropType<'task' | 'container' | 'service'>,
    required: true,
  },
  id: {
    type: String,
    required: true,
  },
})
const filters = reactive({
  lines: 500,
  timestamps: false,
  refresh: false,
})
const logs = ref({
  stdout: "",
  stderr: "",
} as Logs)
const timer = ref();

function changeRefresh(value: boolean) {
  if (value) {
    refreshData()
  } else {
    clearTimeout(timer.value)
  }
}

async function fetchData() {
  var r: Result<Logs>;
  switch (props.type) {
    case 'container':
      r = await containerApi.fetchLogs({ id: props.id, lines: filters.lines, timestamps: filters.timestamps });
      break
    case 'task':
      r = await taskApi.fetchLogs({ id: props.id, lines: filters.lines, timestamps: filters.timestamps });
      break
    case 'service':
      r = await serviceApi.fetchLogs({ id: props.id, lines: filters.lines, timestamps: filters.timestamps });
      break
    default:
      return
  }
  logs.value = r.data as Logs;
}

function refreshData() {
  if (filters.refresh) {
    timer.value = setTimeout(() => {
      fetchData()
      refreshData()
    }, 3000);
  }
}

onMounted(() => {
  fetchData()
  refreshData()
});

onUnmounted(() => clearTimeout(timer.value));
</script>
