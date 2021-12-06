<template>
  <x-page-header :subtitle="node.name ?? node.hostname">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/nodes')">
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
        @click="$router.push(`/swarm/nodes/${node.id}/edit`)"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')">
        <n-space vertical :size="16">
          <x-description label-placement="left" label-align="right" :label-width="105">
            <x-description-item :label="t('fields.id')">{{ node.id }}</x-description-item>
            <x-description-item :label="t('fields.hostname')">{{ node.hostname }}</x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ node.createdAt }}</x-description-item>
            <x-description-item :label="t('fields.updated_at')">{{ node.updatedAt }}</x-description-item>
            <x-description-item :label="t('fields.arch')">{{ node.arch }}</x-description-item>
            <x-description-item :label="t('fields.os')">{{ node.os }}</x-description-item>
            <x-description-item :label="t('fields.cpu')">{{ node.cpu }}</x-description-item>
            <x-description-item :label="t('fields.memory')">{{ node.memory?.toFixed(2) }} GB</x-description-item>
            <x-description-item :label="t('fields.engine_version')">{{ node.engineVersion }}</x-description-item>
            <x-description-item :label="t('fields.address')">{{ node.address }}</x-description-item>
            <x-description-item :label="t('fields.state')">
              <n-tag
                round
                size="small"
                :type="node.state === 'ready' ? 'success' : 'default'"
              >{{ node.state }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.availability')">
              <n-tag
                round
                size="small"
                :type="node.availability === 'active' ? 'success' : 'default'"
              >{{ node.availability }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('fields.role')">
              <n-space :size="6">
                <n-tag
                  round
                  size="small"
                  :type="node.role === 'manager' ? 'primary' : 'default'"
                >{{ node.role }}</n-tag>
                <n-tag round size="small" type="error" v-if="node.manager?.leader">leader</n-tag>
              </n-space>
            </x-description-item>
          </x-description>
          <x-panel title="Manager" v-if="node.manager">
            <x-description label-placement="left" label-align="right">
              <x-description-item :label="t('fields.reachability')">{{ node.manager.reachability }}</x-description-item>
              <x-description-item :label="t('fields.address')">{{ node.manager.addr }}</x-description-item>
            </x-description>
          </x-panel>
          <x-panel title="Labels" v-if="node.labels && node.labels.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>名称</th>
                  <th>值</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in node.labels">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel title="Tasks" v-if="tasks && tasks.length">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.id') }}</th>
                  <th>{{ t('fields.state') }}</th>
                  <th>{{ t('objects.image') }}</th>
                  <th>{{ t('fields.updated_at') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="t in tasks">
                  <td>
                    <x-anchor :url="`/swarm/tasks/${t.id}`">{{ t.id }}</x-anchor>
                  </td>
                  <td>
                    <n-tag
                      round
                      size="small"
                      :type="t.state === 'running' ? 'success' : 'error'"
                    >{{ t.state }}</n-tag>
                  </td>
                  <td>{{ t.image }}</td>
                  <td>{{ t.updatedAt }}</td>
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
import { onMounted, reactive, ref } from "vue";
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
import XAnchor from "@/components/Anchor.vue";
import XCode from "@/components/Code.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import { useRoute } from "vue-router";
import nodeApi from "@/api/node";
import type { Node } from "@/api/node";
import taskApi from "@/api/task";
import type { Task } from "@/api/task";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const node = ref({} as Node)
const tasks = ref([] as Task[])
const raw = ref('')

async function fetchData() {
  const id = route.params.id as string
  let results = await Promise.all([
    nodeApi.find(id),
    taskApi.search({ node: id, pageIndex: 1, pageSize: 100 }),
  ])

  node.value = results[0].data?.node as Node;
  raw.value = results[0].data?.raw as string;
  tasks.value = results[1].data?.items as Task[];
}

onMounted(fetchData);
</script>