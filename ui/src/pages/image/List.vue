<template>
  <x-page-header />
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
} from "naive-ui";
import XPageHeader from "@/components/PageHeader.vue";
import imageApi from "@/api/image";
import type { Image } from "@/api/image";
import nodeApi from "@/api/node";
import { useDataTable } from "@/utils/data-table";
import { formatSize, renderButton, renderLink, renderTags } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  node: '',
  name: '',
});
const nodes: any = ref([])
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    fixed: "left" as const,
    render: (i: Image) => renderLink({ name: 'image_detail', params: { node: filter.node || '-', id: i.id } }, i.id.substring(7, 19)),
  },
  {
    title: t('fields.tags'),
    key: "tags",
    render(i: Image) {
      if (i.tags) {
        return renderTags(i.tags?.map(t => {
          return { text: t, type: 'default' }
        }), true, 6)
      }
    },
  },
  {
    title: t('fields.size'),
    key: "size",
    render(i: Image) {
      return formatSize(i.size)
    }
  },
  {
    title: t('fields.created_at'),
    key: "created"
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(i: Image, index: number) {
      return renderButton('error', t('buttons.delete'), () => deleteImage(i.id, index), t('prompts.delete'))
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(imageApi.search, filter, false)

async function deleteImage(id: string, index: number) {
  await imageApi.delete(filter.node, id, "");
  state.data.splice(index, 1)
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