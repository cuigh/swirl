<template>
  <x-page-header :subtitle="model.title">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/system/charts')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
      <n-button
        secondary
        size="small"
        @click="$router.push(`/system/charts/${model.id}/edit`)"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="16">
    <x-description :label-width="90">
      <x-description-item :label="t('fields.id')">{{ model.id }}</x-description-item>
      <x-description-item :label="t('fields.name')">{{ model.title }}</x-description-item>
      <x-description-item :span="2" :label="t('fields.desc')">{{ model.desc }}</x-description-item>
      <x-description-item :label="t('fields.width')">{{ model.width }}</x-description-item>
      <x-description-item :label="t('fields.height')">{{ model.height }}</x-description-item>
      <x-description-item :label="t('fields.unit')">{{ model.unit }}</x-description-item>
      <x-description-item :label="t('fields.margin')">
        <n-space :size="4" v-if="model.margin">
          <x-pair-tag
            type="warning"
            :label="t('fields.left')"
            :value="model.margin.left.toString()"
            v-if="model.margin?.left"
          />
          <x-pair-tag
            type="warning"
            :label="t('fields.right')"
            :value="model.margin.right.toString()"
            v-if="model.margin?.right"
          />
          <x-pair-tag
            type="warning"
            :label="t('fields.top')"
            :value="model.margin.top.toString()"
            v-if="model.margin?.top"
          />
          <x-pair-tag
            type="warning"
            :label="t('fields.bottom')"
            :value="model.margin.bottom.toString()"
            v-if="model.margin?.bottom"
          />
        </n-space>
      </x-description-item>
      <x-description-item :label="t('fields.dashboard')">{{ model.dashboard }}</x-description-item>
      <x-description-item :label="t('fields.type')">{{ model.type }}</x-description-item>
      <x-description-item :label="t('fields.created_by')">
        <x-anchor :url="`/system/users/${model.createdBy?.id}`">{{ model.createdBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        <n-time :time="model.createdAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
      <x-description-item :label="t('fields.updated_by')">
        <x-anchor :url="`/system/users/${model.updatedBy?.id}`">{{ model.updatedBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        <n-time :time="model.updatedAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
    </x-description>
    <x-panel :title="t('fields.metrics')">
      <n-table size="small" :bordered="true" :single-line="false">
        <thead>
          <tr>
            <th>{{ t('fields.legend') }}</th>
            <th>{{ t('fields.query') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in model.metrics">
            <td>{{ m.legend }}</td>
            <td>{{ m.query }}</td>
          </tr>
        </tbody>
      </n-table>
    </x-panel>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NTable,
  NTime,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPairTag from "@/components/PairTag.vue";
import XPanel from "@/components/Panel.vue";
import XAnchor from "@/components/Anchor.vue";
import chartApi from "@/api/chart";
import type { Chart } from "@/api/chart";
import { useRoute } from "vue-router";
import { XDescription, XDescriptionItem } from "@/components/description";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Chart);

async function fetchData() {
  const id = route.params.id as string;
  let r = await chartApi.find(id);
  model.value = r.data as Chart;
}

onMounted(fetchData);
</script>