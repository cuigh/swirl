<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button size="small" @click="$router.push('/local/volumes')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')">
        <n-space vertical :size="16">
          <x-description label-placement="left" label-align="right" :label-width="90">
            <x-description-item :label="t('fields.name')" :span="2">{{ model.name }}</x-description-item>
            <x-description-item :label="t('fields.driver')">
              <n-tag round size="small">{{ model.driver }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.scope')">
              <n-tag round size="small">{{ model.scope }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.size')" v-if="model.size != -1">{{ model.size }}</x-description-item>
            <x-description-item
              :label="t('fields.ref_count')"
              v-if="model.refCount != -1"
            >{{ model.refCount }}</x-description-item>
            <x-description-item :label="t('fields.mount_point')" :span="2">{{ model.mountPoint }}</x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ model.createdAt }}</x-description-item>
          </x-description>
          <x-panel :title="t('fields.options')" v-if="model.options && model.options.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="opt in model.options">
                  <td>{{ opt.name }}</td>
                  <td>{{ opt.value }}</td>
                </tr>
              </tbody>
            </n-table>
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
  NTag,
  NSpace,
  NIcon,
  NTable,
  NTabs,
  NTabPane,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import XCode from "@/components/Code.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import volumeApi from "@/api/volume";
import type { Volume } from "@/api/volume";
import { useRoute } from "vue-router";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Volume);
const raw = ref('');

async function fetchData() {
  const name = route.params.name as string;
  let r = await volumeApi.find(name);
  model.value = r.data?.volume as Volume;
  raw.value = r.data?.raw as string;
}

onMounted(fetchData);
</script>