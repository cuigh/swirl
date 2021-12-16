<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/stacks/new')">
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
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-table size="small" :bordered="true" :single-line="false">
      <thead>
        <tr>
          <th>{{ t('fields.name') }}</th>
          <th>{{ t('objects.service', {}, 2) }}</th>
          <th>{{ t('fields.created_at') }}</th>
          <th>{{ t('fields.updated_at') }}</th>
          <th>{{ t('fields.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, index) of model" :key="r.name">
          <td>
            <x-anchor :url="`/swarm/stacks/${r.name}`">{{ r.name }}</x-anchor>
          </td>
          <td>
            <n-space :size="4" v-if="r.services && r.services.length">
              <n-tag size="small" type="primary" v-for="s in r.services">
                <x-anchor :url="`/swarm/services/${s}`">{{ s.substring(r.name.length + 1) }}</x-anchor>
              </n-tag>
            </n-space>
          </td>
          <td>
            <n-time :time="r.createdAt" format="y-MM-dd HH:mm:ss" />
          </td>
          <td>
            <n-time :time="r.updatedAt" format="y-MM-dd HH:mm:ss" />
          </td>
          <td>
            <n-button
              size="tiny"
              quaternary
              type="warning"
              @click="$router.push({ name: 'stack_edit', params: { name: r.name } })"
            >{{ t('buttons.edit') }}</n-button>
            <n-popconfirm :show-icon="false" @positive-click="deployStack(r)">
              <template #trigger>
                <n-button size="tiny" quaternary type="warning">{{ t('buttons.deploy') }}</n-button>
              </template>
              {{ t('prompts.deploy') }}
            </n-popconfirm>
            <n-popconfirm
              :show-icon="false"
              @positive-click="shutdownStack(r)"
              v-if="r.services && r.services.length"
            >
              <template #trigger>
                <n-button size="tiny" quaternary type="error">{{ t('buttons.shutdown') }}</n-button>
              </template>
              {{ t('prompts.shutdown') }}
            </n-popconfirm>
            <n-popconfirm :show-icon="false" @positive-click="deleteStack(r.name, index)">
              <template #trigger>
                <n-button size="tiny" quaternary type="error">{{ t('buttons.delete') }}</n-button>
              </template>
              {{ t('prompts.delete') }}
            </n-popconfirm>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import {
  NSpace,
  NButton,
  NIcon,
  NInput,
  NTable,
  NPopconfirm,
  NTag,
  NTime,
} from "naive-ui";
import { AddOutline as AddIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import stackApi from "@/api/stack";
import type { Stack } from "@/api/stack";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref([] as Stack[]);
const filter = reactive({
  name: "",
  filter: "",
});

async function deleteStack(name: string, index: number) {
  await stackApi.delete(name);
  model.value.splice(index, 1)
}

async function shutdownStack(s: Stack) {
  await stackApi.shutdown(s.name);
  s.services = []
}

async function deployStack(s: Stack) {
  await stackApi.deploy(s.name);
  s.services = []
}

async function fetchData() {
  let r = await stackApi.search(filter);
  model.value = r.data || [];
}

onMounted(fetchData);
</script>