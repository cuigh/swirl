<template>
  <x-page-header :subtitle="model.name || model.id">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'config_list' })">
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
        @click="$router.push({ name: 'config_edit', params: { id: model.id } })"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')">
        <n-space vertical :size="16">
          <x-description :label-width="90">
            <x-description-item :label="t('fields.id')">{{ model.id }}</x-description-item>
            <x-description-item :label="t('fields.name')">{{ model.name }}</x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ model.createdAt }}</x-description-item>
            <x-description-item :label="t('fields.updated_at')">{{ model.updatedAt }}</x-description-item>
          </x-description>
          <x-panel :title="t('fields.content')">
            <x-code :code="model.data" />
          </x-panel>
          <x-panel :title="t('fields.labels')" v-if="model.labels && model.labels.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in model.labels">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('fields.template')" v-if="model.templating.name">
            <x-description label-align="left" :label-width="40">
              <x-description-item :label="t('fields.engine')" :span="2">
                <n-tag round size="small" type="warning">{{ model.templating.name }}</n-tag>
              </x-description-item>
            </x-description>
            <p style="margin: 6px 0">{{ t('fields.options') }}</p>
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in model.templating.options">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="raw" :tab="t('fields.raw')">
        <x-code :code="raw" language="json" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NTable,
  NTabs,
  NTabPane,
  NTag,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import configApi from "@/api/config";
import type { Config } from "@/api/config";
import { useRoute } from "vue-router";
import { XDescription, XDescriptionItem } from "@/components/description";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({ templating: {} } as Config);
const raw = ref('');

async function fetchData() {
  const id = route.params.id as string;
  let r = await configApi.find(id);
  model.value = r.data?.config as Config;
  raw.value = r.data?.raw as string;
}

onMounted(fetchData);
</script>