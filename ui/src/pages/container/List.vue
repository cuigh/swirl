<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" type="warning" @click="prune">
        <template #icon>
          <n-icon>
            <close-icon />
          </n-icon>
        </template>
        {{ t('buttons.prune') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-select
        filterable
        size="small"
        :consistent-menu-width="false"
        :placeholder="t('objects.node')"
        v-model:value="filter.node"
        :options="nodes"
        style="width: 200px"
        v-if="nodes && nodes.length"
      />
      <n-input size="small" v-model:value="filter.name" :placeholder="t('fields.name')" clearable />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="row => row.name"
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
import { onMounted, reactive, ref } from "vue";
import {
  NSpace,
  NButton,
  NDataTable,
  NInput,
  NSelect,
  NIcon,
} from "naive-ui";
import { CloseOutline as CloseIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import containerApi from "@/api/container";
import type { Container } from "@/api/container";
import nodeApi from "@/api/node";
import { useDataTable } from "@/utils/data-table";
import { formatSize, renderButton, renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  node: '',
  name: '',
});
const nodes: any = ref([])
const columns = [
  {
    title: t('fields.name'),
    key: "name",
    fixed: "left" as const,
    render: (c: Container) => {
      const node = c.labels?.find(l => l.name === 'com.docker.swarm.node.id')
      const name = c.name.length > 32 ? c.name.substring(0, 32) + '...' : c.name
      return renderLink({ name: 'container_detail', params: { id: c.id, node: node?.value || '-' } }, name)
    },
  },
  {
    title: t('objects.service'),
    key: "service",
    render: (c: Container) => {
      const service = c.labels?.find(l => l.name === 'com.docker.swarm.service.name')?.value
      return service ? renderLink({ name: 'service_detail', params: { name: service } }, service) : ''
    },
  },
  {
    title: t('objects.image'),
    key: "image",
  },
  {
    title: t('fields.state'),
    key: "state",
    render(c: Container) {
      return renderTag(c.state, c.state === 'running' ? 'success' : 'error')
    }
  },
  {
    title: t('fields.created_at'),
    key: "createdAt"
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(i: Container, index: number) {
      return renderButton('error', t('buttons.delete'), () => remove(i, index), t('prompts.delete'))
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(containerApi.search, filter, false)

async function remove(c: Container, index: number) {
  const node = c.labels?.find(l => l.name === 'com.docker.swarm.node.id')
  await containerApi.delete(node?.value || '', c.id, c.name);
  state.data.splice(index, 1)
}

async function prune() {
  window.dialog.warning({
    title: t('dialogs.prune_container.title'),
    content: t('dialogs.prune_container.body'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      const r = await containerApi.prune(filter.node);
      window.message.info(t('texts.prune_container_success', {
        count: r.data?.count,
        size: formatSize(r.data?.size as number),
      }));
      fetchData();
    }
  })
}

onMounted(async () => {
  const r = await nodeApi.list(true)
  nodes.value = r.data?.map(n => ({ label: n.name, value: n.id }))
  if (r.data?.length) {
    filter.node = r.data[0].id
  }
  fetchData()
})
</script>