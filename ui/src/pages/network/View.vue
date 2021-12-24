<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'network_list' })">
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
          <x-description label-placement="left" label-align="right">
            <x-description-item :label="t('fields.id')" :span="2">{{ model.id }}</x-description-item>
            <x-description-item :label="t('fields.created_at')" :span="2">{{ model.created }}</x-description-item>
            <x-description-item :label="t('fields.driver')">
              <n-tag
                round
                size="small"
                :type="model.driver === 'overlay' ? 'success' : 'default'"
              >{{ model.driver }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.scope')">
              <n-tag
                round
                size="small"
                :type="model.scope === 'swarm' ? 'success' : 'default'"
              >{{ model.scope }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.internal')">
              <n-tag
                round
                size="small"
                :type="model.internal ? 'success' : 'default'"
              >{{ model.internal ? t('enums.yes') : t('enums.no') }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.attachable')">
              <n-tag
                round
                size="small"
                :type="model.attachable ? 'success' : 'default'"
              >{{ model.attachable ? t('enums.yes') : t('enums.no') }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.ingress')">
              <n-tag
                round
                size="small"
                :type="model.ingress ? 'success' : 'default'"
              >{{ model.ingress ? t('enums.yes') : t('enums.no') }}</n-tag>
            </x-description-item>
            <x-description-item label="IPv6">
              <n-tag
                round
                size="small"
                :type="model.ipv6 ? 'success' : 'default'"
              >{{ model.ipv6 ? t('enums.enabled') : t('enums.disabled') }}</n-tag>
            </x-description-item>
          </x-description>
          <x-panel
            :title="t('fields.ipam')"
            :subtitle="`(${t('fields.driver')}: ${model.ipam?.driver})`"
            v-if="model.ipam?.config && model.ipam?.config.length"
          >
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.subnet') }}</th>
                  <th>{{ t('fields.gateway') }}</th>
                  <th>{{ t('fields.range') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in model.ipam?.config">
                  <td>{{ c.subnet }}</td>
                  <td>{{ c.gateway }}</td>
                  <td>{{ c.range }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
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
                  <th>名称</th>
                  <th>值</th>
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
          <x-panel title="Containers" v-if="model.containers && model.containers.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>IPv4</th>
                  <th>IPv6</th>
                  <th>Mac</th>
                  <th>{{ t('fields.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in model.containers">
                  <td>
                    <x-anchor
                      :url="{ name: 'container_detail', params: { node: '-', id: c.id } }"
                    >{{ c.name }}</x-anchor>
                  </td>
                  <td>{{ c.ipv4 }}</td>
                  <td>{{ c.ipv6 }}</td>
                  <td>{{ c.mac }}</td>
                  <td>
                    <n-popconfirm :show-icon="false" @positive-click="disconnect(c.id)">
                      <template #trigger>
                        <n-button size="tiny" ghost type="error">{{ t('buttons.disconnect') }}</n-button>
                      </template>
                      {{ t('prompts.disconnect') }}
                    </n-popconfirm>
                  </td>
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
  NPopconfirm,
  NTabs,
  NTabPane,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import networkApi from "@/api/network";
import type { Network } from "@/api/network";
import { useRoute } from "vue-router";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Network);
const raw = ref('');

async function disconnect(container: string) {
  await networkApi.disconnect(model.value.id, model.value.name, container);
}

async function fetchData() {
  const name = route.params.name as string;
  let r = await networkApi.find(name);
  model.value = r.data?.network as Network;
  raw.value = r.data?.raw as string;
}

onMounted(fetchData);
</script>