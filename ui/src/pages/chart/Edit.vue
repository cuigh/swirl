<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push('/system/charts')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-form :model="model" :rules="rules" ref="form" label-placement="top" label-width="90">
      <n-grid cols="1 640:3" :x-gap="24">
        <n-form-item-gi :label="t('fields.title')" path="title">
          <n-input :placeholder="t('fields.title')" v-model:value="model.title" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.desc')" path="desc" span="2">
          <n-input :placeholder="t('fields.title')" v-model:value="model.desc" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.width')" path="width">
          <n-select v-model:value="model.width" :options="widths" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.height')" path="height">
          <n-input-number
            :placeholder="t('fields.height')"
            v-model:value="model.height"
            style="width: 100%"
          />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.unit')" path="unit">
          <n-select v-model:value="model.unit" :options="units" />
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.dashboard')"
          path="dashboard"
          span="3"
          label-placement="left"
          label-align="left"
        >
          <n-radio-group v-model:value="model.dashboard">
            <n-radio :key="d.key" :value="d.value" v-for="d in dashboards">{{ d.label }}</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.type')"
          path="type"
          span="3"
          label-placement="left"
          label-align="left"
        >
          <n-radio-group v-model:value="model.type">
            <n-radio key="line" value="line">Line</n-radio>
            <n-radio key="bar" value="bar">Bar</n-radio>
            <n-radio key="pie" value="pie">Pie</n-radio>
            <n-radio key="gauge" value="gauge">Gauge</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.margin')"
          path="margin"
          span="2"
          label-placement="left"
          label-align="left"
          v-if="model.type === 'line' || model.type === 'bar'"
        >
          <n-input-group>
            <n-input-group-label>{{ t('fields.left') }}</n-input-group-label>
            <n-input-number :min="0" :placeholder="''" v-model:value="model.margin.left" />
            <n-input-group-label>{{ t('fields.right') }}</n-input-group-label>
            <n-input-number :min="0" :placeholder="''" v-model:value="model.margin.right" />
            <n-input-group-label>{{ t('fields.top') }}</n-input-group-label>
            <n-input-number :min="0" :placeholder="''" v-model:value="model.margin.top" />
            <n-input-group-label>{{ t('fields.bottom') }}</n-input-group-label>
            <n-input-number :min="0" :placeholder="''" v-model:value="model.margin.bottom" />
          </n-input-group>
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.metrics')" path="metrics">
          <n-dynamic-input
            v-model:value="model.metrics"
            #="{ index, value }"
            :on-create="newMetric"
          >
            <n-input
              :placeholder="t('tips.legend')"
              v-model:value="value.legend"
              style="width: 250px"
            />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('tips.query')" v-model:value="value.query" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-gi :span="2">
          <n-button
            @click.prevent="submit"
            type="primary"
            :disabled="submiting"
            :loading="submiting"
          >
            <template #icon>
              <n-icon>
                <save-icon />
              </n-icon>
            </template>
            {{ t('buttons.save') }}
          </n-button>
        </n-gi>
      </n-grid>
    </n-form>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NIcon,
  NForm,
  NGrid,
  NGi,
  NFormItemGi,
  NInputNumber,
  NSelect,
  NDynamicInput,
  NRadioGroup,
  NRadio,
  NInputGroup,
  NInputGroupLabel,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import chartApi from "@/api/chart";
import type { Chart } from "@/api/chart";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({
  type: 'line',
  dashboard: 'home',
  width: 12,
  height: 200,
  unit: '',
  margin: {},
} as Chart);
const widths = [
  { label: '1', value: 1 },
  { label: '2', value: 2 },
  { label: '3', value: 3 },
  { label: '4', value: 4 },
  { label: '5', value: 5 },
  { label: '6', value: 6 },
  { label: '7', value: 7 },
  { label: '8', value: 8 },
  { label: '9', value: 9 },
  { label: '10', value: 10 },
  { label: '11', value: 11 },
  { label: '12', value: 12 },
]
const units: any = [
  { label: 'None', value: '' },
  {
    type: 'group',
    key: 'Percent',
    label: 'Percent',
    children: [
      { label: '0-100', value: 'percent:100' },
      { label: '0.0-1.0', value: 'percent:1' },
    ],
  },
  {
    type: 'group',
    key: 'Time',
    label: 'Time',
    children: [
      { label: 'Nanoseconds', value: 'time:ns' },
      { label: 'Microseconds', value: 'time:µs' },
      { label: 'Milliseconds', value: 'time:ms' },
      { label: 'Seconds', value: 'time:s' },
      { label: 'Minutes', value: 'time:m' },
      { label: 'Hours', value: 'time:h' },
      { label: 'Days', value: 'time:d' },
    ],
  },
  {
    type: 'group',
    key: 'Size',
    label: 'Size',
    children: [
      { label: 'Bits', value: 'size:bits' },
      { label: 'Bytes', value: 'size:bytes' },
      { label: 'Kilobytes', value: 'size:kilobytes' },
      { label: 'Megabytes', value: 'size:megabytes' },
      { label: 'Gigabytes', value: 'size:gigabytes' },
    ],
  },
];
const dashboards = [
  { label: 'Any', value: '', key: 'any' },
  { label: 'Home', value: 'home', key: 'home' },
  { label: 'Service', value: 'service', key: 'service' },
]
const rules: any = {
  title: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => chartApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push("/system/charts")
})

function newMetric() {
  return { legend: '', query: '' }
}

async function fetchData() {
  const id = route.params.id as string
  if (id) {
    let r = await chartApi.find(id);
    model.value = r.data as Chart;
  }
}

onMounted(fetchData);
</script>